package ports

import (
	"fmt"
	"github.com/hardiksachan/kanban_board/backend/internal/users"
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userStore    UserStore
	sessionStore SessionStore
}

func NewAuthService(userStore UserStore, sessionStore SessionStore) *AuthService {
	return &AuthService{userStore, sessionStore}
}

// SignUp returns user ID after adding it to store
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

// LogIn creates and returns a session if credentials are valid
// Otherwise, returns ENOTFOUND if email is incorrect
// returns ECONFLICT if passwords do not match
func (a *AuthService) LogIn(email, password string) (*domain.Session, error) {
	op := "ports.AuthService.Login"
	msg := fmt.Sprintf("email(%s) or password incorrect", email)

	storedUser, err := a.userStore.FindByEmail(email)
	if err != nil {
		return nil, &users.Error{Op: op, Message: msg, Err: err}
	}

	_, err = VerifyPassword(password, storedUser.Password)
	if err != nil {
		return nil, &users.Error{Op: op, Code: users.ECONFLICT, Message: msg, Err: err}
	}

	session, err := a.sessionStore.Create(storedUser.ID)
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	return session, nil
}

// LogOut deletes a user session
// Otherwise, returns error
func (a *AuthService) LogOut(sessionId string) error {
	op := "ports.AuthService.LogOut"

	err := a.sessionStore.Delete(sessionId)
	if err != nil {
		return &users.Error{Op: op, Err: err}
	}
	return nil
}

// GetSession returns a user session
// Otherwise, returns ENOTFOUND if the session Id is invalid
// returns error if fails
func (a *AuthService) GetSession(sessionId string) (*domain.Session, error) {
	op := "ports.AuthService.LogOut"

	session, err := a.sessionStore.Get(sessionId)
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}
	return session, nil
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
