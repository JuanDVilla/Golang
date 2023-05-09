package main

import (
	"net/http"

	// "github.com/go-chi/chi"
	// "github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger) // middleware para registrar las solicitudes HTTP
	r.Use(middleware.Recoverer) // middleware para recuperarse de los errores de la aplicación

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("¡Bienvenido a mi API!"))
	})

	r.Get("/users/{userID}", func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		w.Write([]byte("Has solicitado el usuario: " + userID))
	})

	http.ListenAndServe(":3000", r)
}
