-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url , user_id)
VALUES(
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds
ORDER BY name;

-- name: GetFeed :one
SELECT * FROM feeds WHERE id = $1;

-- name: GetNextFeedsToFetch :many
SELECT * 
FROM feeds
WHERE user_id =$1
ORDER BY last_fetched_at IS NULL DESC, last_fetched_at ASC
LIMIT $2;

-- name: DeleteFeed :exec
DELETE FROM feeds WHERE id=$1 AND user_id =$2;

-- name: MarkFeedFetched :exec
UPDATE feeds
 set updated_at = $2,
  last_fetched_at = $3
WHERE id = $1;

