package handler

import (
	httplogwrap "github.com/SyaibanAhmadRamadhan/http-log-wrap"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/mini-e-commerce-microservice/user-service/internal/services"
	"net/http"
	"reflect"
)

type handler struct {
	r         *chi.Mux
	validator *validator.Validate
	service   *services.Dependency
}

func NewHandler(r *chi.Mux, service *services.Dependency) {
	r.Use(func(h http.Handler) http.Handler {
		return httplogwrap.HttpOtel(h,
			httplogwrap.WithExtraHeaders("X-User-Id"),
			httplogwrap.WithOutSetRequestIDHeader(),
		)
	})

	v := validator.New()
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("json")
	})

	h := &handler{
		r:         r,
		validator: v,
		service:   service,
	}
	h.route()
}

func (h *handler) route() {
	h.r.Post("/api/v1/register", httplogwrap.TraceHttpOtel(
		h.V1RegisterPost,
		httplogwrap.WithLogRequestBody(false),
	))

	h.r.Post("/api/v1/otp", httplogwrap.TraceHttpOtel(
		h.V1SendOtpPost,
	))

	h.r.Put("/api/v1/otp", httplogwrap.TraceHttpOtel(
		h.V1VerifyOtpPut,
	))
}
