PRAGMA foreign_keys = ON;

INSERT INTO users (id, email, password_hash, display_name)
VALUES
    ('usr_demo_001', 'demo@pantrypal.local', '$2a$10$demo.hash.placeholder', 'Demo User')
ON CONFLICT(id) DO NOTHING;

INSERT INTO user_body_metrics (id, user_id, height_cm, weight_kg, age, sex, activity_level, goal)
VALUES
    ('ubm_demo_001', 'usr_demo_001', 175, 72, 30, 'male', 'moderate', 'maintain')
ON CONFLICT(user_id) DO NOTHING;

INSERT INTO user_preferences (id, user_id, diet_type, allergies_json, dislikes_json, likes_json, daily_calorie_target, notes)
VALUES
    (
        'upr_demo_001',
        'usr_demo_001',
        'omnivore',
        '["peanuts"]',
        '["mushrooms"]',
        '["rice", "chicken", "yogurt"]',
        2200,
        'Demo profile for deterministic local testing'
    )
ON CONFLICT(user_id) DO NOTHING;

INSERT INTO budgets (id, user_id, month, currency, amount_cents)
VALUES
    ('bdg_demo_2026_05', 'usr_demo_001', '2026-05', 'USD', 45000)
ON CONFLICT(user_id, month) DO NOTHING;

INSERT INTO ingredients (
    id,
    name,
    canonical_unit,
    kcal_per_unit,
    protein_g_per_unit,
    carbs_g_per_unit,
    fat_g_per_unit,
    estimated_cost_cents_per_unit
)
VALUES
    ('ing_oats', 'oats', 'g', 3.89, 0.17, 0.66, 0.07, 1),
    ('ing_milk', 'milk', 'ml', 0.61, 0.033, 0.048, 0.033, 1),
    ('ing_banana', 'banana', 'g', 0.89, 0.011, 0.23, 0.003, 1),
    ('ing_chicken_breast', 'chicken breast', 'g', 1.65, 0.31, 0.0, 0.036, 2),
    ('ing_rice', 'rice', 'g', 1.30, 0.027, 0.28, 0.003, 1),
    ('ing_broccoli', 'broccoli', 'g', 0.35, 0.024, 0.07, 0.004, 1),
    ('ing_olive_oil', 'olive oil', 'ml', 8.84, 0.0, 0.0, 1.0, 2),
    ('ing_egg', 'egg', 'piece', 72, 6.3, 0.4, 4.8, 30),
    ('ing_greek_yogurt', 'greek yogurt', 'g', 0.59, 0.10, 0.036, 0.004, 2),
    ('ing_apple', 'apple', 'g', 0.52, 0.003, 0.14, 0.002, 1)
ON CONFLICT(id) DO NOTHING;

INSERT INTO recipes (
    id,
    name,
    meal_type,
    servings,
    instructions,
    total_kcal,
    total_protein_g,
    total_carbs_g,
    total_fat_g,
    estimated_cost_cents
)
VALUES
    (
        'rcp_breakfast_oats',
        'Banana Oat Bowl',
        'breakfast',
        1,
        'Cook oats in milk, top with sliced banana.',
        430,
        17,
        69,
        10,
        180
    ),
    (
        'rcp_lunch_chicken_rice',
        'Chicken Rice Plate',
        'lunch',
        1,
        'Cook rice, grill chicken, steam broccoli, finish with olive oil.',
        640,
        52,
        57,
        20,
        420
    ),
    (
        'rcp_dinner_eggs_rice',
        'Egg Fried Rice Lite',
        'dinner',
        1,
        'Scramble eggs, stir with cooked rice and broccoli in olive oil.',
        560,
        24,
        63,
        22,
        260
    ),
    (
        'rcp_snack_yogurt_apple',
        'Yogurt Apple Snack',
        'snacks',
        1,
        'Slice apple and serve with greek yogurt.',
        210,
        16,
        30,
        1,
        170
    )
ON CONFLICT(id) DO NOTHING;

INSERT INTO recipe_ingredients (id, recipe_id, ingredient_id, quantity, unit)
VALUES
    ('rci_001', 'rcp_breakfast_oats', 'ing_oats', 80, 'g'),
    ('rci_002', 'rcp_breakfast_oats', 'ing_milk', 200, 'ml'),
    ('rci_003', 'rcp_breakfast_oats', 'ing_banana', 120, 'g'),
    ('rci_004', 'rcp_lunch_chicken_rice', 'ing_chicken_breast', 180, 'g'),
    ('rci_005', 'rcp_lunch_chicken_rice', 'ing_rice', 170, 'g'),
    ('rci_006', 'rcp_lunch_chicken_rice', 'ing_broccoli', 100, 'g'),
    ('rci_007', 'rcp_lunch_chicken_rice', 'ing_olive_oil', 10, 'ml'),
    ('rci_008', 'rcp_dinner_eggs_rice', 'ing_egg', 2, 'piece'),
    ('rci_009', 'rcp_dinner_eggs_rice', 'ing_rice', 150, 'g'),
    ('rci_010', 'rcp_dinner_eggs_rice', 'ing_broccoli', 80, 'g'),
    ('rci_011', 'rcp_dinner_eggs_rice', 'ing_olive_oil', 8, 'ml'),
    ('rci_012', 'rcp_snack_yogurt_apple', 'ing_greek_yogurt', 200, 'g'),
    ('rci_013', 'rcp_snack_yogurt_apple', 'ing_apple', 160, 'g')
ON CONFLICT(id) DO NOTHING;

INSERT INTO pantry_items (id, user_id, ingredient_id, quantity, unit)
VALUES
    ('pnt_001', 'usr_demo_001', 'ing_oats', 500, 'g'),
    ('pnt_002', 'usr_demo_001', 'ing_milk', 1500, 'ml'),
    ('pnt_003', 'usr_demo_001', 'ing_banana', 600, 'g'),
    ('pnt_004', 'usr_demo_001', 'ing_chicken_breast', 1000, 'g'),
    ('pnt_005', 'usr_demo_001', 'ing_rice', 2000, 'g'),
    ('pnt_006', 'usr_demo_001', 'ing_broccoli', 800, 'g'),
    ('pnt_007', 'usr_demo_001', 'ing_olive_oil', 500, 'ml'),
    ('pnt_008', 'usr_demo_001', 'ing_egg', 12, 'piece'),
    ('pnt_009', 'usr_demo_001', 'ing_greek_yogurt', 1000, 'g'),
    ('pnt_010', 'usr_demo_001', 'ing_apple', 900, 'g')
ON CONFLICT(user_id, ingredient_id, unit) DO UPDATE SET
    quantity = excluded.quantity,
    updated_at = CURRENT_TIMESTAMP;
