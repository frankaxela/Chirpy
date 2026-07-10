-- +goose Up
-- Add a hashed password column to the users table.
-- Existing users will have an unset/blank value.
ALTER TABLE users
    ADD COLUMN hashed_password TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE users
    DROP COLUMN hashed_password;
