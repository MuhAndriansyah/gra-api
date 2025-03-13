-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "order_number" VARCHAR(25) UNIQUE NOT NULL,
    "user_id" INT NOT NULL,
    "payment_status" VARCHAR(25) NOT NULL CHECK(payment_status IN ('Pending', 'Paid', 'Failed')),
    "payment_date" TIMESTAMPTZ,
    "payment_method" VARCHAR(100),
    "created_at" TIMESTAMPTZ DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id)

);


CREATE TABLE IF NOT EXISTS order_details (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "order_id" INT NOT NULL,
    "book_id" INT NOT NULL,
    "borrowing_date" TIMESTAMPTZ,
    "return_date" TIMESTAMPTZ,
    "created_at" TIMESTAMPTZ DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (book_id) REFERENCES books(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
DROP TABLE order_details;
-- +goose StatementEnd
