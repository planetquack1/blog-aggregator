-- name: CreateFeedFollows :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- -- name: GetUser :one
-- SELECT * FROM users
-- WHERE api_key = $1
-- LIMIT 1;