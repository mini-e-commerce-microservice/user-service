package handler

import (
	whttp "github.com/SyaibanAhmadRamadhan/http-wrapper"
	"github.com/go-chi/chi/v5"
	"github.com/mini-e-commerce-microservice/user-service/internal/services"
)

type handler struct {
	r        *chi.Mux
	service  *services.Dependency
	httpOtel *whttp.Opentelemetry
}

func NewHandler(r *chi.Mux, service *services.Dependency) {

	h := &handler{
		r:       r,
		service: service,
		httpOtel: whttp.NewOtel(
			whttp.WithRecoverMode(true),
			whttp.WithPropagator(),
			whttp.WithValidator(nil, nil),
		),
	}
	h.route()
}

func (h *handler) route() {
	h.r.Post("/api/v1/register", h.httpOtel.Trace(
		h.V1RegisterPost,
		whttp.WithLogRequestBody(false),
	))

	h.r.Post("/api/v1/otp", h.httpOtel.Trace(
		h.V1SendOtpPost,
	))

	h.r.Put("/api/v1/otp", h.httpOtel.Trace(
		h.V1VerifyOtpPut,
	))

	h.r.Put("/api/v1/verify-email-user", h.httpOtel.Trace(
		h.V1VerifyEmailUser,
	))
}
