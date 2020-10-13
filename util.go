// CMPT315 Macewan University
// Assignment 1: RESTful API for Text Sharing
// Author: Jayden Laturnus

package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetAccessID(r *http.Request) (uuid.UUID, error) {
	// Get AccessID from uri
	vars := mux.Vars(r)
	accessID, err := uuid.Parse(vars["uuid"])
	if err != nil {
		return accessID, err
	}
	return accessID, nil
}
