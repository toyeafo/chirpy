-- +goose Up
alter table users add hashed_password text not null default 'unset';

-- +goose Down
alter table users drop COLUMN hashed_password;