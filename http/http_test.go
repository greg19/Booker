package http

import (
	"booker/models"
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

var counter int

func checkResponseCode(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("%d Expected response code %d. Got %d\n", counter, expected, actual)
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

func checkEmptyRequestWithCookies(
	t *testing.T,
	s *server,
	requestType string,
	url string,
	cookies string,
	expectedCode int,
) {
	counter += 1
	w := httptest.NewRecorder()
	r := httptest.NewRequest(requestType, url, nil)
	if cookies != "" {
		r.Header.Set("Cookie", cookies)
	}
	s.router.ServeHTTP(w, r)
	checkResponseCode(t, expectedCode, w.Code)
}

func TestBookUnbook(t *testing.T) {
	s := initTestingServer()

	checkEmptyRequestWithCookies(t, s, "POST", "/book/1/", "", http.StatusFound)
	admin := loginAndReturnCookies(t, s, "username=admin&password=admin")
	bob := loginAndReturnCookies(t, s, "username=bob&password=123")
	checkEmptyRequestWithCookies(t, s, "POST", "/book/1/", admin, http.StatusFound)
	checkEmptyRequestWithCookies(t, s, "POST", "/book/1/", admin, http.StatusBadRequest)
	checkEmptyRequestWithCookies(t, s, "POST", "/unbook/1/", admin, http.StatusFound)
	checkEmptyRequestWithCookies(t, s, "POST", "/unbook/2/", admin, http.StatusBadRequest)
	checkEmptyRequestWithCookies(t, s, "POST", "/book/1/", admin, http.StatusFound)
	checkEmptyRequestWithCookies(t, s, "POST", "/book/1/", bob, http.StatusBadRequest)
	checkEmptyRequestWithCookies(t, s, "POST", "/unbook/1/", bob, http.StatusBadRequest)
	checkEmptyRequestWithCookies(t, s, "POST", "/unbook/1/", admin, http.StatusFound)
	checkEmptyRequestWithCookies(t, s, "POST", "/book/1/", bob, http.StatusFound)
}

func TestLogin(t *testing.T) {
	s := initTestingServer()
	_ = loginAndReturnCookies(t, s, "username=admin&password=admin")
}
