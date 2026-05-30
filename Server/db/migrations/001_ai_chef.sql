-- Migration: 001_ai_chef
-- Purpose: Add AI Chef support — structured recipes, ingredient tracking, interaction logging

-- 1. Add AI Chef columns to recipes table
ALTER TABLE recipes
    ADD COLUMN IF NOT EXISTS ai_generated BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS source_prompt TEXT,
    ADD COLUMN IF NOT EXISTS generated_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS estimated_minutes INT,
    ADD COLUMN IF NOT EXISTS difficulty VARCHAR(20) CHECK (difficulty IN ('easy', 'medium', 'hard') OR difficulty IS NULL),
    ADD COLUMN IF NOT EXISTS cuisine VARCHAR(50),
    ADD COLUMN IF NOT EXISTS tags TEXT[] DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS steps_json JSONB;

-- 2. Create recipe_ingredients table (structured ingredients with quantities/units)
CREATE TABLE IF NOT EXISTS recipe_ingredients (
    id SERIAL PRIMARY KEY,
    recipe_id INT REFERENCES recipes(id) ON DELETE CASCADE,
    ingredient_name TEXT NOT NULL,
    quantity NUMERIC(10,2),
    unit VARCHAR(20)
);

-- 3. Create ai_chef_interactions table (usage tracking & personalization signals)
CREATE TABLE IF NOT EXISTS ai_chef_interactions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    ingredients_submitted TEXT[],
    results_count INT DEFAULT 0,
    recipe_selected_id INT REFERENCES recipes(id) ON DELETE SET NULL,
    was_ai_generated BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 4. GIN index on tags for fast ingredient matching queries
CREATE INDEX IF NOT EXISTS idx_recipes_tags ON recipes USING GIN (tags);

-- 5. Index on ai_generated for filtering
CREATE INDEX IF NOT EXISTS idx_recipes_ai_generated ON recipes (ai_generated);

-- 6. Index on interactions for analytics queries
CREATE INDEX IF NOT EXISTS idx_ai_chef_interactions_user ON ai_chef_interactions (user_id);
CREATE INDEX IF NOT EXISTS idx_ai_chef_interactions_created ON ai_chef_interactions (created_at);
