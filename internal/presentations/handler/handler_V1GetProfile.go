package handler

import (
	"errors"
	"github.com/mini-e-commerce-microservice/user-service/generated/api"
	"github.com/mini-e-commerce-microservice/user-service/internal/services/user"
	"net/http"
)

func (h *handler) V1GetProfile(w http.ResponseWriter, r *http.Request) {
	authData, ok := h.getUserFromBearerAuth(w, r)
	if !ok {
		return
	}

	userData, err := h.service.UserService.GetUser(r.Context(), user.GetUserInput{
		UserID: authData.UserId,
	})
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			h.httpOtel.Err(w, r, http.StatusNotFound, err, user.ErrUserNotFound.Error())
		default:
			h.httpOtel.Err(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	resp := api.V1GetProfile200Response{
		BackgroundImage: userData.BackgroundImage.Ptr(),
		Email:           userData.Email,
		FullName:        userData.FullName,
		Id:              userData.ID,
		Image:           userData.ImageProfile.Ptr(),
		IsEmailVerified: userData.IsEmailVerified,
	}

	h.httpOtel.WriteJson(w, r, http.StatusOK, resp)
}
