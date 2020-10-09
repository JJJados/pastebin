// CMPT315 Macewan University
// Assignment 1: RESTful API for Text Sharing
// Author: Jayden Laturnus

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Server struct {
	Router    *mux.Router
	SubRouter *mux.Router
	DB        *sqlx.DB
}

type Post struct {
	PostID        uuid.UUID `json:"postId" db:"post_uuid"`
	ReadAccessID  uuid.UUID `json:"readAccessId" db:"read_access_uuid"`
	AdminAccessID uuid.UUID `json:"adminAccessId" db:"admin_access_uuid"`
	PostTitle     string    `json:"postTitle,omitempty" db:"title"`
	PostContent   string    `json:"postContent,omitempty" db:"content"`
	PublicAccess  bool      `json:"publicAccess,omitempty" db:"public_access"`
	Reported      bool      `json:"reported,omitempty" db:"reported"`
	Created       time.Time `json:"created,omitempty" db:"created"`
	Updated       time.Time `json:"updated,omitempty" db:"updated"`
}

func (s *Server) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	// Get AccessID from uri
	vars := mux.Vars(r)
	accessID := uuid.MustParse(vars["uuid"])
	// Create a new Post struct
	p := Post{}

	query := `	SELECT * 
				FROM post_references
				NATURAL JOIN posts
				WHERE read_access_uuid = $1 OR admin_access_uuid = $1;
			`

	err := s.DB.Get(&p, query, accessID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	switch accessID {
	case p.ReadAccessID:
		w.Header().Set("Content-Type", "application/json")
		// Send Post information back to the client
		json.NewEncoder(w).Encode(p)
		fmt.Println("matches read only")

	case p.AdminAccessID:
		w.Header().Set("Content-Type", "application/json")
		// Send Post information back to the client
		json.NewEncoder(w).Encode(p)
		fmt.Println("matches admin only")

	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func (s *Server) GetPostsHandler(w http.ResponseWriter, r *http.Request) {

	posts := []Post{}

	query := `	SELECT posts.title, post_references.read_access_uuid, posts.created,
					posts.updated
				FROM post_references, posts 
				WHERE post_references.post_uuid = posts.post_uuid 
					and public_access = true AND reported != true;
			`

	err := s.DB.Select(&posts, query)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
	//http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (s *Server) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new Post struct
	p := Post{}

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	// Generate UUIDs for the postID and uri's
	p.PostID = uuid.New()
	p.ReadAccessID = uuid.New()
	p.AdminAccessID = uuid.New()
	// Assign default Reported value
	p.Reported = false
	// Assign Post creation and update time
	p.Created = time.Now()
	p.Updated = time.Now()

	query := `	WITH new_post as (
					INSERT INTO post_references (
						post_uuid, read_access_uuid, admin_access_uuid, 
						public_access
					)
						VALUES (
							:post_uuid, :read_access_uuid, :admin_access_uuid, 
							:public_access
						)
					RETURNING post_uuid
				)
				INSERT INTO posts (title, content, post_uuid)
					VALUES (:title, :content, (select post_uuid from new_post)
				);
			`

	result, err := s.DB.NamedExec(query, p)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	// Check the number of affected rows.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	fmt.Printf("%d record(s) created.\n", rowsAffected)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// Send Post information back to the client
	json.NewEncoder(w).Encode(p)
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
