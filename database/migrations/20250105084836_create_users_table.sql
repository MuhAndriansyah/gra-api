-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
 "id" SERIAL PRIMARY KEY,
 "name" VARCHAR(255) NOT NULL,
 "email" VARCHAR(255) NOT NULL UNIQUE,
 "password" VARCHAR(255) NOT NULL,
 "photo" VARCHAR(255),
 "email_verify_token" VARCHAR(8),
 "email_token_expired_at" TIMESTAMPTZ,
 "verified_at" TIMESTAMPTZ,
 "created_at" TIMESTAMPTZ DEFAULT NOW(),
 "updated_at" TIMESTAMPTZ DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user;
-- +goose StatementEnd
