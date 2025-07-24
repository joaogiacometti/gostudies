package flashcards

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaogiacometti/gostudies/pgstore"
)

type FlashcardService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewFlashcardService(pool *pgxpool.Pool) *FlashcardService {
	return &FlashcardService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (fs *FlashcardService) CreateFlashcard(ctx context.Context, userID uuid.UUID, question, answer string) (uuid.UUID, error) {
	flashcardID, err := fs.queries.CreateFlashcard(ctx, pgstore.CreateFlashcardParams{
		UserID:   userID,
		Question: question,
		Answer:   answer,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return flashcardID, nil
}
