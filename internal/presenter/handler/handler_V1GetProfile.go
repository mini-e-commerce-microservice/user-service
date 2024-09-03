package handler

import (
	"net/http"
)

func (h *handler) V1GetProfile(w http.ResponseWriter, r *http.Request) {
	// TODO: need implement
	h.writeJson(w, r, http.StatusNotImplemented, map[string]string{
		"message": "need impl",
	})
	return
}
