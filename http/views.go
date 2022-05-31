package http

import (
	"booker/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func (s *server) indexView(w http.ResponseWriter, r *http.Request) {
	dates, err := models.GetDatesWithNamesNotBooked(s.db)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError)
		log.Print(err)
		return
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
		log.Print(err)
		return
	}

	renderTemplate(w, r, "booked.html", dates)
}

func (s *server) loginView(w http.ResponseWriter, r *http.Request) {
	if getUser(r) == nil {
		renderTemplate(w, r, "login.html", nil)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (s *server) bookHandler(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if user == nil {
		http.Redirect(w, r, "/login/", http.StatusFound)
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
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *server) unbookHandler(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
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
	if err != nil || (date.BookedBy != user.Id && !user.IsAdmin()) || date.BookedBy == -1 {
		renderError(w, r, http.StatusBadRequest)
		return
	}

	err = models.SetDateBookedBy(s.db, dateId, -1)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/booked/", http.StatusFound)
}

func (s *server) addDateView(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if user == nil || !user.IsEmployee() {
		renderError(w, r, http.StatusForbidden)
	} else if user.IsAdmin() {
		emps, err := models.GetUsersByType(s.db, models.UserTypeEmployee)
		if err != nil {
			renderError(w, r, http.StatusInternalServerError)
			log.Println(err)
		} else {
			renderTemplate(w, r, "add_date.html", map[string]interface{}{
				"emps":   emps,
				"userId": user.Id,
			})
		}
	} else {
		renderTemplate(w, r, "add_date.html", nil)
	}
}

func (s *server) addDateHandler(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if user == nil || !user.IsEmployee() {
		renderError(w, r, http.StatusForbidden)
		return
	}

	if !verifyForm(r, "start-time", "end-time") {
		renderError(w, r, http.StatusBadRequest)
		return
	}

	empId := user.Id
	if user.IsAdmin() {
		if !r.Form.Has("employee") {
			renderError(w, r, http.StatusBadRequest)
			return
		}

		var err error
		empId, err = strconv.Atoi(r.Form.Get("employee"))
		if err != nil {
			renderError(w, r, http.StatusBadRequest)
			return
		}
	} else if r.Form.Has("employee") {
		renderError(w, r, http.StatusForbidden)
		return
	}

	const layout = "2006-01-02T15:04"
	startTime, err := time.Parse(layout, r.Form.Get("start-time"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse(layout, r.Form.Get("end-time"))
	if err != nil || startTime.After(endTime) {
		renderError(w, r, http.StatusBadRequest)
		return
	}

	// TODO: check if user[empID].isEmployee()

	err = models.CreateDate(s.db, startTime, endTime, empId)
	if err != nil {
		log.Println(err)
		renderError(w, r, http.StatusInternalServerError)
		return
	}

	s.addDateView(w, r)
}

func (s *server) addUserView(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if user == nil || !user.IsAdmin() {
		renderError(w, r, http.StatusForbidden)
	} else {
		renderTemplate(w, r, "add_user.html", nil)
	}
}

func (s *server) addUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if user == nil || !user.IsAdmin() {
		renderError(w, r, http.StatusForbidden)
		return
	}

	r.ParseForm()

	var err error
	var userType int
	userType, err = strconv.Atoi(r.Form.Get("type"));
	if err != nil {
		renderError(w, r, http.StatusBadRequest)
		return
	}

	err = models.CreateUser(s.db, r.Form.Get("name"), r.Form.Get("username"), r.Form.Get("password"), userType)
	if err != nil {
		log.Println(err);
		renderError(w, r, http.StatusBadRequest)
		return
	}

	log.Println("User added");

	s.addUserView(w, r)
}

func (s *server) assignedView(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if user == nil || !user.IsEmployee() {
		renderError(w, r, http.StatusForbidden)
		return
	}

	var dates []*models.DateWithNames
	var err error
	if user.IsAdmin() {
		dates, err = models.GetDatesWithNamesAll(s.db)
		if err != nil {
			log.Println(err)
			renderError(w, r, http.StatusInternalServerError)
			return
		}
	} else {
		dates, err = models.GetDatesWithNamesAssignedTo(s.db, user.Id)
		if err != nil {
			log.Println(err)
			renderError(w, r, http.StatusInternalServerError)
			return
		}
	}

	renderTemplate(w, r, "assigned.html", dates)
}
