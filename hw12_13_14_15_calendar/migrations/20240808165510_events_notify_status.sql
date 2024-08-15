-- +goose Up
-- +goose StatementBegin
ALTER TABLE events ADD notify_status int DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events DROP COLUMN notify_status;
-- +goose StatementEnd
