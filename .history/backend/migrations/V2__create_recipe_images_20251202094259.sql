-- Migration to create recipe_images table for multiple images per recipe
CREATE TABLE recipe_images (
    id SERIAL PRIMARY KEY,
    recipe_id INTEGER REFERENCES recipes(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    is_featured BOOLEAN DEFAULT FALSE
);