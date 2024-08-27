package handler

import (
	"net/http"
	"user-service/generated/api"
)

func (h *handler) V1RegisterPost(w http.ResponseWriter, r *http.Request) {
	req := api.V1RegisterPostRequestBody{}

	if !h.bodyRequestBindToStruct(w, r, &req) {
		return
	}

	resp := api.V1RegisterPost200Response{
		Id:                    0,
		UploadBackgroundImage: nil,
		UploadProfileImage:    nil,
	}

	h.writeJson(w, r, http.StatusOK, resp)
}
