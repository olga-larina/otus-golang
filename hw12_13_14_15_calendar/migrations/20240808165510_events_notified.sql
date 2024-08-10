-- +goose Up
-- +goose StatementBegin
ALTER TABLE events ADD notified boolean DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events DROP COLUMN notified;
-- +goose StatementEnd
