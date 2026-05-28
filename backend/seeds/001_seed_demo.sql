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

INSERT INTO recipe_ingredients (id, recipe_id, fdc_id, quantity, unit)
VALUES
    ('rci_001', 'rcp_breakfast_oats', 2346396, 80, 'g'),
    ('rci_002', 'rcp_breakfast_oats', 746778, 200, 'ml'),
    ('rci_003', 'rcp_breakfast_oats', 1105314, 120, 'g'),
    ('rci_004', 'rcp_lunch_chicken_rice', 331960, 180, 'g'),
    ('rci_005', 'rcp_lunch_chicken_rice', 2512381, 170, 'g'),
    ('rci_006', 'rcp_lunch_chicken_rice', 747447, 100, 'g'),
    ('rci_007', 'rcp_lunch_chicken_rice', 748608, 10, 'ml'),
    ('rci_008', 'rcp_dinner_eggs_rice', 323604, 2, 'piece'),
    ('rci_009', 'rcp_dinner_eggs_rice', 2512381, 150, 'g'),
    ('rci_010', 'rcp_dinner_eggs_rice', 747447, 80, 'g'),
    ('rci_011', 'rcp_dinner_eggs_rice', 748608, 8, 'ml'),
    ('rci_012', 'rcp_snack_yogurt_apple', 330137, 200, 'g'),
    ('rci_013', 'rcp_snack_yogurt_apple', 1750341, 160, 'g')
ON CONFLICT(id) DO NOTHING;

INSERT INTO pantry_items (id, user_id, fdc_id, quantity, unit)
VALUES
    ('pnt_001', 'usr_demo_001', 2346396, 500, 'g'),
    ('pnt_002', 'usr_demo_001', 746778, 1500, 'ml'),
    ('pnt_003', 'usr_demo_001', 1105314, 600, 'g'),
    ('pnt_004', 'usr_demo_001', 331960, 1000, 'g'),
    ('pnt_005', 'usr_demo_001', 2512381, 2000, 'g'),
    ('pnt_006', 'usr_demo_001', 747447, 800, 'g'),
    ('pnt_007', 'usr_demo_001', 748608, 500, 'ml'),
    ('pnt_008', 'usr_demo_001', 323604, 12, 'piece'),
    ('pnt_009', 'usr_demo_001', 330137, 1000, 'g'),
    ('pnt_010', 'usr_demo_001', 1750341, 900, 'g')
ON CONFLICT(user_id, fdc_id, unit) DO UPDATE SET
    quantity = excluded.quantity,
    updated_at = CURRENT_TIMESTAMP;
