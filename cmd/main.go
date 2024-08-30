package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"go-notes/api"
	"net/http"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(api.AuthMiddleware)

	r.Post("/login", api.Login)
	r.Post("/notes", api.AddNote)
	r.Get("/notes", api.GetNotes)

	logrus.Info("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
