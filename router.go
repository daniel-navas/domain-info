package main

import (
	"encoding/json"
	"net/http"

	"github.com/dfnavas/domain-info/controllers"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func createRouter(ctrl *controllers.DomainInfoCtrl) *chi.Mux {
	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	r.Get("/{url}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		url := chi.URLParam(r, "url")
		info, err := ctrl.Get(url)
		if err != nil {
			w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			json.NewEncoder(w).Encode(info)
			w.WriteHeader(http.StatusOK)
		}
	})

	return r
}
