package handler

import (
	httplogwrap "github.com/SyaibanAhmadRamadhan/http-log-wrap"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/mini-e-commerce-microservice/user-service/internal/services/user"
	"reflect"
)

type handler struct {
	r           *chi.Mux
	validator   *validator.Validate
	userService user.Service
}

func NewHandler(r *chi.Mux, userService user.Service) {
	r.Use(httplogwrap.HttpOtel)
	v := validator.New()
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("json")
	})

	h := &handler{
		r:           r,
		validator:   v,
		userService: userService,
	}
	h.route()
}

func (h *handler) route() {
	h.r.Post("/api/v1/register", httplogwrap.TraceHttpOtel(
		h.V1RegisterPost,
		httplogwrap.WithLogRequestBody(false),
	))
}
