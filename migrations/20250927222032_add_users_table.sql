-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE IF NOT EXISTS users
(
    uuid  uuid not null unique default gen_random_uuid() primary key,
    name  varchar(255)
);

-- +goose Down
drop table if exists users;
-- This is a core migration. DO NOT DROP TABLES OR DELETE DATA.
-- Down migration is intentionally left empty to prevent accidental data loss.