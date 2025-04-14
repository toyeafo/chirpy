-- +goose Up
create table chirps (
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    body text not null,
    user_id uuid not null REFERENCES users(id) on DELETE CASCADE
);

-- +goose Down
drop table chirps;