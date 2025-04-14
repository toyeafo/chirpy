-- name: CreateRefreshToken :one
insert into refresh_tokens (token, created_at, updated_at, user_id, expires_at)
values (
    $1,
    now(),
    now(),
    $2,
    $3
)
returning *;

-- name: GetUserRefreshToken :one
select users.* from users 
join refresh_tokens on users.id = refresh_tokens.user_id 
where refresh_tokens.token = $1 
AND revoked_at IS NULL 
AND expires_at > NOW();

-- name: UpdateRefreshToken :exec
update refresh_tokens set updated_at = now(), revoked_at = now() where token = $1 returning *;