-- name: CreateUser :one
INSERT INTO users (nickname)
VALUES ($1)
RETURNING *;

-- name: GetUser :one
SELECT  *
FROM    users
WHERE   nickname = $1;
