package controllers

import (
	"musicSharingAPp/db"
	"musicSharingAPp/models"
	"strconv"
)

func CreatePost(id string, post *models.Post) string {
	var currentUser models.User
	err := db.DB.First(&currentUser, id).Error
	if err != nil {
		return "User not found"
	}
	newPost := models.Post{
		Title:   post.Title,
		Likes:   0,
		Content: post.Content,
		UserID:  currentUser.ID,
		Author:  currentUser.Username,
	}
	err2 := db.DB.Create(&newPost).Error
	if err2 != nil {
		return "Failed to create post"
	}
	currentUser.Posts = append(currentUser.Posts, newPost)
	currentUser.PostCount++
	err = db.DB.Save(&currentUser).Error
	if err != nil {
		return "Failed to update user's posts"
	}
	return "success"
}
func GetPostsByUserID(id string) ([]models.Post, error) {
	var user models.User
	if err := db.DB.Preload("Posts").First(&user, id).Error; err != nil {
		return nil, err
	}
	return user.Posts, nil
}
func DeletePost(userid string, postid string) string {
	var currentUser models.User
	err := db.DB.First(&currentUser, userid).Error
	if err != nil {
		return "User not found"
	}
	var post models.Post
	err = db.DB.First(&post, postid).Error
	if strconv.FormatUint(uint64(post.UserID), 10) != userid {
		return "Incorrect user"
	}
	if err != nil {
		return "Post not found"
	}
	err = db.DB.Delete(&post).Error
	if err != nil {
		return "Failed to delete post"
	}

	currentUser.PostCount--
	err = db.DB.Save(&currentUser).Error
	if err != nil {
		return "Failed to update user's posts"
	}

	return "success"
}

func LikePost(userID, postID string) string {
	var currentUser models.User
	err := db.DB.Preload("LikedPosts").First(&currentUser, userID).Error
	if err != nil {
		return "User not found"
	}
	var post models.Post
	err = db.DB.First(&post, postID).Error
	if err != nil {
		return "Post not found"
	}
	for _, likedPost := range currentUser.LikedPosts {
		if likedPost.ID == post.ID {
			return "Post already liked"
		}
	}
	tx := db.DB.Begin()
	currentUser.LikedPosts = append(currentUser.LikedPosts, post)
	err = tx.Save(&currentUser).Error
	if err != nil {
		tx.Rollback()
		return "Failed to like the post"
	}
	post.Likes++
	err = tx.Save(&post).Error
	if err != nil {
		tx.Rollback()
		return "Failed to like the post"
	}
	tx.Commit()
	return "Success"
}
func UnlikePost(userID, postID string) string {
	var currentUser models.User
	var postsToUnlike models.Post
	err := db.DB.First(&currentUser, userID).Error
	err2 := db.DB.First(&postsToUnlike, postID).Error
	if err != nil || err2 != nil {
		return "User and/or post not found"
	}
	tx := db.DB.Begin()
	if err := tx.Model(&currentUser).Association("LikedPosts").Delete(&postsToUnlike); err != nil {
		tx.Rollback()
		return "Failed to unlike the post"
	}
	if err := tx.Commit().Error; err != nil {
		return "Failed to unlike the post"
	}
	postsToUnlike.Likes--
	db.DB.Save(&currentUser)
	db.DB.Save(&postsToUnlike)
	return "User successfully unfollowed"
}

