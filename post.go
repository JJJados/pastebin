// CMPT315 Macewan University
// Assignment 1: RESTful API for Text Sharing
// Author: Jayden Laturnus

package main

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	PostID        *uuid.UUID `json:"postId,omitempty" db:"post_uuid"`
	ReadAccessID  *uuid.UUID `json:"readAccessId,omitempty" db:"read_access_uuid"`
	AdminAccessID *uuid.UUID `json:"adminAccessId,omitempty" db:"admin_access_uuid"`
	PostTitle     string     `json:"postTitle,omitempty" db:"title"`
	PostContent   string     `json:"postContent,omitempty" db:"content"`
	PublicAccess  *bool      `json:"publicAccess,omitempty" db:"public_access"`
	Reported      bool       `json:"reported,omitempty" db:"reported"`
	Created       time.Time  `json:"created,omitempty" db:"created"`
	Updated       time.Time  `json:"updated,omitempty" db:"updated"`
}

type Posts []Post

type ReportedPost struct {
	ReadAccessID   *uuid.UUID `json:"readAccessId" db:"read_access_uuid"`
	ReportedID     *uuid.UUID `json:"reportedID" db:"reported_uuid"`
	Reported       bool       `json:"reported" db:"reported"`
	ReportedReason string     `json:"reportedReason" db:"reported_reason"`
}
