package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nytlabs/st-core/core"
)

type ConnectionNode struct {
	Id    int `json:"id"`
	Route int `json:"route"`
}

type ConnectionLedger struct {
	Source ConnectionNode `json:"source"`
	Target ConnectionNode `json:"target"`
	Id     int            `json:"id"`
}

func (s *Server) ListConnections() []ConnectionLedger {
	var connections []ConnectionLedger
	for _, c := range s.connections {
		connections = append(connections, *c)
	}
	return connections
}

func (s *Server) ConnectionIndex(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()

	c := s.ListConnections()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(c); err != nil {
		panic(err)
	}
}

// CreateConnectionHandler responds to a POST request to instantiate a new connection
func (s *Server) ConnectionCreate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	var newConn ConnectionLedger
	json.Unmarshal(body, &newConn)

	s.Lock()
	defer s.Unlock()

	source, ok := s.blocks[newConn.Source.Id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"source block does not exist"})
		return
	}

	target := s.blocks[newConn.Target.Id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"target block does not exist"})
		return
	}

	sourceRoute := core.RouteID(newConn.Source.Route)
	targetRoute, err := target.Block.GetRoute(core.RouteID(newConn.Target.Route))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	err = source.Block.Connect(sourceRoute, targetRoute.C)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	newConn.Id = s.GetNextID()
	s.connections[newConn.Id] = &newConn

	out, _ := json.Marshal(newConn)
	// broadcast state update
	s.broadcast <- out

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) ConnectionDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ids, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"no ID supplied"})
		return
	}

	id, err := strconv.Atoi(ids)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	s.Lock()
	defer s.Unlock()

	c, ok := s.connections[id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"connection not found"})
		return
	}

	source, ok := s.blocks[c.Source.Id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"source block does not exist"})
		return
	}

	target, ok := s.blocks[c.Target.Id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{"target block does not exist"})
		return
	}

	route, err := target.Block.GetRoute(core.RouteID(c.Target.Route))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	err = source.Block.Disconnect(core.RouteID(c.Source.Route), route.C)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
