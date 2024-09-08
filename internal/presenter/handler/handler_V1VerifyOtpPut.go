package handler

import (
	"errors"
	"github.com/mini-e-commerce-microservice/user-service/generated/api"
	"github.com/mini-e-commerce-microservice/user-service/internal/services/otp"
	"github.com/mini-e-commerce-microservice/user-service/internal/util/primitive"
	"net/http"
)

func (h *handler) V1VerifyOtpPut(w http.ResponseWriter, r *http.Request) {
	req := api.V1VerifyOtpPutRequestBody{}

	if !h.bodyRequestBindToStruct(w, r, &req) {
		return
	}

	verifyOutput, err := h.service.OtpService.VerifyOtp(r.Context(), otp.VerifyOtpInput{
		Usecase:            primitive.OtpUseCase(req.Usecase),
		Type:               primitive.OtpType(req.Type),
		Code:               req.Code,
		DestinationAddress: req.DestinationAddress,
	})
	if err != nil {
		if errors.Is(err, otp.ErrOtpExpired) {
			Error(w, r, http.StatusBadRequest, err, otp.ErrOtpExpired.Error())
		} else if errors.Is(err, otp.ErrOtpCounterExceeded) {
			Error(w, r, http.StatusBadRequest, err, otp.ErrOtpCounterExceeded.Error())
		} else if errors.Is(err, otp.ErrCodeOtpInvalid) || errors.Is(err, otp.ErrOtpNotFound) || errors.Is(err, otp.ErrDestinationAddressNotFound) {
			Error(w, r, http.StatusBadRequest, err, otp.ErrCodeOtpInvalid.Error())
		} else {
			Error(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	resp := api.V1VerifyOtpPutResponse{
		Token: verifyOutput.Token,
	}

	h.writeJson(w, r, http.StatusOK, resp)
}
