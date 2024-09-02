package handler

import (
	"errors"
	"github.com/mini-e-commerce-microservice/user-service/generated/api"
	"github.com/mini-e-commerce-microservice/user-service/internal/services/otp"
	"github.com/mini-e-commerce-microservice/user-service/internal/util/primitive"
	"net/http"
)

func (h *handler) V1SendOtpPost(w http.ResponseWriter, r *http.Request) {
	req := api.V1SendOtpPostRequestBody{}

	if !h.bodyRequestBindToStruct(w, r, &req) {
		return
	}

	userID, ok := h.getUserID(w, r)
	if !ok {
		return
	}

	err := h.service.OtpService.SendOtp(r.Context(), otp.SendOtpInput{
		Usecase:            primitive.OtpUseCase(req.Usecase),
		UserID:             userID,
		Type:               primitive.OtpType(req.Type),
		DestinationAddress: req.DestinationAddress,
	})
	if err != nil {
		if errors.Is(err, otp.ErrDestinationAddressNotFound) {
			Error(w, r, http.StatusBadRequest, err, otp.ErrDestinationAddressNotFound.Error())
		} else {
			Error(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
