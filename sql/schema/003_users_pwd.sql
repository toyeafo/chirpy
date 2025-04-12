-- +goose Up
alter table users add hashed_password text not null default 'unset';