package presenter

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strconv"
	"time"
	"user-service/internal/util"
)

func New() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Post("/", WithOtel(largeResponseHandler))

	err := http.ListenAndServe(":3333", r)
	util.Panic(err)
}

// Profile is a sample profile structure
type Profile struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

// LargeResponse represents a larger JSON response structure
type LargeResponse struct {
	Profiles []Profile `json:"profiles"`
	Meta     MetaData  `json:"meta"`
}

// MetaData represents additional metadata
type MetaData struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

// Example handler that sends a large JSON response
func largeResponseHandler(w http.ResponseWriter, r *http.Request) {
	// Creating a large list of profiles
	profiles := make([]Profile, 0, 1000)
	for i := 1; i <= 1000; i++ {
		profiles = append(profiles, Profile{
			ID:      i,
			Name:    "Name " + strconv.Itoa(i),
			Address: "Address " + strconv.Itoa(i),
		})
	}

	// Creating metadata
	meta := MetaData{
		Page:       1,
		PageSize:   1000,
		TotalCount: 1000,
		TotalPages: 1,
	}

	// Creating large response
	largeResponse := LargeResponse{
		Profiles: profiles,
		Meta:     meta,
	}

	// Marshalling large response to JSON
	b, err := json.Marshal(largeResponse)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Writing response
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
