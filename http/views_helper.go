package http

import (
	"booker/models"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"runtime"
)

func getTemplatesDir() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	basepath = path.Join(basepath, "../web/templates/")
	return basepath
}

var templatesDir = getTemplatesDir()

type templateContext struct {
	User  *models.User
	Data  interface{}
	Error interface{}
}

func readErrorMessage(r *http.Request) interface{} {
	if r.Context().Value("error") == nil {
		return nil
	} else {
		return r.Context().Value("error").(string)
	}
}

func getUser(r *http.Request) *models.User {
	val := r.Context().Value("user")
	if val != nil {
		return val.(*models.User)
	} else {
		return nil
	}
}

func renderTemplate(w http.ResponseWriter, r *http.Request, filename string, data interface{}) {
	t, err := template.ParseFiles(
		path.Join(templatesDir, "layout.html"),
		path.Join(templatesDir, filename),
	)
	if err != nil {
		log.Println(err)
		return
	}

	err = t.Execute(w, templateContext{
		User:  getUser(r),
		Error: readErrorMessage(r),
		Data:  data,
	})
	if err != nil {
		log.Println(err)
	}
}

func addError(w http.ResponseWriter, r *http.Request, statusCode int, err string) {
	w.WriteHeader(statusCode)
	ctx := context.WithValue(r.Context(), "error", err)
	*r = *r.WithContext(ctx)
}

func errorMessageFromStatus(statusCode int) string {
	return fmt.Sprintf("%d: %s", statusCode, http.StatusText(statusCode))
}

func addStatusCodeError(w http.ResponseWriter, r *http.Request, statusCode int) {
	addError(w, r, statusCode, errorMessageFromStatus(statusCode))
}

func renderError(w http.ResponseWriter, r *http.Request, statusCode int) {
	addStatusCodeError(w, r, statusCode)
	renderTemplate(w, r, "error.html", nil)
}

func verifyForm(r *http.Request, fields ...string) bool {
	err := r.ParseForm()
	if err != nil {
		return false
	}

	for _, field := range fields {
		if !r.Form.Has(field) {
			return false
		}
	}

	return true
}
