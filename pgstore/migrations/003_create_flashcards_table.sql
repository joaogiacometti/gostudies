-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS flashcards (
  id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  question TEXT NOT NULL,
  answer TEXT NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (user_id, question)
)
---- create above / drop below ----
DROP TABLE IF EXISTS flashcards;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
