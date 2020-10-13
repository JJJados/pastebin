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
	PostID        *uuid.UUID `json:"postId,omitempty" db:"post_uuid"`
	ReadAccessID  *uuid.UUID `json:"readAccessId,omitempty" db:"read_access_uuid"`
	AdminAccessID *uuid.UUID `json:"adminAccessId,omitempty" db:"admin_access_uuid"`
	PostTitle     string     `json:"postTitle,omitempty" db:"title"`
	PostContent   string     `json:"postContent,omitempty" db:"content"`
	PublicAccess  *bool      `json:"publicAccess,omitempty" db:"public_access"`
	Reported      *bool      `json:"reported,omitempty" db:"reported"`
	Created       time.Time  `json:"created,omitempty" db:"created"`
	Updated       time.Time  `json:"updated,omitempty" db:"updated"`
}

type ReportedPost struct {
	ReadAccessID   *uuid.UUID `json:"readAccessId" db:"read_access_uuid"`
	ReportedID     *uuid.UUID `json:"reportedID" db:"reported_uuid"`
	Reported       bool       `json:"reported" db:"reported"`
	ReportedReason string     `json:"reportedReason" db:"reported_reason"`
}

func GetAccessID(r *http.Request) (uuid.UUID, error) {
	// Get AccessID from uri
	vars := mux.Vars(r)
	accessID, err := uuid.Parse(vars["uuid"])
	if err != nil {
		return accessID, err
	}
	return accessID, nil
}

func (s *Server) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	// Get AccessID from uri
	accessID, err := GetAccessID(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Create a new Post struct
	p := Post{}

	query := `SELECT * FROM post_references NATURAL JOIN posts
				WHERE read_access_uuid = $1 OR admin_access_uuid = $1;`

	err = s.DB.Get(&p, query, accessID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	switch accessID {
	case *p.ReadAccessID:
		// Create a new post
		publicP := Post{}
		// Assign only public facing info to new post
		publicP.ReadAccessID = p.ReadAccessID
		publicP.PostTitle = p.PostTitle
		publicP.PostContent = p.PostContent
		publicP.PublicAccess = p.PublicAccess
		publicP.Reported = p.Reported

		w.Header().Set("Content-Type", "application/json")
		// Send Post information back to the client
		json.NewEncoder(w).Encode(publicP)

	case *p.AdminAccessID:
		w.Header().Set("Content-Type", "application/json")
		// Send Post information back to the client
		json.NewEncoder(w).Encode(p)

	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func (s *Server) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	// Get query params
	q := r.URL.Query()
	// Get the limit param
	limit, exists := q["limit"]
	if !exists {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// Get the offset param
	offset, exists := q["offset"]
	if !exists {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	posts := []Post{}

	query := `SELECT posts.title, post_references.read_access_uuid, posts.created,
					posts.updated FROM post_references, posts 
				WHERE post_references.post_uuid = posts.post_uuid 
					AND public_access = true AND reported = false
				ORDER BY created DESC
				LIMIT $1 OFFSET $2;`

	err := s.DB.Select(&posts, query, limit[0], offset[0])
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (s *Server) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new Post struct
	p := Post{}

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// Check if there is a post title and content
	if len(p.PostTitle) == 0 || len(p.PostContent) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Generate UUIDs for the postID and uri's
	postID := uuid.New()
	readAccessID := uuid.New()
	adminAccessID := uuid.New()
	// Assign values to pointers
	p.PostID = &postID
	p.ReadAccessID = &readAccessID
	p.AdminAccessID = &adminAccessID
	// Assign Post creation and update time
	p.Created = time.Now()
	p.Updated = time.Now()

	// Checks if public access was set, if not defaults to true
	if p.PublicAccess == nil {
		publicAccess := true
		p.PublicAccess = &publicAccess
	}

	query := `WITH new_post as (
					INSERT INTO post_references (
						post_uuid, read_access_uuid, admin_access_uuid)
						VALUES (:post_uuid, :read_access_uuid, :admin_access_uuid)
					RETURNING post_uuid
				)
				INSERT INTO posts (title, content, public_access, post_uuid)
					VALUES (:title, :content, :public_access, (select post_uuid from new_post));`

	result, err := s.DB.NamedExec(query, p)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Check the number of affected rows.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	fmt.Printf("%d record(s) created.\n", rowsAffected)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// Send Post information back to the client
	json.NewEncoder(w).Encode(p)
}

func (s *Server) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Get AccessID from uri
	accessID, err := GetAccessID(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Create a new Post struct
	p := Post{}

	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// Assign AccessID to the AdminAccessID to check
	p.AdminAccessID = &accessID
	// Update the last time anything changed on the post
	p.Updated = time.Now()

	query := `UPDATE posts SET title = :title, content = :content, 
					public_access = :public_access, updated = :updated
				FROM post_references
				WHERE posts.post_uuid = post_references.post_uuid AND 
				  	admin_access_uuid = :admin_access_uuid;`

	result, err := s.DB.NamedExec(query, p)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Check the number of affected rows.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// If no rows were affected, post does not exist or read access link was used
	if rowsAffected == 0 {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	fmt.Printf("%d record(s) updated.\n", rowsAffected)
}

func (s *Server) ReportPostHandler(w http.ResponseWriter, r *http.Request) {
	// Get AccessID from uri
	accessID, err := GetAccessID(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Create a new Post struct
	p := ReportedPost{}

	reportedID := uuid.New()
	// Set report uuid
	p.ReportedID = &reportedID
	p.ReadAccessID = &accessID

	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Check if there is a reported reason
	if len(p.ReportedReason) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	query := `with new_report as (
				UPDATE posts SET reported = :reported
				FROM post_references 
				WHERE posts.post_uuid = post_references.post_uuid 
					AND read_access_uuid = :read_access_uuid
				RETURNING posts.post_uuid
			)
			INSERT INTO reported_posts (reported_uuid, reported_reason, post_uuid)
					VALUES (:reported_uuid, :reported_reason, (select post_uuid from new_report));`

	result, err := s.DB.NamedExec(query, p)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Check the number of affected rows.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// If no rows were affected, post does not exist
	if rowsAffected == 0 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	fmt.Printf("%d record(s) reported.\n", rowsAffected)
}

func (s *Server) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	// Get AccessID from uri
	accessID, err := GetAccessID(r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	query := `DELETE FROM post_references WHERE admin_access_uuid = $1;`

	result, err := s.DB.Exec(query, accessID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Check the number of affected rows.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// If no rows were affected, post does not exist and or read access link was used
	if rowsAffected == 0 {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	fmt.Printf("%d record(s) deleted.\n", rowsAffected)
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
