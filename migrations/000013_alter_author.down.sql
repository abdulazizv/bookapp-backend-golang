ALTER TABLE author
    DROP COLUMN IF EXISTS about_text,
    DROP COLUMN IF EXISTS creativity,
    DROP COLUMN IF EXISTS died_year;