-- name: CreateFeedFollows :one
INSERT INTO  feed_follows (id, user_id, feed_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteFeedFollows :one
DELETE FROM feed_follows
WHERE id=$1
RETURNING *;

-- name: GetAllFeedFollows :many
SELECT *
FROM feed_follows;


