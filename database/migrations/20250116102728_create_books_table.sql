-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS books (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "title" VARCHAR(255) NOT NULL,
    "slug" VARCHAR(255) NOT NULL,
    "author_id" INT NOT NULL,
    "publisher_id" INT NOT NULL,
    "publish_year" SMALLINT NOT NULL,
    "total_page" INT NOT NULL,
    "description" TEXT,
    "sku" VARCHAR(255) NOT NULL,
    "isbn" VARCHAR(255) NOT NULL UNIQUE,
    "price" NUMERIC(12, 2),
    "created_at" TIMESTAMPTZ DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (author_id) REFERENCES authors(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    FOREIGN KEY (publisher_id) REFERENCES publishers(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE books;
-- +goose StatementEnd
