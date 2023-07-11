package controllers

import (
	"musicSharingAPp/db"
	"musicSharingAPp/models"
)

func GetFeed(id string) []models.Post {
	var user models.User
	var posts []models.Post

	err := db.DB.Find(&user, id).Error
	if err != nil {
		return []models.Post{}
	}

	tx := db.DB.Begin()
	if err := tx.Model(&user).Association("Following").Find(&user.Following); err != nil {
		tx.Rollback()
		return []models.Post{}
	}

	for _, followedUser := range user.Following {
		err := tx.Model(&followedUser).Association("Posts").Find(&followedUser.Posts)
		if err != nil {
			tx.Rollback()
			return []models.Post{}
		}
		posts = append(posts, followedUser.Posts...)
	}

	tx.Commit()
	return posts
}