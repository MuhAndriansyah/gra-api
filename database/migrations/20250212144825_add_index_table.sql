-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_book_title ON books(title);
CREATE INDEX idx_books_author ON books(author_id);
CREATE INDEX idx_books_publisher ON books(publisher_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_book_title ON books;
DROP INDEX idx_books_author ON books;
DROP INDEX idx_books_publisher ON books;
-- +goose StatementEnd
