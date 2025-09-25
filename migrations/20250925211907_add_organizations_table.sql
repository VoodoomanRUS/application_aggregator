-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS organizations(
    ID int PRIMARY key,
    UUID varchar(255) unique,
    name varchar(255) unique
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- SELECT 'down SQL query';
-- +goose StatementEnd
