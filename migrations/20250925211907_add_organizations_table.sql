-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS organizations(
    ID serial primary key,
    uuid uuid not null unique default gen_random_uuid(),
    name varchar(255) unique
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists organizations;
-- +goose StatementEnd
