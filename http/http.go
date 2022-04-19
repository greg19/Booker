package http

import (
	"booker/models"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type server struct {
	router *chi.Mux
	db     *sql.DB
}

func NewServer() *server {
	s := server{
		router: chi.NewRouter(),
		db:     models.ConnectToDatabase(),
	}
	s.registerHandlers()
	return &s
}

func (s *server) Run(addr string) {
	log.Println("Starting server on " + addr)
	log.Fatal(http.ListenAndServe(addr, s.router))
}

func (s *server) registerHandlers() {
	r := s.router

	r.Use(s.readUser)
	r.Use(middleware.Logger)

	r.Get("/", s.indexView)
	r.Get("/booked/", s.bookedView)
	r.Get("/login/", s.loginView)
	r.Get("/logout/", s.logoutHandler)

	r.Post("/login/", s.loginHandler)
	r.Post("/book/{dateId:[0-9]+}/", s.bookHandler)
	r.Post("/unbook/{dateId:[0-9]+}/", s.unbookHandler)

	fs := http.FileServer(http.Dir("web/static/"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))
}
