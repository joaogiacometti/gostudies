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
SELECT id, question, answer, created_at, updated_at
FROM flashcards
WHERE id = $1 AND user_id = $2;