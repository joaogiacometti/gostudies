-- Write your migrate up statements here
ALTER TABLE flashcards
ADD COLUMN IF NOT EXISTS last_reviewed_at TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS next_review_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
ADD COLUMN IF NOT EXISTS success_count INT DEFAULT 0;
---- create above / drop below ----
ALTER TABLE flashcards
DROP COLUMN IF EXISTS last_reviewed_at,
DROP COLUMN IF EXISTS next_review_at,
DROP COLUMN IF EXISTS success_count;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
