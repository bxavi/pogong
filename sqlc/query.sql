-- name: GetAccounts :one
SELECT * FROM accounts
WHERE id = $1  LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY email;

-- name: CreateAccounts :one
INSERT INTO accounts (
	email, password
) VALUES (
	$1, $2
)
RETURNING *;

-- name: UpdateAccounts :one
UPDATE accounts
set email = $2,
password = $3
WHERE id = $1
RETURNING *;


-- name: DeleteAccounts :exec
DELETE FROM accounts
WHERE id = $1 ;

-- Add queries here:
