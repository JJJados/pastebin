// CMPT315 - Assignment 2
// Macewan University
// Jayden Laturnus

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

// This was taken from the example in the gorilla mux docs https://github.com/gorilla/mux
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI, "Method: ", r.Method, "Body: ", r.Body)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func (s *Server) Initialize() {
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres dbname=assign1 sslmode=disable password=password")
	if err != nil {
		log.Fatal(err)
	}

	s.Router = mux.NewRouter()
	s.SubRouter = s.Router.PathPrefix("/api/v1").Subrouter()
	s.DB = db

	s.Router.Use(LoggingMiddleware)

	s.InitializeRoutes()
}

func (s *Server) InitializeRoutes() {
	// GET
	s.SubRouter.HandleFunc("/posts", s.GetPostsHandler).Methods(http.MethodGet)
	s.SubRouter.HandleFunc("/posts/{uuid}", s.GetPostHandler).Methods(http.MethodGet)
	// POST
	s.SubRouter.HandleFunc("/posts", s.CreatePostHandler).Methods(http.MethodPost)
	s.SubRouter.HandleFunc("/posts/{uuid}", s.UpdatePostHandler).Methods(http.MethodPut)
	s.SubRouter.HandleFunc("/posts/{uuid}/reports", s.ReportPostHandler).Methods(http.MethodPost)
	// DELETE
	s.SubRouter.HandleFunc("/posts/{uuid}", s.DeletePostHandler).Methods(http.MethodDelete)

	// Serve webpage files
	s.Router.HandleFunc("/pastes/{uuid}", s.PostPageHandler).Methods(http.MethodGet)
	s.Router.HandleFunc("/pastes", s.PostsPageHandler).Methods(http.MethodGet)
	s.Router.HandleFunc("/", s.IndexPageHandler).Methods(http.MethodGet)
	s.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("dist")))

}

func main() {
	s := Server{}

	s.Initialize()

	log.Fatal(http.ListenAndServe(":3333", s.Router))
}
