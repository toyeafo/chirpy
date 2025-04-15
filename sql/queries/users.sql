-- name: CreateUser :one
insert into users (id, created_at, updated_at, email, hashed_password)
values (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2
)
returning *;

-- name: DeleteUsers :exec
delete from users;

-- name: RetrieveUserPwd :one
select * from users where email = $1;

-- name: RetrieveUserByID :one
select * from users where id = $1;

-- name: UpdateUsernamePassword :one
update users set email = $1, hashed_password = $2, updated_at = now() where id = $3 returning *;

-- name: UpgradeUser :one
update users set is_chirpy_red = $1 where id = $2 returning *;