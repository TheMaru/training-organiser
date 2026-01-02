-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN platform_role TEXT NOT NULL DEFAULT 'member';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN platform_role;
-- +goose StatementEnd
