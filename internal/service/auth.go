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


func (a *AuthService) VerifyEmail(ctx context.Context, input types.VerifyUserInput) (types.LoginOutput, error) {
	email, err := a.cache.GetVerificationCode(ctx, input.Token)
    if err != nil {
		return types.LoginOutput{}, ErrInvalidCode
	}

	if email != input.Email {
		return types.LoginOutput{}, ErrInvalidCode
	}
	user, err := a.repo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return types.LoginOutput{}, ErrUserNotFound
	}

	err = a.repo.UpdateUserVerified(ctx, repository.UpdateUserVerifiedParams{
		ID: user.ID,
		IsVerified: true,
	})

	if err != nil {
		return types.LoginOutput{}, ErrUserNotFound
	}

	token, err := a.jwt.GenerateToken(user.ID)
	if err != nil {
		return types.LoginOutput{}, ErrGenerateToken
	}

	return types.LoginOutput{
		AccessToken: token,
	}, nil

}


func (a *AuthService) ResendVerificationCode(ctx context.Context, email string) error {
	user, err := a.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return ErrUserNotFound
	}

	userDetails, err := a.getUserDetails(ctx, user)
	if err != nil {
		return ErrUserNotFound
	}

	if user.IsVerified {
		return ErrUserAlreadyVerified
	}

	verificationCode := utils.Generate6DigitCode()

	err = a.cache.SetVerificationCode(ctx, verificationCode, email)
	if err != nil {
	return ErrSetVerification 
	}

	a.queue.EnqueueSendVerificationEmail(email, userDetails.FirstName, verificationCode, 0)
return nil
}


type getUserDetails struct {
	FirstName string
	LastName  string
	Email     string
}

func (a *AuthService) getUserDetails(ctx context.Context, user interface{}) (getUserDetails, error) {
	var details getUserDetails
	var email string
	// var userID *uuid.UUID

	switch u := user.(type) {

	case repository.GetUserByEmailRow:
		email = u.Email
		details = getUserDetails{
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     email,
		}

	case repository.GetUserByIDRow:
		email = u.Email
		details = getUserDetails{
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     email,
		}

	default:
		return getUserDetails{}, errors.New("unknown user type")
	}

	return details, nil
}




func (a *AuthService) Login(ctx context.Context, input types.LoginInput) (types.LoginOutput, error ) {
	user, err := a.repo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return types.LoginOutput{}, ErrInvalidCredential
	}

	if !user.IsVerified {
		return types.LoginOutput{}, ErrUserNotVerified
	}

	if err := utils.CheckPasswordHash(input.Password, user.PasswordHash); err != nil {
        return types.LoginOutput{}, ErrInvalidCredential
	}

	token, err := a.jwt.GenerateToken(user.ID)
	if err != nil {
		return types.LoginOutput{}, ErrGenerateToken
	}

	return types.LoginOutput{
		AccessToken: token,
	}, nil

}


func (a *AuthService) ForgotPassword(ctx context.Context, email string) error {
	user, err := a.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return ErrUserNotFound
	}

	userDetails, err := a.getUserDetails(ctx, user)
	if err != nil {
		return ErrUserNotFound
	}
	resetPasswordCode := utils.Generate6DigitCode()

	err = a.cache.SetResetPasswordCode(ctx, resetPasswordCode, email)
    if err != nil {
		return ErrSetVerification
	}
	a.queue.EnqueueSendResetPasswordEmail(email, userDetails.FirstName, resetPasswordCode, 0)
    return nil
}


func (a *AuthService) ResetPassword(ctx context.Context, input types.ResetPasswordInput) error {
	email, err := a.cache.GetResetPasswordCode(ctx, input.Token)
	if err != nil {
		return ErrInvalidCode
	}
	user, err := a.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return ErrUserNotFound
	}

	passwordHash, err := utils.HashPassword(input.Password)
    if err != nil {
		return err
	}

	err = a.repo.UpdateUserPassword(ctx, repository.UpdateUserPasswordParams{
		ID:	user.ID,
		PasswordHash: passwordHash, 
	})
	if err != nil {
		return ErrInvalidEmail
	}

	return nil 
}