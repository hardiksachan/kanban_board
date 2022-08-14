package user

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/hardiksachan/kanban_board/backend/internal/users"
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/ports"
	"github.com/hardiksachan/kanban_board/backend/internal/users/handlers/auth"
	jsonHelper "github.com/hardiksachan/kanban_board/backend/shared/json"
	"github.com/hardiksachan/kanban_board/backend/shared/logging"
	"net/http"
)

type Handler struct {
	user     *ports.UserService
	log      logging.Logger
	validate *validator.Validate
}

func NewUserHandler(user *ports.UserService, log logging.Logger, validator *validator.Validate) *Handler {
	return &Handler{user, log, validator}
}

func (h *Handler) Update(rw http.ResponseWriter, r *http.Request) {
	loggedInUserID := r.Context().Value(&auth.UserIDKey{}).(string)
	rUserID := mux.Vars(r)["user_id"]

	if loggedInUserID != rUserID {
		h.log.Debug(fmt.Sprintf("user %s cannot update user %s", loggedInUserID, rUserID))

		http.Error(rw, "cannot update a different", http.StatusBadRequest)
		return
	}

	rm, err := jsonHelper.Parse[UpdateRequest](r.Body)
	if err != nil {
		h.log.Debug(fmt.Sprintf("unable to parse request body. err: %s", err.Error()))

		http.Error(rw, "unable to parse request body", http.StatusBadRequest)
		return
	}

	// validate and sanitize input
	validationErr := h.validate.Struct(rm)
	if validationErr != nil {
		h.log.Debug(fmt.Sprintf("invalid request Body. err: %s", validationErr.Error()))

		http.Error(rw, fmt.Sprintf("invalid request. %s", validationErr.Error()), http.StatusBadRequest)
		return
	}

	err = h.user.Update(&domain.UserMetadata{
		UserId:      rUserID,
		DisplayName: rm.DisplayName,
		ImageURL:    rm.ProfileURL,
	})
	if err != nil {
		h.log.Debug(fmt.Sprintf("unable to update metadata. err: %s", err.Error()))

		http.Error(rw, users.ErrorMessage(err), http.StatusInternalServerError)
		return
	}

	h.log.Debug(fmt.Sprintf("metadata update succesfull. userId: %s", rUserID))
	rw.WriteHeader(http.StatusCreated)
}

func (h *Handler) Get(rw http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["user_id"]

	user, err := h.user.Find(userId)
	if err != nil {
		switch users.ErrorCode(err) {
		case users.ENOTFOUND:
			h.log.Debug(fmt.Sprintf("invalid user. err: %s", err.Error()))

			http.Error(rw, users.ErrorMessage(err), http.StatusNotFound)
			return
		}
	}

	h.log.Debug(fmt.Sprintf("user fetch succesfull. userId: %s", userId))
	json.NewEncoder(rw).Encode(user)
}
