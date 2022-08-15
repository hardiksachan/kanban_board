package repository

import (
	"github.com/golang-jwt/jwt"
	"github.com/hardiksachan/kanban_board/backend/internal/users"
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
	"time"
)

type Claims struct {
	domain.AccessClaims
	jwt.StandardClaims
}

type JWTAccessProvider struct {
	jwtKey         []byte
	expiryDuration time.Duration
}

func NewAccessTokenStore(jwtKey string, expiryDuration time.Duration) *JWTAccessProvider {
	return &JWTAccessProvider{[]byte(jwtKey), expiryDuration}
}

func (p *JWTAccessProvider) Create(credential *domain.AccessClaims) (*domain.AccessToken, error) {
	op := "JWTAccessProvider.Create"

	claims := &Claims{
		AccessClaims: domain.AccessClaims{
			UserID: credential.UserID,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(p.expiryDuration).Unix(),
		},
	}

	signedToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(p.jwtKey)
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	return (*domain.AccessToken)(&signedToken), nil
}

func (p *JWTAccessProvider) Verify(accessToken *domain.AccessToken) (*domain.AccessClaims, error) {
	op := "JWTAccessProvider.Verify"

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(string(*accessToken), claims, func(token *jwt.Token) (interface{}, error) {
		return p.jwtKey, nil
	})
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	if !token.Valid {
		return nil, nil
	}

	return &claims.AccessClaims, nil
}
