-- +goose Up
-- +goose StatementBegin
CREATE INDEX events_user_id_idx ON events (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS events_user_id_idx;
-- +goose StatementEnd
