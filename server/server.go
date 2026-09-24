package server

import (
	"log"
	"net"
	"net/http"

	"github.com/sarthak-shubham/cherry/store"
)

type Server struct {
	store *store.Store
}

func NewServer(store *store.Store) *Server {
	return &Server{
		store: store,
	}
}

func (s *Server) Start(address string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/kv/", s.handleKV)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("Cherry server is listening on %s", address)

	server := &http.Server{
		Handler: mux,
	}

	return server.Serve(listener)
}
