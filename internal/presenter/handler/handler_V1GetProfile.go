package handler

import (
	"net/http"
	"user-service/generated/api"
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
