package reviews

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaogiacometti/gostudies/pgstore"
)

type ReviewService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewReviewService(pool *pgxpool.Pool) *ReviewService {
	return &ReviewService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (fs *ReviewService) ReviewFlashcard(ctx context.Context, flashcardID uuid.UUID, userID uuid.UUID, nextReview time.Time, successCount int32, remembered bool) error {
	pgSuccessCount := pgtype.Int4{Int32: successCount, Valid: true}
	pgNextReview := pgtype.Timestamptz{Time: nextReview, Valid: true}

	err := fs.queries.ReviewFlashcard(ctx, pgstore.ReviewFlashcardParams{
		ID:           flashcardID,
		UserID:       userID,
		SuccessCount: pgSuccessCount,
		NextReviewAt: pgNextReview,
	})
	if err != nil {
		return err
	}

	fs.queries.AddReview(ctx, pgstore.AddReviewParams{
		FlashcardID: flashcardID,
		UserID:      userID,
		Remembered:  remembered,
	})

	return nil
}
