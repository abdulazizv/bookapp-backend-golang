CREATE TABLE IF NOT EXISTS book_likes(
    id SERIAL PRIMARY KEY,
    user_id INT,
    book_id INT,
    CONSTRAINT unique_like UNIQUE (user_id, book_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);