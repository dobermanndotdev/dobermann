-- +goose Up
-- +goose StatementBegin
ALTER TABLE monitors ADD COLUMN is_paused BOOLEAN NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE monitors DROP COLUMN is_paused;
-- +goose StatementEnd
