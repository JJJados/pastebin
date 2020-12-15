// CMPT315 - Pastebin
// Macewan University
// Jayden Laturnus

package main

import "database/sql"

func (s *Server) GetPost(accessID string) (Post, error) {
	// Create a new Post struct
	p := Post{}

	query := `SELECT * FROM post_references NATURAL JOIN posts
				WHERE read_access_uuid = $1 OR admin_access_uuid = $1;`

	err := s.DB.Get(&p, query, accessID)
	if err != nil {
		return p, err
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

		return publicP, nil

	case *p.AdminAccessID:
		return p, nil

	default:
		return p, err
	}
}

func (s *Server) GetPosts(limit []string, offset []string) (Posts, error) {
	posts := Posts{}

	query := `SELECT posts.title, post_references.read_access_uuid, posts.created,
					posts.updated FROM post_references, posts 
				WHERE post_references.post_uuid = posts.post_uuid 
					AND public_access = true AND reported = false
				ORDER BY created DESC
				LIMIT $1 OFFSET $2;`

	err := s.DB.Select(&posts, query, limit[0], offset[0])
	if err != nil {
		return posts, err
	}

	return posts, nil
}

func (s *Server) CreatePost(p Post) (sql.Result, error) {
	query := `WITH new_post as (
		INSERT INTO post_references (
			post_uuid, read_access_uuid, admin_access_uuid)
			VALUES (:post_uuid, :read_access_uuid, :admin_access_uuid)
		RETURNING post_uuid
	)
	INSERT INTO posts (title, content, public_access, post_uuid)
		VALUES (:title, :content, :public_access, (select post_uuid from new_post));`

	return s.DB.NamedExec(query, p)
}

func (s *Server) UpdatePost(p Post) (sql.Result, error) {
	query := `UPDATE posts SET title = :title, content = :content,
					public_access = :public_access, updated = :updated
				FROM post_references
				WHERE posts.post_uuid = post_references.post_uuid AND
				  	admin_access_uuid = :admin_access_uuid;`

	return s.DB.NamedExec(query, p)
}

func (s *Server) ReportPost(p ReportedPost) (sql.Result, error) {
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

	return s.DB.NamedExec(query, p)
}

func (s *Server) DeletePost(accessID string) (sql.Result, error) {
	query := `DELETE FROM post_references WHERE admin_access_uuid = $1;`

	return s.DB.Exec(query, accessID)
}
