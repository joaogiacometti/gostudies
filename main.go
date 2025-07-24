package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaogiacometti/gostudies/api"
	"github.com/joaogiacometti/gostudies/flashcards"
	"github.com/joaogiacometti/gostudies/reviews"
	"github.com/joaogiacometti/gostudies/users"
	"github.com/joho/godotenv"
)

func init() {
	gob.Register(uuid.UUID{})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	sessionManager := api.NewSessionManager(pool)

	userHandlers := users.NewUserHandler(
		users.NewUserService(pool),
		sessionManager,
	)
	flashcardHandlers := flashcards.NewFlashcardHandler(
		flashcards.NewFlashcardService(pool),
		sessionManager,
	)
	reviewHandlers := reviews.NewReviewHandler(
		reviews.NewReviewService(pool),
		flashcards.NewFlashcardService(pool),
		sessionManager,
	)

	app := api.API{
		Router:            chi.NewMux(),
		UserHandlers:      userHandlers,
		FlashcardHandlers: flashcardHandlers,
		ReviewHandlers:    reviewHandlers,
		Sessions:          sessionManager,
	}

	app.BindRoutes()

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", app.Router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
