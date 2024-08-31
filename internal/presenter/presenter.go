package presenter

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/mini-e-commerce-microservice/user-service/internal/presenter/handler"
	"github.com/mini-e-commerce-microservice/user-service/internal/services"
	"net/http"
	"time"
)

type Presenter struct {
	Dependency *services.Dependency
	Port       int
}

func New(presenter *Presenter) *http.Server {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3002"},
		AllowedHeaders:   []string{"Origin", "Test", "Content-Type", "Accept", "X-Request-Id", "Authorization"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	}))

	handler.NewHandler(r, presenter.Dependency.UserService)

	s := &http.Server{
		Addr:              fmt.Sprintf(":%d", presenter.Port),
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	return s
}
