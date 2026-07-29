-- name: CreateUser :one
INSERT INTO users (nickname)
VALUES ($1)
RETURNING *;
