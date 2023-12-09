package routes

import (
	"github.com/go-chi/chi"
	"net/http"
)

func ApiRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("users"))
	})

	return r
}
