-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS carts(
    "id" SERIAL NOT NULL PRIMARY KEY,
    "user_id" INT NOT NULL,
    "book_id" INT NOT NULL,
    "quantity" INT NOT NULL CHECK(quantity >= 1),
    "created_at" TIMESTAMPTZ DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE carts;
-- +goose StatementEnd
