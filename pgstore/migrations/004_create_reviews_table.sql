-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS reviews(
    ID UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    flashcard_id UUID NOT NULL REFERENCES flashcards(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    remembered BOOLEAN NOT NULL,
    reviewed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
---- create above / drop below ----
DROP TABLE IF EXISTS reviews;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
