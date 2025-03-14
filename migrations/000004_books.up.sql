CREATE TABLE IF NOT EXISTS books (
    id SERIAL NOT NULL PRIMARY KEY,
    download_url VARCHAR,
    audio_url VARCHAR,
    category_id INT REFERENCES category(id) ON DELETE CASCADE,
    subcategory_id INT REFERENCES sub_category(id) ON DELETE CASCADE,
    author_id INT REFERENCES author(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
)