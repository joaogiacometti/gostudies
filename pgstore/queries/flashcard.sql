-- name: CreateFlashcard :one
INSERT INTO flashcards (user_id, question, answer)
VALUES ($1, $2, $3)
RETURNING id;

-- name: IsDuplicateFlashcard :one
SELECT EXISTS (
    SELECT 1 FROM flashcards
    WHERE user_id = $1 AND question = $2
);

-- name: GetFlashcards :many
SELECT id, question, answer, created_at, updated_at
FROM flashcards
WHERE user_id = $1;

-- name: GetFlashcardByID :one
SELECT id, question, answer, created_at, updated_at, success_count
FROM flashcards
WHERE id = $1 AND user_id = $2;

-- name: GetNextFlashcardToReview :one
SELECT id, question, answer
FROM flashcards
WHERE user_id = $1 and next_review_at <= NOW()
ORDER BY next_review_at ASC
LIMIT 1;

-- name: ReviewFlashcard :exec  
UPDATE flashcards
SET success_count = $1,
    next_review_at = $2,
    last_reviewed_at = NOW(),
    updated_at = NOW()
WHERE id = $3 AND user_id = $4;
