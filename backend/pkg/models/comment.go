package models

import "database/sql"

type Comment struct {
	Id             int    `json:"id"`
	Post_id        int    `json:"post_id"`
	OwnerId        int    `json:"owner_id"`
	Content        string `json:"content"`
	Image          string `json:"image"`
	CreatedAt      string `json:"created_at"`
}

type CommentModel struct {
	DB *sql.DB
}

type CommentFilter struct {
	Start int    `json:"start"`
	Ncomment int    `json:"n_comment"`
}
