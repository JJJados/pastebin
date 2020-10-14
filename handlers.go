// CMPT315 Macewan University
// Assignment 1: RESTful API for Text Sharing
// Author: Jayden Laturnus

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

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
		publicP.Created = p.Created
		publicP.Updated = p.Updated

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

	posts := Posts{}

	query := `SELECT posts.title, post_references.read_access_uuid, posts.created,
					posts.updated FROM post_references, posts 
				WHERE post_references.post_uuid = posts.post_uuid 
					AND public_access = true AND reported = false
				ORDER BY created DESC
				LIMIT $1 OFFSET $2;`

	err := s.DB.Select(&posts, query, limit[0], offset[0])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
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
					AND post_references.read_access_uuid = :read_access_uuid
				RETURNING posts.post_uuid
			)
			INSERT INTO reported_posts (reported_uuid, reported_reason, post_uuid)
			SELECT :reported_uuid, :reported_reason, (SELECT post_uuid FROM new_report)
			WHERE EXISTS(SELECT post_uuid FROM new_report);`

	//VALUES (:reported_uuid, :reported_reason, (SELECT post_uuid FROM new_report));
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
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
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
