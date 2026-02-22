-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
	gen_random_uuid(),
	NOW(),
	NOW(),
	$1, -- parameter passed by application
	$2
)
RETURNING *;

-- name: DeleteChirps :exec
TRUNCATE TABLE chirps;

-- name: Chirps :many
SELECT * from chirps order by created_at asc;


-- name: Chirp :one
SELECT * from chirps where id = $1;
