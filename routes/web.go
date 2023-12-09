package routes

import (
	"github.com/go-chi/chi"
	"net/http"
)

func WebRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("web"))
	})

	return r
}
