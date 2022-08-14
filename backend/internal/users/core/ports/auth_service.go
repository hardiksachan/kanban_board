package ports

import (
	"fmt"
	"github.com/hardiksachan/kanban_board/backend/internal/users"
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	userStore         UserStore
	refreshTokenStore RefreshTokenStore
	accessTokenStore  AccessTokenStore
}

func NewAuthService(userStore UserStore, accessTokenStore AccessTokenStore, refreshTokenStore RefreshTokenStore) *AuthService {
	return &AuthService{userStore, refreshTokenStore, accessTokenStore}
}

// SignUp returns User after adding it to store
// Otherwise, returns ECONFLICT if user with same email exists
// returns wrapped error if fails
func (a *AuthService) SignUp(user *domain.User) (*domain.User, error) {
	op := "ports.AuthService.SignUp"

	exists, err := a.userStore.CheckByEmail(user.Email)
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}
	if exists {
		return nil, &users.Error{Op: op, Code: users.ECONFLICT, Message: fmt.Sprintf("user with email (%s) exists", user.Email)}
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return nil, &users.Error{Op: op, Code: users.EINTERNAL}
	}
	user.Password = hashedPassword

	storedUser, err := a.userStore.Insert(user)
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	return storedUser, nil
}

// LogIn creates and returns a AccessToken and RefreshToken if credentials are valid
// Otherwise, returns ENOTFOUND if email is incorrect
// returns ECONFLICT if passwords do not match
func (a *AuthService) LogIn(email, password string) (encodedAccessToken, encodedRefreshToken string, err error) {
	op := "ports.AuthService.Login"
	msg := fmt.Sprintf("email(%s) or password incorrect", email)

	storedUser, err := a.userStore.FindByEmail(email)
	if err != nil {
		return "", "", &users.Error{Op: op, Message: msg, Err: err}
	}

	_, err = VerifyPassword(password, storedUser.Password)
	if err != nil {
		return "", "", &users.Error{Op: op, Code: users.ECONFLICT, Message: msg, Err: err}
	}

	refreshToken, err := a.refreshTokenStore.Create(storedUser.ID)
	if err != nil {
		return "", "", &users.Error{Op: op, Err: err}
	}
	encodedRefreshToken, err = a.refreshTokenStore.Encode(refreshToken)
	if err != nil {
		return "", "", &users.Error{Op: op, Err: err}
	}

	accessToken, err := a.accessTokenStore.Create(storedUser.ID)
	if err != nil {
		return "", "", &users.Error{Op: op, Err: err}
	}

	encodedAccessToken, err = a.accessTokenStore.Encode(accessToken)
	if err != nil {
		return "", "", &users.Error{Op: op, Err: err}
	}

	return encodedAccessToken, encodedRefreshToken, nil
}

// LogOut deletes a RefreshToken
// Otherwise, returns error
func (a *AuthService) LogOut(encodedRefreshToken string) error {
	op := "ports.AuthService.LogOut"

	err := a.refreshTokenStore.Delete(encodedRefreshToken)
	if err != nil {
		return &users.Error{Op: op, Err: err}
	}
	return nil
}

func (a *AuthService) DecodeAccessToken(encodedAccessToken string) (*domain.AccessToken, error) {
	op := "ports.AuthService.DecodeAccessToken"

	accessToken, err := a.accessTokenStore.VerifyAndDecode(encodedAccessToken)
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}
	if accessToken == nil {
		return nil, &users.Error{Op: op, Message: "invalid access token", Code: users.EINVALID}
	}

	if accessToken.ExpiresAt.Unix() < time.Now().Unix() {
		return nil, &users.Error{Op: op, Message: "access token expired", Code: users.EEXPIRED}
	}

	return accessToken, nil
}

func (a *AuthService) RegenerateAccessToken(encodedRefreshToken string) (encodedAccessToken string, err error) {
	op := "ports.AuthService.RegenerateAccessToken"

	refreshToken, err := a.refreshTokenStore.VerifyAndDecode(encodedRefreshToken)
	if err != nil {
		return "", &users.Error{Op: op, Err: err}
	}
	if refreshToken == nil {
		return "", &users.Error{Op: op, Code: users.EINVALID, Err: err}
	}

	accessToken, err := a.accessTokenStore.Create(refreshToken.UserID)
	if err != nil {
		return "", &users.Error{Op: op, Err: err}
	}

	encodedAccessToken, err = a.accessTokenStore.Encode(accessToken)
	if err != nil {
		return "", &users.Error{Op: op, Err: err}
	}

	return encodedAccessToken, nil
}

func HashPassword(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyPassword(password, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
