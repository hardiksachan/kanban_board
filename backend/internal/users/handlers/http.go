package handlers

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/hardiksachan/kanban_board/backend/internal/users"
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/ports"
	"github.com/hardiksachan/kanban_board/backend/shared/json"
	"github.com/hardiksachan/kanban_board/backend/shared/logging"
	"net/http"
	"time"
)

type SessionKey struct {
}

type UsersHandler struct {
	auth      *ports.AuthService
	log       logging.Logger
	validator *validator.Validate
}

func NewUsersHandler(auth *ports.AuthService, log logging.Logger, validator *validator.Validate) *UsersHandler {
	return &UsersHandler{auth, log, validator}
}

func (h *UsersHandler) SignUp(rw http.ResponseWriter, r *http.Request) {
	// Read request body as SignUpRequest
	rm, err := json.Parse[SignUpRequest](r.Body)
	if err != nil {
		h.log.Debug(fmt.Sprintf("unable to parse request body. err: %s", err.Error()))

		http.Error(rw, "unable to parse request body", http.StatusBadRequest)
		return
	}

	// validate and sanitize input
	validationErr := h.validator.Struct(rm)
	if validationErr != nil {
		h.log.Debug(fmt.Sprintf("invalid request Body. err: %s", validationErr.Error()))

		http.Error(rw, fmt.Sprintf("invalid request. %s", validationErr.Error()), http.StatusBadRequest)
		return
	}

	// call application layer to SignUp user
	signedUpUser, err := h.auth.SignUp(rm.toDomain())
	if err != nil {
		switch users.ErrorCode(err) {
		case users.ECONFLICT:
			h.log.Debug(err.Error())
			http.Error(rw, users.ErrorMessage(err), http.StatusBadRequest)
		default:
			h.log.Warn(err.Error())
			http.Error(rw, users.ErrorMessage(err), http.StatusInternalServerError)
		}
		return
	}

	h.log.Debug(fmt.Sprintf("user signed up successfully. user: %+v", signedUpUser))
	rw.WriteHeader(http.StatusCreated)
}

func (h *UsersHandler) LogIn(rw http.ResponseWriter, r *http.Request) {
	// Read request body as LogInRequest
	rm, err := json.Parse[LogInRequest](r.Body)
	if err != nil {
		h.log.Debug(fmt.Sprintf("unable to parse request body. err: %s", err.Error()))

		http.Error(rw, "unable to parse request body", http.StatusBadRequest)
		return
	}

	// validate and sanitize input
	validationErr := h.validator.Struct(rm)
	if validationErr != nil {
		h.log.Debug(fmt.Sprintf("invalid request Body. err: %s", validationErr.Error()))

		http.Error(rw, fmt.Sprintf("invalid request. %s", validationErr.Error()), http.StatusBadRequest)
		return
	}

	// call application layer to Log In user
	session, err := h.auth.LogIn(rm.Email, rm.Password)
	if err != nil {
		switch users.ErrorCode(err) {
		case users.ECONFLICT, users.ENOTFOUND:
			h.log.Debug(err.Error())
			http.Error(rw, users.ErrorMessage(err), http.StatusBadRequest)
		default:
			h.log.Warn(err.Error())
			http.Error(rw, users.ErrorMessage(err), http.StatusInternalServerError)
		}
		return
	}

	h.log.Debug(fmt.Sprintf("user logged in successfully. session: %+v", session))

	http.SetCookie(rw, &http.Cookie{
		Name:    "session",
		Value:   session.ID,
		Expires: session.ExpiresAt,
	})
	rw.WriteHeader(http.StatusOK)
}

func (h *UsersHandler) LogOut(rw http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(&SessionKey{}).(*domain.Session)

	err := h.auth.LogOut(session.ID)
	if err != nil {
		h.log.Warn(err.Error())
		http.Error(rw, users.ErrorMessage(err), http.StatusInternalServerError)
		return
	}

	h.log.Debug(fmt.Sprintf("user logged out successfully. session: %+v", session))
	rw.WriteHeader(http.StatusNoContent)
}

func (h *UsersHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session")
		if err != nil {
			h.log.Debug("Session Cookie not found")

			http.Error(rw, "not authenticated", http.StatusUnauthorized)
			return
		}

		session, err := h.auth.GetSession(sessionCookie.Value)
		if err != nil {
			switch users.ErrorCode(err) {
			case users.ENOTFOUND:
				h.log.Debug(err.Error())
				http.Error(rw, "not authenticated", http.StatusUnauthorized)
			default:
				h.log.Warn(err.Error())
				http.Error(rw, users.ErrorMessage(err), http.StatusInternalServerError)
			}
			return
		}

		if session.ExpiresAt.Unix() < time.Now().Unix() {
			h.log.Debug(fmt.Sprintf("Session timed out. session: %+v", session))

			http.Error(rw, "session timed out", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, &SessionKey{}, session)

		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
