package http

import (
	"booker/models"
	"database/sql"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type httpServer struct {
	templates map[string]*template.Template
	mux       *http.ServeMux
	db        *sql.DB
}

func NewServer() *httpServer {
	s := httpServer{
		templates: make(map[string]*template.Template),
		mux:       http.NewServeMux(),
		db:        nil,
	}
	s.loadTemplates()
	s.registerHandlers()
	s.connectDatabase()
	return &s
}

func (s *httpServer) Run(addr string) {
	log.Println("Starting server on " + addr)
	log.Fatal(http.ListenAndServe(addr, s.mux))
}

func (s *httpServer) loadTemplates() {
	// Loads all .html files in templatesDir as templates
	const templatesDir = "web/templates/"

	files, err := ioutil.ReadDir(templatesDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".html") && filename != "layout.html" {
			s.templates[filename] = template.Must(
				template.ParseFiles(templatesDir+"layout.html", templatesDir+filename))
		}
	}
}

func (s *httpServer) renderTemplate(w http.ResponseWriter, filename string, data interface{}) {
	t, ok := s.templates[filename]
	if !ok {
		log.Fatal("template '" + filename + "' not found")
	}
	t.Execute(w, data)
}

func (s *httpServer) connectDatabase() {
	var err error
	s.db, err = sql.Open("sqlite3", "booker.db")
	if err != nil {
		log.Fatal(err)
	}

	models.InitTables(s.db)
}
