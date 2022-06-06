package http

import (
	"booker/models"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func initTestingServer() *server {
	s := NewServer("testing.db")
	models.CreateNewTables(s.db)
	models.FillWithSampleData(s.db)
	return s
}

func checkResponseCode(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func loginAndReturnCookies(t *testing.T, s *server, credentials string) string {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/login/", strings.NewReader(credentials))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.router.ServeHTTP(w, r)
	checkResponseCode(t, http.StatusSeeOther, w.Code)
	cookies := strings.Join(w.Result().Header["Set-Cookie"], "; ")
	return cookies
}

func loginAsAdmin(t *testing.T, s *server) string {
	return loginAndReturnCookies(t, s, "username=admin&password=admin")
}

func loginAsBob(t *testing.T, s *server) string {
	return loginAndReturnCookies(t, s, "username=bob&password=123")
}

func loginAsAndrzej(t *testing.T, s *server) string {
	return loginAndReturnCookies(t, s, "username=pracownik&password=roku")
}

func checkEmptyRequestWithCookies(
	t *testing.T,
	s *server,
	requestType string,
	url string,
	cookies string,
	expectedCode int,
) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(requestType, url, nil)
	if cookies != "" {
		r.Header.Set("Cookie", cookies)
	}
	s.router.ServeHTTP(w, r)
	checkResponseCode(t, expectedCode, w.Code)
	return w
}

func checkResponseBodySubstring(t *testing.T, pattern string, w *httptest.ResponseRecorder) {
	if !strings.Contains(w.Body.String(), pattern) {
		t.Errorf("Expected string '%s' in response body", pattern)
	}
}

func TestIndexView(t *testing.T) {
	s := initTestingServer()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	s.router.ServeHTTP(w, r)

	checkResponseCode(t, http.StatusOK, w.Code)
	checkResponseBodySubstring(t, "Available dates:", w)
	checkResponseBodySubstring(t, "Sign in", w)
}

func TestBookedView(t *testing.T) {
	s := initTestingServer()

	checkEmptyRequestWithCookies(t, s, "GET", "/booked/", "", http.StatusForbidden)

	bob := loginAsBob(t, s)
	w := checkEmptyRequestWithCookies(t, s, "GET", "/booked/", bob, http.StatusOK)
	checkResponseBodySubstring(t, "Booked dates:", w)
}

func TestLoginView(t *testing.T) {
	s := initTestingServer()

	w := checkEmptyRequestWithCookies(t, s, "GET", "/login/", "", http.StatusOK)
	checkResponseBodySubstring(t, "Username:", w)
	checkResponseBodySubstring(t, "Password:", w)

	bob := loginAsBob(t, s)
	checkEmptyRequestWithCookies(t, s, "GET", "/login/", bob, http.StatusFound)
}

func TestBookUnbook(t *testing.T) {
	s := initTestingServer()

	bob := loginAsBob(t, s)
	admin := loginAsAdmin(t, s)

	checkEmptyRequestWithCookies(t, s, "POST", "/book/1/", "", http.StatusFound)
	checkEmptyRequestWithCookies(t, s, "POST", "/book/1/", admin, http.StatusFound)
	checkEmptyRequestWithCookies(t, s, "POST", "/book/10000000000000000000000000000000000000/", admin, http.StatusBadRequest)
	checkEmptyRequestWithCookies(t, s, "POST", "/unbook/10000000000000000000000000000000000000/", admin, http.StatusBadRequest)
	checkEmptyRequestWithCookies(t, s, "POST", "/book/1/", admin, http.StatusBadRequest)
	checkEmptyRequestWithCookies(t, s, "POST", "/unbook/1/", admin, http.StatusFound)
	checkEmptyRequestWithCookies(t, s, "POST", "/unbook/2/", admin, http.StatusBadRequest)
	checkEmptyRequestWithCookies(t, s, "POST", "/book/1/", admin, http.StatusFound)
	checkEmptyRequestWithCookies(t, s, "POST", "/book/1/", bob, http.StatusBadRequest)
	checkEmptyRequestWithCookies(t, s, "POST", "/unbook/1/", bob, http.StatusBadRequest)
	checkEmptyRequestWithCookies(t, s, "POST", "/unbook/1/", admin, http.StatusFound)
	checkEmptyRequestWithCookies(t, s, "POST", "/book/1/", bob, http.StatusFound)
}

func TestAddDateView(t *testing.T) {
	s := initTestingServer()

	bob := loginAsBob(t, s)
	checkEmptyRequestWithCookies(t, s, "GET", "/add-date/", "", http.StatusForbidden)
	checkEmptyRequestWithCookies(t, s, "GET", "/add-date/", bob, http.StatusForbidden)

	andrzej := loginAsAndrzej(t, s)
	w := checkEmptyRequestWithCookies(t, s, "GET", "/add-date/", andrzej, http.StatusOK)
	checkResponseBodySubstring(t, "start time:", w)
	checkResponseBodySubstring(t, "end time:", w)

	admin := loginAsAdmin(t, s)
	w = checkEmptyRequestWithCookies(t, s, "GET", "/add-date/", admin, http.StatusOK)
	checkResponseBodySubstring(t, "start time:", w)
	checkResponseBodySubstring(t, "end time:", w)
	checkResponseBodySubstring(t, "assign to:", w)
}

func postAddDate(
	t *testing.T,
	s *server,
	emp int,
	startTime string,
	endTime string,
	cookies string,
	expectedCode int,
) {
	var formData string
	if emp == 0 {
		formData = fmt.Sprintf("start-time=%s&end-time=%s", startTime, endTime)
	} else {
		formData = fmt.Sprintf("employee=%d&start-time=%s&end-time=%s", emp, startTime, endTime)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/add-date/", strings.NewReader(formData))
	if cookies != "" {
		r.Header.Set("Cookie", cookies)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.router.ServeHTTP(w, r)
	checkResponseCode(t, expectedCode, w.Code)
}

func TestAddingDates(t *testing.T) {
	s := initTestingServer()

	bob := loginAsBob(t, s)
	postAddDate(t, s, 0, "2022-07-13T11:30", "2022-07-13T12:30", bob, http.StatusForbidden)
	postAddDate(t, s, 1, "2022-07-13T11:30", "2022-07-13T12:30", bob, http.StatusForbidden)

	andrzej := loginAsAndrzej(t, s)
	postAddDate(t, s, 0, "2022-07-13T11:30", "2022-07-13T12:30", andrzej, http.StatusFound)
	postAddDate(t, s, 0, "2022-07-13T12:30", "2022-07-13T11:30", andrzej, http.StatusBadRequest)
	postAddDate(t, s, 1, "2022-07-13T11:30", "2022-07-13T12:30", andrzej, http.StatusForbidden)

	admin := loginAsAdmin(t, s)
	postAddDate(t, s, 0, "2022-07-13T11:30", "2022-07-13T12:30", admin, http.StatusBadRequest)
	postAddDate(t, s, 1, "2022-07-13T11:30", "2022-07-13T12:30", admin, http.StatusFound)
	postAddDate(t, s, 2, "2022-07-13T11:30", "2022-07-13T12:30", admin, http.StatusFound)
}

func TestInvalidLogin(t *testing.T) {
	s := initTestingServer()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/login/", strings.NewReader("username=admin&password=123"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.router.ServeHTTP(w, r)
	checkResponseCode(t, http.StatusBadRequest, w.Code)

	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/login/", strings.NewReader("username=admin"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.router.ServeHTTP(w, r)
	checkResponseCode(t, http.StatusBadRequest, w.Code)
}

func TestLogout(t *testing.T) {
	s := initTestingServer()

	checkEmptyRequestWithCookies(t, s, "POST", "/logout/", "", http.StatusUnauthorized)

	bob := loginAsBob(t, s)
	checkEmptyRequestWithCookies(t, s, "POST", "/logout/", bob, http.StatusFound)
}

func TestAssignedView(t *testing.T) {
	s := initTestingServer()

	bob := loginAsBob(t, s)
	checkEmptyRequestWithCookies(t, s, "GET", "/assigned/", bob, http.StatusForbidden)

	andrzej := loginAsAndrzej(t, s)
	w := checkEmptyRequestWithCookies(t, s, "GET", "/assigned/", andrzej, http.StatusOK)
	checkResponseBodySubstring(t, "Andrzej:", w)

	admin := loginAsAdmin(t, s)
	w = checkEmptyRequestWithCookies(t, s, "GET", "/assigned/", admin, http.StatusOK)
	checkResponseBodySubstring(t, "Andrzej:", w)
	checkResponseBodySubstring(t, "Fabian:", w)
}

func postAddUser(
	t *testing.T,
	s *server,
	name string,
	username string,
	password string,
	userType int,
	cookies string,
	expectedCode int,
) {
	formData := fmt.Sprintf("name=%s&username=%s&password=%s&type=%d",
		name, username, password, userType)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/add-user/", strings.NewReader(formData))
	if cookies != "" {
		r.Header.Set("Cookie", cookies)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.router.ServeHTTP(w, r)
	checkResponseCode(t, expectedCode, w.Code)
}

func TestAddUser(t *testing.T) {
	s := initTestingServer()

	admin := loginAsAdmin(t, s)
	andrzej := loginAsAndrzej(t, s)

	checkEmptyRequestWithCookies(t, s, "GET", "/add-user/", admin, http.StatusOK)
	checkEmptyRequestWithCookies(t, s, "GET", "/add-user/", andrzej, http.StatusForbidden)

	postAddUser(t, s, "Tomasz", "totomek", "kemotot", models.UserTypeEmployee, andrzej, http.StatusForbidden)
	postAddUser(t, s, "Tomasz", "totomek", "kemotot", models.UserTypeEmployee, admin, http.StatusOK)
	postAddUser(t, s, "Tomasz", "totomek", "kemotot", models.UserTypeEmployee, admin, http.StatusBadRequest)
}

func TestDisconnectedDatabase(t *testing.T) {
	s := initTestingServer()

	admin := loginAsAdmin(t, s)

	if err := s.db.Close(); err != nil {
		t.Errorf("Closing the database failed")
	}

	checkEmptyRequestWithCookies(t, s, "GET", "/", admin, http.StatusInternalServerError)
	checkEmptyRequestWithCookies(t, s, "GET", "/booked/", admin, http.StatusInternalServerError)
	checkEmptyRequestWithCookies(t, s, "POST", "/book/1/", admin, http.StatusInternalServerError)
	checkEmptyRequestWithCookies(t, s, "POST", "/unbook/1/", admin, http.StatusInternalServerError)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/login/", strings.NewReader("username=admin&password=admin"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.router.ServeHTTP(w, r)
	checkResponseCode(t, http.StatusInternalServerError, w.Code)
}
