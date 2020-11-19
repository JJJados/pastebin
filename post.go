// CMPT315 - Assignment 2
// Macewan University
// Jayden Laturnus

package main

import (
	"time"
)

type Post struct {
	PostID        *string   `json:"postId,omitempty" db:"post_uuid"`
	ReadAccessID  *string   `json:"readAccessId,omitempty" db:"read_access_uuid"`
	AdminAccessID *string   `json:"adminAccessId,omitempty" db:"admin_access_uuid"`
	PostTitle     string    `json:"postTitle,omitempty" db:"title"`
	PostContent   string    `json:"postContent,omitempty" db:"content"`
	PublicAccess  *bool     `json:"publicAccess,omitempty" db:"public_access"`
	Reported      bool      `json:"reported,omitempty" db:"reported"`
	Created       time.Time `json:"created,omitempty" db:"created"`
	Updated       time.Time `json:"updated,omitempty" db:"updated"`
}

type Posts []Post

type ReportedPost struct {
	ReadAccessID   *string `json:"readAccessId" db:"read_access_uuid"`
	ReportedID     *string `json:"reportedID" db:"reported_uuid"`
	Reported       bool    `json:"reported" db:"reported"`
	ReportedReason string  `json:"reportedReason" db:"reported_reason"`
}
