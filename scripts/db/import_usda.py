#!/usr/bin/env python3

import argparse
import json
import sqlite3
from pathlib import Path


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Import USDA FoundationFoods JSON into SQLite")
    parser.add_argument("--db-path", required=True, help="Path to sqlite database")
    parser.add_argument("--dataset-path", required=True, help="Path to USDA FoundationFoods json")
    return parser.parse_args()


def chunks(values, chunk_size: int):
    bucket = []
    for value in values:
        bucket.append(value)
        if len(bucket) >= chunk_size:
            yield bucket
            bucket = []
    if bucket:
        yield bucket


def main() -> None:
    args = parse_args()
    db_path = Path(args.db_path)
    dataset_path = Path(args.dataset_path)

    if not dataset_path.exists():
        raise SystemExit(f"Dataset not found: {dataset_path}")

    with dataset_path.open("r", encoding="utf-8") as f:
        payload = json.load(f)

    foods = payload.get("FoundationFoods", [])
    dataset_name = "usda_foundation_foods"
    dataset_version = dataset_path.stem

    conn = sqlite3.connect(str(db_path))
    try:
        conn.execute("PRAGMA foreign_keys = ON;")
        conn.execute("BEGIN")

        food_rows = []
        nutrient_rows = []
        food_nutrient_rows = []

        for food in foods:
            if not isinstance(food, dict):
                continue

            fdc_id = food.get("fdcId")
            if fdc_id is None:
                continue

            food_rows.append(
                (
                    int(fdc_id),
                    str(food.get("description") or ""),
                    str(food.get("foodClass") or ""),
                    dataset_name,
                    dataset_version,
                )
            )

            for fn in food.get("foodNutrients", []):
                if not isinstance(fn, dict):
                    continue

                nutrient = fn.get("nutrient") or {}
                nutrient_id = nutrient.get("id")
                if nutrient_id is None:
                    continue

                nutrient_rows.append(
                    (
                        int(nutrient_id),
                        str(nutrient.get("number") or ""),
                        str(nutrient.get("name") or ""),
                        str(nutrient.get("unitName") or ""),
                        nutrient.get("rank"),
                    )
                )

                derivation = fn.get("foodNutrientDerivation") or {}
                source = derivation.get("foodNutrientSource") or {}
                food_nutrient_rows.append(
                    (
                        fn.get("id"),
                        int(fdc_id),
                        int(nutrient_id),
                        fn.get("amount"),
                        fn.get("median"),
                        fn.get("min"),
                        fn.get("max"),
                        fn.get("dataPoints"),
                        derivation.get("code"),
                        derivation.get("description"),
                        source.get("code"),
                        source.get("description"),
                    )
                )

        for batch in chunks(food_rows, 1000):
            conn.executemany(
                """
                INSERT INTO usda_foods (fdc_id, description, food_class, source_dataset, source_version)
                VALUES (?, ?, ?, ?, ?)
                ON CONFLICT(fdc_id) DO UPDATE SET
                    description = excluded.description,
                    food_class = excluded.food_class,
                    source_dataset = excluded.source_dataset,
                    source_version = excluded.source_version,
                    updated_at = CURRENT_TIMESTAMP
                """,
                batch,
            )

        for batch in chunks(nutrient_rows, 1000):
            conn.executemany(
                """
                INSERT INTO usda_nutrients (nutrient_id, number, name, unit_name, rank)
                VALUES (?, ?, ?, ?, ?)
                ON CONFLICT(nutrient_id) DO UPDATE SET
                    number = excluded.number,
                    name = excluded.name,
                    unit_name = excluded.unit_name,
                    rank = excluded.rank
                """,
                batch,
            )

        for batch in chunks(food_nutrient_rows, 1000):
            conn.executemany(
                """
                INSERT INTO usda_food_nutrients (
                    food_nutrient_id,
                    fdc_id,
                    nutrient_id,
                    amount,
                    median,
                    min,
                    max,
                    data_points,
                    derivation_code,
                    derivation_description,
                    source_code,
                    source_description
                )
                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                ON CONFLICT(fdc_id, nutrient_id) DO UPDATE SET
                    food_nutrient_id = excluded.food_nutrient_id,
                    amount = excluded.amount,
                    median = excluded.median,
                    min = excluded.min,
                    max = excluded.max,
                    data_points = excluded.data_points,
                    derivation_code = excluded.derivation_code,
                    derivation_description = excluded.derivation_description,
                    source_code = excluded.source_code,
                    source_description = excluded.source_description,
                    updated_at = CURRENT_TIMESTAMP
                """,
                batch,
            )

        conn.execute(
            """
            INSERT INTO dataset_imports (
                dataset_name,
                dataset_version,
                source_path,
                row_count_foods,
                row_count_food_nutrients,
                imported_at
            )
            VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
            ON CONFLICT(dataset_name) DO UPDATE SET
                dataset_version = excluded.dataset_version,
                source_path = excluded.source_path,
                row_count_foods = excluded.row_count_foods,
                row_count_food_nutrients = excluded.row_count_food_nutrients,
                imported_at = CURRENT_TIMESTAMP
            """,
            (
                dataset_name,
                dataset_version,
                str(dataset_path),
                len(food_rows),
                len(food_nutrient_rows),
            ),
        )

        conn.commit()
        print(
            f"Imported USDA foods={len(food_rows)} nutrients={len(nutrient_rows)} "
            f"food_nutrients={len(food_nutrient_rows)}"
        )
    except Exception:
        conn.rollback()
        raise
    finally:
        conn.close()


if __name__ == "__main__":
    main()
