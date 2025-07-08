-- +goose Up
-- +goose StatementBegin
ALTER TABLE readings ADD COLUMN public BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE readings DROP COLUMN public;
-- +goose StatementEnd
