package server

import (
	"encoding/json"
	"net/http"

	"github.com/nytlabs/st-core/core"
)

// naming confusion between "name" and "type" ~_~
type LibraryEntry struct {
	Type   string          `json:"type"`
	Source core.SourceType `json:"source"`
	// type if we need that later
}

func (s *Server) BlockLibraryHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()

	l := []LibraryEntry{}

	for _, v := range s.library {
		l = append(l, LibraryEntry{
			v.Name,
			v.Source,
		})
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(l); err != nil {
		panic(err)
	}
}

func (s *Server) SourceLibraryHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()

	l := []LibraryEntry{}

	for _, v := range s.sourceLibrary {
		l = append(l, LibraryEntry{
			v.Name,
			v.Type,
		})
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(l); err != nil {
		panic(err)
	}
}
