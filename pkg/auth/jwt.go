package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/yusufaniki/muslim_tech/internal/repository"
)


type JWTManager struct {
	secretKey      string
	tokenDuration  time.Duration
	repo           repository.Queries
}



type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}


func NewJWTManager(secretKey string, tokenDuration time.Duration, repo repository.Queries) *JWTManager {
	return &JWTManager{
		secretKey: secretKey,
		tokenDuration: tokenDuration,
		repo: repo,
	}
}



func (j *JWTManager) GenerateToken(userId uuid.UUID) (string, error) {
	claims := &Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer: "muslim_tech_app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}



func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrInvalidKeyType
	}
	return claims, nil
}