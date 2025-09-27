-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE IF NOT EXISTS organizations(
    uuid uuid not null unique default gen_random_uuid() primary key,
    name varchar(255),
    email varchar(255) unique
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists organizations;
-- +goose StatementEnd
