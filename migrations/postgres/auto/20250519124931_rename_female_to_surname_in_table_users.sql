-- +goose Up
-- Rename 'female' column to 'surname' in users table
ALTER TABLE users RENAME COLUMN female TO surname;

-- +goose Down
-- Revert the change by renaming back to 'female'
ALTER TABLE users RENAME COLUMN surname TO female;
