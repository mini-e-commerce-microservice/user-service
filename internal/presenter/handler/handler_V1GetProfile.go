package handler

import (
	"github.com/mini-e-commerce-microservice/user-service/generated/api"
	"net/http"
)

func (h *handler) V1GetProfile(w http.ResponseWriter, r *http.Request) {

	resp := api.V1GetProfile200Response{
		BackgroundImage: "",
		Email:           "",
		FullName:        "",
		Id:              0,
		Image:           "",
	}

	h.writeJson(w, r, http.StatusOK, resp)
}
