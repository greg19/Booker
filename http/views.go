package http

import (
	"booker/models"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const templatesDir = "web/templates/"

type templateContext struct {
	User *models.User
	Data interface{}
}

func renderTemplate(w http.ResponseWriter, r *http.Request, filename string, data interface{}) {
	t, err := template.ParseFiles(templatesDir+"layout.html", templatesDir+filename)
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, templateContext{
		User: r.Context().Value("user").(*models.User),
		Data: data,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func renderError(w http.ResponseWriter, r *http.Request, statusCode int) {
	w.WriteHeader(statusCode)
	msg := fmt.Sprintf("%d: %s", statusCode, http.StatusText(statusCode))
	renderTemplate(w, r, "error.html", msg)
}

func (s *server) indexView(w http.ResponseWriter, r *http.Request) {
	dates, err := models.GetDatesNotBooked(s.db)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError)
		log.Fatal(err)
	}
	renderTemplate(w, r, "index.html", dates)
}

func (s *server) bookedView(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	if user == nil {
		renderError(w, r, http.StatusForbidden)
		return
	}

	dates, err := models.GetDatesBookedBy(s.db, user.Id)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError)
		log.Fatal(err)
	}
	renderTemplate(w, r, "booked.html", dates)
}

func (s *server) loginView(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "login.html", nil)
}

func (s *server) bookHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	if user == nil {
		renderError(w, r, http.StatusForbidden)
		return
	}

	dateId, err := strconv.Atoi(chi.URLParam(r, "dateId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest)
		return
	}

	date, err := models.GetDateById(s.db, dateId)
	if err != nil || date.BookedBy != -1 {
		renderError(w, r, http.StatusBadRequest)
		return
	}

	// TODO: remove race condition
	err = models.SetDateBookedBy(s.db, dateId, user.Id)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *server) unbookHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	if user == nil {
		renderError(w, r, http.StatusForbidden)
		return
	}

	dateId, err := strconv.Atoi(chi.URLParam(r, "dateId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest)
		return
	}

	date, err := models.GetDateById(s.db, dateId)
	if err != nil || date.BookedBy != user.Id {
		renderError(w, r, http.StatusBadRequest)
		return
	}

	// TODO: remove race condition
	err = models.SetDateBookedBy(s.db, dateId, -1)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/booked/", http.StatusFound)
}
