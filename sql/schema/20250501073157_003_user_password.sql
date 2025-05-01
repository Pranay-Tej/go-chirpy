-- +goose Up
-- +goose StatementBegin

ALTER TABLE users
ADD COLUMN hashed_password TEXT NOT NULL
DEFAULT 'unset';

-- ALTER TABLE users
-- ADD COLUMN hashed_password TEXT;

-- -- set default for existing users
-- UPDATE users SET hashed_password = 'unset';

-- -- add NOT NULL constraint
-- ALTER TABLE users
-- ALTER COLUMN hashed_password SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN hashed_password;
-- +goose StatementEnd
