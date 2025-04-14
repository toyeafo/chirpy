-- +goose Up
create table refresh_tokens (
    token text primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    user_id uuid not null REFERENCES users(id) on DELETE CASCADE,
    expires_at timestamp not null,
    revoked_at timestamp null,
    CONSTRAINT check_expires_at_future CHECK (expires_at > NOW())
);

-- +goose Down
drop table refresh_tokens;