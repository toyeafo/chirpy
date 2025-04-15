-- +goose Up
alter table users add column is_chirpy_red boolean not null default false;

-- +goose Down
alter table users drop COLUMN is_chirpy_red;