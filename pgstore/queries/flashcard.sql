-- name: CreateFlashcard :one
INSERT INTO flashcards (user_id, question, answer)
VALUES ($1, $2, $3)
RETURNING id;

-- name: IsDuplicateFlashcard :one
SELECT EXISTS (
    SELECT 1 FROM flashcards
    WHERE user_id = $1 AND question = $2
);