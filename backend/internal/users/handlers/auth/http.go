package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/hardiksachan/kanban_board/backend/internal/users"
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/ports"
	jsonHelper "github.com/hardiksachan/kanban_board/backend/shared/json"
	"github.com/hardiksachan/kanban_board/backend/shared/logging"
	"net/http"
)

type CredentialKey struct {
}

type UserIDKey struct {
}

type Handler struct {
	auth      *ports.AuthService
	log       logging.Logger
	validator *validator.Validate
}

func NewAuthHandler(auth *ports.AuthService, log logging.Logger, validator *validator.Validate) *Handler {
	return &Handler{auth, log, validator}
}

func (h *Handler) SignUp(rw http.ResponseWriter, r *http.Request) {
	// Read request body as SignUpRequest
	rm, err := jsonHelper.Parse[SignUpRequest](r.Body)
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

func (h *Handler) LogIn(rw http.ResponseWriter, r *http.Request) {
	// Read request body as LogInRequest
	rm, err := jsonHelper.Parse[LogInRequest](r.Body)
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
	storedUser, refreshToken, accessToken, err := h.auth.LogIn(rm.Email, rm.Password)
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

	h.log.Debug(fmt.Sprintf("user logged in successfully. accessToken: %s, refreshToken: %s", accessToken, refreshToken))

	json.NewEncoder(rw).Encode(LogInResponse{
		AccessToken:  string(*accessToken),
		RefreshToken: string(*refreshToken),
		UserId:       storedUser.UserID,
	})
}

func (h *Handler) LogOut(rw http.ResponseWriter, r *http.Request) {
	rm, err := jsonHelper.Parse[LogOutRequest](r.Body)
	if err != nil {
		h.log.Debug(fmt.Sprintf("unable to parse request body. err: %s", err.Error()))

		http.Error(rw, "unable to parse request body", http.StatusBadRequest)
		return
	}

	refreshToken := (domain.RefreshToken)(rm.RefreshToken)

	err = h.auth.LogOut(&refreshToken)
	if err != nil {
		h.log.Warn(err.Error())
		http.Error(rw, users.ErrorMessage(err), http.StatusInternalServerError)
		return
	}

	h.log.Debug(fmt.Sprintf("user logged out successfully. refreshToken: %s", rm.RefreshToken))
	rw.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RefreshAccessToken(rw http.ResponseWriter, r *http.Request) {
	rm, err := jsonHelper.Parse[RefreshAccessTokenRequest](r.Body)
	if err != nil {
		h.log.Debug(fmt.Sprintf("unable to parse request body. err: %s", err.Error()))

		http.Error(rw, "unable to parse request body", http.StatusBadRequest)
		return
	}

	refreshToken := (domain.RefreshToken)(rm.RefreshToken)
	accessToken, err := h.auth.RegenerateAccessToken(&refreshToken)
	if err != nil {
		h.log.Warn(err.Error())
		http.Error(rw, users.ErrorMessage(err), http.StatusInternalServerError)
		return
	}

	h.log.Debug(fmt.Sprintf("access token generated successfully. accessToken: %s", accessToken))
	json.NewEncoder(rw).Encode(&RefreshAccessTokenResponse{
		AccessToken: string(*accessToken),
	})
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		accessTokenStr := r.Header.Get("Authorization")
		if accessTokenStr == "" {
			h.log.Debug(fmt.Sprintf("access token not provided"))

			http.Error(rw, "no access token", http.StatusUnauthorized)
			return
		}

		accessToken := (domain.AccessToken)(accessTokenStr)

		credential, err := h.auth.DecodeAccessToken(&accessToken)
		if err != nil {
			switch users.ErrorCode(err) {
			case users.EINVALID:
				h.log.Debug(err.Error())
				http.Error(rw, "not authenticated", http.StatusUnauthorized)
			case users.EEXPIRED:
				h.log.Debug(err.Error())
				http.Error(rw, "token expired", http.StatusUnauthorized)
			default:
				h.log.Warn(err.Error())
				http.Error(rw, users.ErrorMessage(err), http.StatusInternalServerError)
			}
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, &CredentialKey{}, credential)
		ctx = context.WithValue(ctx, &UserIDKey{}, credential.UserID)

		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
