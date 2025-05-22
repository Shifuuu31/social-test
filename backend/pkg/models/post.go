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
	Privacy   string `json:"Privacy"` // [public', 'almost_private', 'private']
	CreatedAt string `json:"created_at"`
	ChosenUsersIds []int `json:"chosen_user"`
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

func (pm *PostModel) GetPosts(filter *PostFilter) error {
	var query string
	switch filter.Type {
	case "group":
		query = ""
	case "privacy":

		query = ""

	}
	_ = query

	// pm.db.exec
	return nil
}
