PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    display_name TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_body_metrics (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL UNIQUE,
    height_cm REAL,
    weight_kg REAL,
    age INTEGER,
    sex TEXT,
    activity_level TEXT,
    goal TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CHECK (sex IN ('female', 'male', 'other', 'prefer_not_to_say') OR sex IS NULL)
);

CREATE TABLE IF NOT EXISTS user_preferences (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL UNIQUE,
    diet_type TEXT NOT NULL DEFAULT 'omnivore',
    allergies_json TEXT NOT NULL DEFAULT '[]',
    dislikes_json TEXT NOT NULL DEFAULT '[]',
    likes_json TEXT NOT NULL DEFAULT '[]',
    daily_calorie_target INTEGER,
    notes TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS budgets (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    month TEXT NOT NULL,
    currency TEXT NOT NULL DEFAULT 'USD',
    amount_cents INTEGER NOT NULL CHECK (amount_cents >= 0),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (user_id, month)
);

CREATE TABLE IF NOT EXISTS purchases (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    budget_id TEXT,
    fdc_id INTEGER,
    description TEXT,
    quantity REAL,
    unit TEXT,
    cost_cents INTEGER NOT NULL CHECK (cost_cents >= 0),
    purchased_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (budget_id) REFERENCES budgets(id) ON DELETE SET NULL,
    FOREIGN KEY (fdc_id) REFERENCES usda_foods(fdc_id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS recipes (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    meal_type TEXT NOT NULL,
    servings INTEGER NOT NULL DEFAULT 1 CHECK (servings > 0),
    instructions TEXT,
    total_kcal REAL NOT NULL DEFAULT 0,
    total_protein_g REAL NOT NULL DEFAULT 0,
    total_carbs_g REAL NOT NULL DEFAULT 0,
    total_fat_g REAL NOT NULL DEFAULT 0,
    estimated_cost_cents INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (meal_type IN ('breakfast', 'lunch', 'dinner', 'snacks'))
);

CREATE TABLE IF NOT EXISTS recipe_ingredients (
    id TEXT PRIMARY KEY,
    recipe_id TEXT NOT NULL,
    fdc_id INTEGER NOT NULL,
    quantity REAL NOT NULL CHECK (quantity >= 0),
    unit TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
    FOREIGN KEY (fdc_id) REFERENCES usda_foods(fdc_id) ON DELETE RESTRICT,
    UNIQUE (recipe_id, fdc_id, unit)
);

CREATE TABLE IF NOT EXISTS pantry_items (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    fdc_id INTEGER NOT NULL,
    quantity REAL NOT NULL CHECK (quantity >= 0),
    unit TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (fdc_id) REFERENCES usda_foods(fdc_id) ON DELETE RESTRICT,
    UNIQUE (user_id, fdc_id, unit)
);

CREATE TABLE IF NOT EXISTS meal_plans (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    period_type TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status TEXT NOT NULL DEFAULT 'accepted',
    source TEXT NOT NULL DEFAULT 'ai',
    proposal_version INTEGER NOT NULL DEFAULT 1,
    ai_cost_cents_total INTEGER NOT NULL DEFAULT 0 CHECK (ai_cost_cents_total >= 0),
    notes TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CHECK (period_type IN ('day', 'week', 'month')),
    CHECK (status IN ('proposal', 'accepted', 'declined', 'archived')),
    CHECK (source IN ('ai', 'fallback', 'manual')),
    CHECK (start_date <= end_date)
);

CREATE TABLE IF NOT EXISTS plan_meals (
    id TEXT PRIMARY KEY,
    meal_plan_id TEXT NOT NULL,
    recipe_id TEXT,
    scheduled_date DATE NOT NULL,
    meal_section TEXT NOT NULL,
    recipe_name TEXT NOT NULL,
    servings REAL NOT NULL DEFAULT 1 CHECK (servings > 0),
    kcal REAL NOT NULL DEFAULT 0,
    protein_g REAL NOT NULL DEFAULT 0,
    carbs_g REAL NOT NULL DEFAULT 0,
    fat_g REAL NOT NULL DEFAULT 0,
    estimated_cost_cents INTEGER NOT NULL DEFAULT 0,
    is_consumed INTEGER NOT NULL DEFAULT 0 CHECK (is_consumed IN (0, 1)),
    consumed_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (meal_plan_id) REFERENCES meal_plans(id) ON DELETE CASCADE,
    FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE SET NULL,
    CHECK (meal_section IN ('breakfast', 'lunch', 'dinner', 'snacks')),
    UNIQUE (meal_plan_id, scheduled_date, meal_section)
);

CREATE TABLE IF NOT EXISTS consumption_log (
    id TEXT PRIMARY KEY,
    plan_meal_id TEXT,
    user_id TEXT NOT NULL,
    fdc_id INTEGER,
    pantry_item_id TEXT,
    quantity_deducted REAL NOT NULL CHECK (quantity_deducted >= 0),
    unit TEXT NOT NULL,
    before_quantity REAL NOT NULL CHECK (before_quantity >= 0),
    after_quantity REAL NOT NULL CHECK (after_quantity >= 0),
    warning TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (plan_meal_id) REFERENCES plan_meals(id) ON DELETE SET NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (fdc_id) REFERENCES usda_foods(fdc_id) ON DELETE SET NULL,
    FOREIGN KEY (pantry_item_id) REFERENCES pantry_items(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS chat_messages (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    meal_plan_id TEXT,
    role TEXT NOT NULL,
    action TEXT,
    content TEXT NOT NULL,
    metadata_json TEXT NOT NULL DEFAULT '{}',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (meal_plan_id) REFERENCES meal_plans(id) ON DELETE SET NULL,
    CHECK (role IN ('system', 'user', 'assistant', 'tool'))
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_pantry_items_user_id ON pantry_items(user_id);
CREATE INDEX IF NOT EXISTS idx_pantry_items_fdc_id ON pantry_items(fdc_id);
CREATE INDEX IF NOT EXISTS idx_meal_plans_user_period ON meal_plans(user_id, period_type, start_date);
CREATE INDEX IF NOT EXISTS idx_plan_meals_meal_plan_id ON plan_meals(meal_plan_id);
CREATE INDEX IF NOT EXISTS idx_plan_meals_schedule ON plan_meals(scheduled_date, meal_section);
CREATE INDEX IF NOT EXISTS idx_consumption_log_user_created_at ON consumption_log(user_id, created_at);
CREATE INDEX IF NOT EXISTS idx_chat_messages_user_created_at ON chat_messages(user_id, created_at);
