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

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (s *Server) Initalize() {
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
	s.SubRouter.HandleFunc("/posts", GetPostsHandler).Methods(http.MethodGet)
}

func main() {
	s := Server{}

	s.Initalize()

	log.Fatal(http.ListenAndServe(":3333", s.Router))
}
