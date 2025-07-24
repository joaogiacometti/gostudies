-- name: AddReview :exec
INSERT INTO reviews (flashcard_id, user_id, remembered, reviewed_at)
VALUES ($1, $2, $3, NOW());