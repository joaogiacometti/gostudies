package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaogiacometti/gostudies/internal/store/pgstore"
)

type UserService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewUserService(pool *pgxpool.Pool) *UserService {
	return &UserService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

var ErrDuplicatedUsernameOrEmail = errors.New("username or email already exists")

func (us *UserService) CreateUser(ctx context.Context, username, email string, password_hash []byte) (uuid.UUID, error) {
	userId, err := us.queries.CreateUser(ctx, pgstore.CreateUserParams{
		UserName:     username,
		Email:        email,
		PasswordHash: password_hash,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && (pgErr.Code == "23505") {
			return uuid.Nil, ErrDuplicatedUsernameOrEmail
		}
		return uuid.Nil, err
	}

	return userId, nil
}
