package routes

import (
	"database/sql"
	"net/http"
	"tech-challenge/internal/handlers"
	"tech-challenge/internal/services"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r chi.Router, db *sql.DB) {
	c := new(handlers.CourseHandler)
	c.CourseService = services.NewCourseService(db)
	p := new(handlers.PersonHandler)
	p.PersonService = services.NewPersonService(db)
	r.Route("/api", func(r chi.Router) {
		r.Route("/course", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) { c.GetAllCourses(w, r) })
			r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) { c.GetCourse(w, r) })
			r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) { c.UpdateCourse(w, r) })
			r.Post("/", func(w http.ResponseWriter, r *http.Request) { c.CreateCourse(w, r) })
			r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) { c.DeleteCourse(w, r) })
		})
		r.Route("/person", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) { p.GetAllPeople(w, r) })
			r.Get("/{name}", func(w http.ResponseWriter, r *http.Request) { p.GetPerson(w, r) })
			r.Put("/{name}", func(w http.ResponseWriter, r *http.Request) { p.UpdatePerson(w, r) })
			r.Post("/", func(w http.ResponseWriter, r *http.Request) { p.CreatePerson(w, r) })
			r.Delete("/{name}", func(w http.ResponseWriter, r *http.Request) { p.DeletePerson(w, r) })
		})
	})
}
