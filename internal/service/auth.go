package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yusufaniki/muslim_tech/internal/cache"
	"github.com/yusufaniki/muslim_tech/internal/queue/tasks"
	repository "github.com/yusufaniki/muslim_tech/internal/repository"
	"github.com/yusufaniki/muslim_tech/internal/types"
	"github.com/yusufaniki/muslim_tech/internal/utils"
	"github.com/yusufaniki/muslim_tech/pkg/auth"
)


type tokenInfoResponse  struct {
    Email     string   `json:"email"`
	GivenName string   `json:"given_name"` 
	Family    string   `json:"family_name"` 
}

type AuthService struct {
	dbc  *pgxpool.Pool
	repo  repository.Queries
	queue *tasks.Queue
	cache cache.RedisCache
	jwt   auth.JWTManager
}


func NewAuthService( dbc *pgxpool.Pool, repo repository.Queries, queue tasks.Queue, cache cache.RedisCache, jwt auth.JWTManager) *AuthService {
	return &AuthService{dbc: dbc, repo: repo, queue: &queue, cache: cache, jwt: jwt}
}

func (a *AuthService) Register(ctx context.Context, input types.RegisterUserInput) (types.RegisterUserOutput, error){
	passwordHash, err := utils.HashPassword(input.Password)
	if err != nil {
		return types.RegisterUserOutput{}, err
	}

	existingUser, err := a.repo.GetUserByEmail(ctx, input.Email)
	if err == nil {
		return types.RegisterUserOutput{UserID: existingUser.ID.String()}, nil
	}
	tx, err := a.dbc.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return types.RegisterUserOutput{}, err
	}

	defer tx.Rollback(ctx)

	qtx := a.repo.WithTx(tx)

	user, err := qtx.CreateUser(ctx, repository.CreateUserParams{
		Email: input.Email,
		PasswordHash: passwordHash,
		FirstName: input.FirstName,
		LastName: input.LastName,
	})

	if err != nil {
		var pgErr  *pgconn.PgError
		if errors.As(err, &pgErr){
			if pgErr.Code == "23505" {
				return types.RegisterUserOutput{}, ErrEmailAlreadyExists
			}
		}
		return types.RegisterUserOutput{}, err
	}

	verificationCode := utils.Generate6DigitCode()
	a.queue.EnqueueSendVerificationEmail(input.Email, input.FirstName, verificationCode, 0)
    
	err = a.cache.SetVerificationCode(ctx, verificationCode, input.Email)
    if err != nil {
		fmt.Println("Error setting verification code in cache:", err)
		return types.RegisterUserOutput{}, ErrSetVerificationCode
	}


	err = tx.Commit(ctx)
	if err != nil {
		return types.RegisterUserOutput{}, err
	}

	return types.RegisterUserOutput{
		UserID: user.ID.String(),
		Email: user.Email,
		FirstName: user.FirstName,
		LastName: user.LastName,
	},
	nil
}

