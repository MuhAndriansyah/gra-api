-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS permissions(
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "display_name" VARCHAR(255) NOT NULL,
    "description" TEXT,
    "created_at" TIMESTAMPTZ DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE("name")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE permissions;
-- +goose StatementEnd
