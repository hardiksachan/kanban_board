package jwt

import (
	"github.com/golang-jwt/jwt"
	"github.com/hardiksachan/kanban_board/backend/internal/users"
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
	"time"
)

type Claims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

type AccessTokenStore struct {
	jwtKey         []byte
	expiryDuration time.Duration
}

func NewAccessTokenStore(jwtKey string, expiryDuration time.Duration) *AccessTokenStore {
	return &AccessTokenStore{[]byte(jwtKey), expiryDuration}
}

func (s *AccessTokenStore) Create(userId string) (*domain.AccessToken, error) {
	return &domain.AccessToken{
		UserID:    userId,
		ExpiresAt: time.Now().Add(s.expiryDuration),
	}, nil
}

// VerifyAndDecode Returns nil if access token is invalid
func (s *AccessTokenStore) VerifyAndDecode(encodedToken string) (*domain.AccessToken, error) {
	op := "jwt.AccessTokenStore.Verify"

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(encodedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtKey, nil
	})
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	if !token.Valid {
		return nil, nil
	}

	return &domain.AccessToken{
		UserID:    claims.UserId,
		ExpiresAt: time.Unix(claims.ExpiresAt, 0),
	}, nil
}

func (s *AccessTokenStore) Encode(token *domain.AccessToken) (string, error) {
	op := "jwt.AccessTokenStore.Encode"

	claims := &Claims{
		UserId: token.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: token.ExpiresAt.Unix(),
		},
	}

	signedToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtKey)
	if err != nil {
		return "", &users.Error{Op: op, Err: err}
	}

	return signedToken, nil
}
