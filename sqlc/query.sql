-- name: GetAccount :one
SELECT * FROM account
WHERE id = $1  LIMIT 1;

-- name: ListAccount :many
SELECT * FROM account
ORDER BY id
LIMIT sqlc.narg('limit')
OFFSET sqlc.narg('offset');


-- name: CreateAccount :one
INSERT INTO account (
	email, password
) VALUES (
	$1, $2
)
RETURNING *;

-- name: UpdateAccount :one
UPDATE account
set email = $2,
password = $3
WHERE id = $1
RETURNING *;


-- name: DeleteAccount :exec
DELETE FROM account
WHERE id = $1 ;

-- Add queries here:

-- name: GetAccountWithEmail :one
SELECT * FROM account
WHERE email = $1;


