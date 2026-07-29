-- name: InsertMessage :one
INSERT INTO messages (room_id, user_id, content)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListRecentMessages :many
SELECT m.id, m.room_id, m.user_id, m.content, m.created_at, u.nickname
FROM messages m
JOIN users u ON u.id = m.user_id
WHERE m.room_id = $1
ORDER BY m.created_at DESC
LIMIT $2;

