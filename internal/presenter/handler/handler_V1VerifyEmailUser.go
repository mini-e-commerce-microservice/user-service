package handler

import (
	"errors"
	"github.com/mini-e-commerce-microservice/user-service/generated/api"
	"github.com/mini-e-commerce-microservice/user-service/internal/services/user"
	"net/http"
)

func (h *handler) V1VerifyEmailUser(w http.ResponseWriter, r *http.Request) {
	req := api.V1VerifyEmailUserRequestBody{}

	if !h.bodyRequestBindToStruct(w, r, &req) {
		return
	}

	err := h.service.UserService.ActivationEmailUser(r.Context(), user.ActivationEmailUserInput{
		Token: req.Token,
	})
	if err != nil {
		if errors.Is(err, user.ErrTokenIsExpired) || errors.Is(err, user.ErrInvalidToken) {
			Error(w, r, http.StatusBadRequest, err, "invalid otp token")
		} else {
			Error(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
