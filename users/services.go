package users

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaogiacometti/gostudies/pgstore"
	"golang.org/x/crypto/bcrypt"
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
var ErrInvalidCredentials = errors.New("invalid credentials")

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

func (us *UserService) AuthenticateUser(ctx context.Context, email, password string) (uuid.UUID, error) {
	user, err := us.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return uuid.Nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return uuid.Nil, ErrInvalidCredentials
		}

		return uuid.Nil, err
	}
	return user.ID, nil
}
