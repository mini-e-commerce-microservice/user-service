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

	if !h.httpOtel.BindBodyRequest(w, r, &req) {
		return
	}

	err := h.service.OtpService.SendOtp(r.Context(), otp.SendOtpInput{
		Usecase:            primitive.OtpUseCase(req.Usecase),
		Type:               primitive.OtpType(req.Type),
		DestinationAddress: req.DestinationAddress,
	})
	if err != nil {
		if errors.Is(err, otp.ErrDestinationAddressNotFound) {
			h.httpOtel.Err(w, r, http.StatusBadRequest, err, otp.ErrDestinationAddressNotFound.Error())
		} else if errors.Is(err, otp.ErrEmailUserIsVerified) {
			h.httpOtel.Err(w, r, http.StatusBadRequest, err, otp.ErrEmailUserIsVerified.Error())
		} else {
			h.httpOtel.Err(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
