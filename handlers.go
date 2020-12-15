// CMPT315 - Pastebin
// Macewan University
// Jayden Laturnus

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./dist/html/index.html")
}

func (s *Server) PostsPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./dist/html/pastes.html")
}

func (s *Server) PostPageHandler(w http.ResponseWriter, r *http.Request) {
	// Get AccessID from uri
	accessID, err := GetAccessID(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.ParseFiles("./dist/html/paste.html"))

	p, err := s.GetPost(accessID)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	tmpl.Execute(w, p)
}

func (s *Server) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	// Get AccessID from uri
	accessID, err := GetAccessID(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	p, err := s.GetPost(accessID)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// Send Post information back to the client
	json.NewEncoder(w).Encode(p)
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

	// Checks to ensure limit is int
	_, err := strconv.Atoi(limit[0])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Get the offset param
	offset, exists := q["offset"]
	if !exists {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Checks to ensure offset is int
	_, err = strconv.Atoi(offset[0])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	posts, err := s.GetPosts(limit, offset)

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
	postID, err := NewUUID()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	readAccessID, err := NewUUID()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	adminAccessID, err := NewUUID()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

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

	result, err := s.CreatePost(p)

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

	result, err := s.UpdatePost(p)
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

	reportedID, err := NewUUID()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
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

	result, err := s.ReportPost(p)
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
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	result, err := s.DeletePost(accessID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
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
