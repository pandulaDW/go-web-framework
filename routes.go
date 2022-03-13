package celeritas

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (c *Celeritas) routes() http.Handler {
	mux := chi.NewRouter()

	// middlewares
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)

	if c.Debug {
		mux.Use(middleware.Logger)
	}

	mux.Use(middleware.Recoverer)

	// routes
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "Welcome to Celeritas")
	})

	return mux
}
