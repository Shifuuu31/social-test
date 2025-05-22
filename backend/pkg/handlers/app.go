package handlers

import "social-network/pkg/models"

type SocialApp struct {
	Posts *models.PostModel
	Comments *models.CommentModel
}
