-- name: CreateChirp :one
insert into chirps (id, created_at, updated_at, body, user_id)
values (
    gen_random_uuid(),
    now(),
    now(),
    $1, 
    $2
)
returning *;

-- name: DeleteChirp :exec
delete from chirps where user_id = $1;

-- name: GetChirps :many
select * from chirps order by created_at ASC;

-- name: GetSingleChirp :one
select * from chirps where id = $1;