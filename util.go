// CMPT315 - Assignment 2
// Macewan University
// Jayden Laturnus

package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type UUID [16]byte

// UUID generation was taken from https://play.golang.org/p/4FkNSiUDMg
func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func GetAccessID(r *http.Request) (string, error) {
	// Get AccessID from uri
	vars := mux.Vars(r)
	accessID := vars["uuid"]
	if len(accessID) == 0 {
		return accessID, errors.New("Cannot find ID")
	}
	return accessID, nil
}

/*
Checks if publicAccess is private for templating
This may seem counter intutive but for some reason within templating
it would not compare true/false correctly without the help of this
function.
*/
func (p Post) IsPrivate(publicAccess bool) bool {
	if publicAccess {
		return true
	}
	return false
}
