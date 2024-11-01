package handler

import (
	whttp "github.com/SyaibanAhmadRamadhan/http-wrapper"
	"github.com/go-chi/chi/v5"
	"github.com/mini-e-commerce-microservice/user-service/generated/proto/secret_proto"
	"github.com/mini-e-commerce-microservice/user-service/internal/services"
)

type handler struct {
	r                  *chi.Mux
	service            *services.Dependency
	httpOtel           *whttp.Opentelemetry
	jwtAccessTokenConf *secret_proto.JwtAccessToken
}

func NewHandler(r *chi.Mux, service *services.Dependency, jwtAccessTokenConf *secret_proto.JwtAccessToken) {

	h := &handler{
		r:                  r,
		service:            service,
		jwtAccessTokenConf: jwtAccessTokenConf,
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

	h.r.Get("/api/v1/profile", h.httpOtel.Trace(
		h.V1GetProfile,
	))
}
