-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS book_category(
    "id" SERIAL NOT NULL PRIMARY KEY,
    "book_id" INT,
    "category_id" INT,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE book_category;
-- +goose StatementEnd
