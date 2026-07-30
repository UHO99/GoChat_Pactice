-- name: CreateRoom :one
INSERT INTO rooms(name)
VALUES ($1)
RETURNING *;

-- name: GetRoomByName :one
SELECT * FROM rooms WHERE name = $1;

-- name: ListRooms :many
SELECT * FROM rooms ORDER BY name;
