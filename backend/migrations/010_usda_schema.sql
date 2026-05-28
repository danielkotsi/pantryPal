PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS usda_foods (
    fdc_id INTEGER PRIMARY KEY,
    description TEXT NOT NULL,
    food_class TEXT,
    source_dataset TEXT NOT NULL,
    source_version TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS usda_nutrients (
    nutrient_id INTEGER PRIMARY KEY,
    number TEXT,
    name TEXT NOT NULL,
    unit_name TEXT,
    rank REAL
);

CREATE TABLE IF NOT EXISTS usda_food_nutrients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    food_nutrient_id INTEGER,
    fdc_id INTEGER NOT NULL,
    nutrient_id INTEGER NOT NULL,
    amount REAL,
    median REAL,
    min REAL,
    max REAL,
    data_points INTEGER,
    derivation_code TEXT,
    derivation_description TEXT,
    source_code TEXT,
    source_description TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (fdc_id) REFERENCES usda_foods(fdc_id) ON DELETE CASCADE,
    FOREIGN KEY (nutrient_id) REFERENCES usda_nutrients(nutrient_id) ON DELETE RESTRICT,
    UNIQUE (fdc_id, nutrient_id),
    UNIQUE (food_nutrient_id)
);

CREATE TABLE IF NOT EXISTS dataset_imports (
    dataset_name TEXT PRIMARY KEY,
    dataset_version TEXT NOT NULL,
    source_path TEXT NOT NULL,
    row_count_foods INTEGER NOT NULL DEFAULT 0,
    row_count_food_nutrients INTEGER NOT NULL DEFAULT 0,
    imported_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_usda_foods_description ON usda_foods(description);
CREATE INDEX IF NOT EXISTS idx_usda_food_nutrients_fdc_id ON usda_food_nutrients(fdc_id);
CREATE INDEX IF NOT EXISTS idx_usda_food_nutrients_nutrient_id ON usda_food_nutrients(nutrient_id);
