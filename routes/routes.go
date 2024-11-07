package routes

import (
	"students/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/students", func (r chi.Router) {
		r.Get("/", handlers.GetAllStudents)
		r.Post("/", handlers.CreateStudent)
		r.Get("/{id}", handlers.GetStudentById)
		r.Put("/{id}", handlers.UpdateStudent)
		r.Delete("/{id}", handlers.DeleteStudent)
		r.Get("/{id}/summary", handlers.GetStudentSummary)
	})

	return r;
}