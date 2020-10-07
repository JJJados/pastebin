// CMPT315 Macewan University
// Assignment 1: RESTful API for Text Sharing
// Author: Jayden Laturnus

package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Server struct {
	Router    *mux.Router
	SubRouter *mux.Router
	DB        *sqlx.DB
}

func (s *Server) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (s *Server) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (s *Server) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (s *Server) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (s *Server) ReportPostHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (s *Server) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (s *Server) Initialize() {
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres dbname=assign1 sslmode=disable password=password")
	if err != nil {
		log.Fatal(err)
	}

	s.Router = mux.NewRouter()
	s.SubRouter = s.Router.PathPrefix("/api/v1").Subrouter()
	s.DB = db

	s.InitializeRoutes()
}

func (s *Server) InitializeRoutes() {
	// GET
	s.SubRouter.HandleFunc("/posts", s.GetPostsHandler).Methods(http.MethodGet)
	s.SubRouter.HandleFunc("/posts/{uuid}", s.GetPostHandler).Methods(http.MethodGet)
	// POST
	s.SubRouter.HandleFunc("/posts/create", s.CreatePostHandler).Methods(http.MethodPost)
	s.SubRouter.HandleFunc("/posts/{uuid}/update", s.UpdatePostHandler).Methods(http.MethodPost)
	s.SubRouter.HandleFunc("/posts/{uuid}/report", s.ReportPostHandler).Methods(http.MethodPost)
	// DELETE
	s.SubRouter.HandleFunc("/posts/{uuid}", s.DeletePostHandler).Methods(http.MethodDelete)
}

func main() {
	s := Server{}

	s.Initialize()

	log.Fatal(http.ListenAndServe(":3333", s.Router))
}
