package http

import (
	"booker/models"
	"net/http"
)

func (s *httpServer) registerHandlers() {
	s.mux.HandleFunc("/", s.indexHandler)

	fs := http.FileServer(http.Dir("web/static/"))
	s.mux.Handle("/static/", http.StripPrefix("/static/", fs))
}

// example handler
func (s *httpServer) indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	counter := models.GetCounter(s.db, "licznik")
	if counter == nil {
		counter = &models.Counter{Name: "licznik", Cnt: 1}
		models.AddCounter(s.db, counter)
	}
	models.UpdateCounter(s.db, "licznik", counter.Cnt+1)
	s.renderTemplate(w, "index.html", counter.Cnt)
}
