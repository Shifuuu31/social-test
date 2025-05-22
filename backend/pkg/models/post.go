package models

import (
	"database/sql"
)

type Post struct {
	Id        int    `json:"id"`
	OwnerId   int    `json:"owner_id"`
	GroupId   int    `json:"group_id"`
	Content   string `json:"content"`
	Image     string `json:"image"`
	Privacy   string `json:"Privacy"`
	CreatedAt string `json:"created_at"` // [public', 'almost_private', 'private']
}

type PostModel struct {
	DB *sql.DB
}

type PostFilter struct {
	Id    int    `json:"id"`
	Type  string `json:"type"`
	Start int    `json:"start"`
	NPost int    `json:"n_post"`
}

func (pm *PostModel) GetPosts(filter *PostFilter) (posts []Post, err error) {
	var query string
	var rows *sql.Rows

	switch filter.Type {
	case "group":
		query = `
			SELECT id, user_id, group_id, content, image, privacy, created_at
			FROM posts
			WHERE group_id = ? AND privacy = '' AND id > ?
			ORDER BY id ASC
			LIMIT ?`

		rows, err = pm.DB.Query(query, filter.Id, filter.Start, filter.NPost)

	case "privacy":
		query = `
			SELECT *
			FROM posts
			LEFT JOIN follows 
				ON follows.followee_id = posts.user_id AND follows.follower_id = ? AND follows.status = 'accepted'
			LEFT JOIN post_privacy
				ON post_privacy.post_id = posts.id AND post_privacy.user_id = ?
			WHERE NOT (posts.group_id IS NOT NULL AND posts.privacy = '')
			  AND (
				posts.privacy IN ('public')
				OR (posts.privacy = 'almost_private' AND f.follower_id IS NOT NULL)
				OR (posts.privacy = 'private' AND post_privacy.user_id IS NOT NULL)
			  )
			  AND posts.id > ?
			ORDER BY posts.id ASC
			LIMIT ?` // TODO needs testing

		rows, err = pm.DB.Query(query, filter.Id, filter.Id, filter.Start, filter.NPost)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.Id,
			&post.OwnerId,
			&post.GroupId,
			&post.Content,
			&post.Image,
			&post.Privacy,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
