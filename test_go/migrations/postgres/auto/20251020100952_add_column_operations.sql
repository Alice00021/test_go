-- +goose Up
-- +goose StatementBegin
alter table operations
    add column IF NOT EXISTS commands JSONB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table operations
    drop column IF EXISTS commands;
-- +goose StatementEnd
