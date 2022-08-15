package ports

import (
	"fmt"
	"github.com/hardiksachan/kanban_board/backend/internal/users"
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	credentialStore CredentialStore
	refreshStore    RefreshStore
	accessProvider  AccessProvider
}

func NewAuthService(credentialStore CredentialStore, accessProvider AccessProvider, refreshStore RefreshStore) *AuthService {
	return &AuthService{credentialStore, refreshStore, accessProvider}
}

func (a *AuthService) SignUp(credential *domain.Credential) (*domain.Credential, error) {
	op := "ports.AuthService.SignUp"

	count, err := a.credentialStore.CountByEmail(credential.Email)
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}
	if count != 0 {
		return nil, &users.Error{Op: op, Code: users.ECONFLICT, Message: fmt.Sprintf("credential with email (%s) exists", credential.Email)}
	}

	hashedPassword, err := HashPassword(credential.Password)
	if err != nil {
		return nil, &users.Error{Op: op, Code: users.EINTERNAL}
	}
	credential.Password = hashedPassword

	storedCred, err := a.credentialStore.Insert(credential)
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	return storedCred, nil
}

func (a *AuthService) LogIn(email, password string) (*domain.Credential, *domain.RefreshToken, *domain.AccessToken, error) {
	op := "ports.AuthService.Login"
	msg := fmt.Sprintf("email(%s) or password incorrect", email)

	credential, err := a.credentialStore.FindByEmail(email)
	if err != nil {
		return nil, nil, nil, &users.Error{Op: op, Message: msg, Err: err}
	}

	passwordValid, err := VerifyPassword(password, credential.Password)
	if err != nil || !passwordValid {
		return nil, nil, nil, &users.Error{Op: op, Code: users.ECONFLICT, Message: msg, Err: err}
	}

	refreshToken, err := a.refreshStore.Create(credential)
	if err != nil {
		return nil, nil, nil, &users.Error{Op: op, Err: err}
	}

	accessToken, err := a.accessProvider.Create(&domain.AccessClaims{UserID: credential.UserID})
	if err != nil {
		return nil, nil, nil, &users.Error{Op: op, Err: err}
	}

	return credential, refreshToken, accessToken, nil
}

func (a *AuthService) LogOut(refreshToken *domain.RefreshToken) error {
	op := "ports.AuthService.LogOut"

	err := a.refreshStore.Revoke(refreshToken)
	if err != nil {
		return &users.Error{Op: op, Err: err}
	}
	return nil
}

func (a *AuthService) DecodeAccessToken(accessToken *domain.AccessToken) (*domain.Credential, error) {
	op := "ports.AuthService.DecodeAccessToken"

	accessClaims, err := a.accessProvider.Verify(accessToken)
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}
	if accessClaims == nil {
		return nil, &users.Error{Op: op, Message: "invalid access token", Code: users.EINVALID}
	}

	credential, err := a.credentialStore.FindById(accessClaims.UserID)
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	return credential, nil
}

func (a *AuthService) RegenerateAccessToken(refreshToken *domain.RefreshToken) (encodedAccessToken *domain.AccessToken, err error) {
	op := "ports.AuthService.RegenerateAccessToken"

	credential, err := a.refreshStore.Verify(refreshToken)
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}
	if credential == nil {
		return nil, &users.Error{Op: op, Code: users.EINVALID, Err: err}
	}

	accessToken, err := a.accessProvider.Create(&domain.AccessClaims{UserID: credential.UserID})
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	return accessToken, nil
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
