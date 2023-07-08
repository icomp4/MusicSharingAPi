package controllers

import (
	"log"
	"musicSharingAPp/db"
	"musicSharingAPp/models"
	"strconv"
	"strings"
)

type LoginStruct struct {
	Username string
	Password string
}

func IsFollowing(userID, followID uint) bool {
	var user models.User
	db.DB.Preload("Following").First(&user, userID)

	for _, followingUser := range user.Following {
		if followingUser.ID == followID {
			return true
		}
	}

	return false
}

func SignUp(user *models.User) error {
	user.Username = strings.ToLower(user.Username)
	err := db.DB.Create(&user).Error
	if err != nil {
		return err
	}
	log.Println("Successfully created user: ", user.Username)
	return nil
}
func Login(login LoginStruct) (string, error) {
	var user models.User
	err := db.DB.Where("Username = ?", login.Username).First(&user).Error
	if err != nil {
		return "", err
	}
	if login.Password != user.Password {
		return "password incorrect", nil
	}
	return strconv.FormatUint(uint64(user.ID), 10), nil

}
func DeleteAcc(id string) error {
	err := db.DB.Delete(&models.User{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
func GetUserInfo(id string) (*models.User, error) {
	var user models.User
	result := db.DB.Preload("FavGenres").Preload("Following").Preload("Followers").Preload("Blocked").Preload("Playlists").Preload("Posts").First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
func GetAllUsersInfo() (*[]models.User, error) {
	var user []models.User
	result := db.DB.Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func FollowAccount(currentID string, FollowID string) string {
	var currentUser models.User
	var userToFollow models.User
	err := db.DB.First(&currentUser, currentID).Error
	err2 := db.DB.First(&userToFollow, FollowID).Error
	if err != nil || err2 != nil {
		return "Users not found"
	}
	following := IsFollowing(currentUser.ID, userToFollow.ID)
	if following {
		return "Already following this user"
	}
	tx := db.DB.Begin()
	if err := tx.Model(&currentUser).Association("Following").Append(&userToFollow); err != nil {
		tx.Rollback()
		return "Failed to follow the user"
	}
	if err := tx.Model(&userToFollow).Association("Followers").Append(&currentUser); err != nil {
		tx.Rollback()
		return "Failed to follow the user"
	}
	if err := tx.Commit().Error; err != nil {
		return "Failed to unfollow the user"
	}
	currentUser.FollowingCount++
	userToFollow.FollowerCount++
	db.DB.Save(&currentUser)
	db.DB.Save(&userToFollow)
	return "User successfully followed!"
}
func UnfollowAccount(currentID string, UnfollowID string) string {
	var currentUser models.User
	var userToUnfollow models.User
	err := db.DB.First(&currentUser, currentID).Error
	err2 := db.DB.First(&userToUnfollow, UnfollowID).Error
	if err != nil || err2 != nil {
		return "Users not found"
	}
	following := IsFollowing(currentUser.ID, userToUnfollow.ID)
	if !following {
		return "Not following this user"
	}
	tx := db.DB.Begin()
	if err := tx.Model(&currentUser).Association("Following").Delete(&userToUnfollow); err != nil {
		tx.Rollback()
		return "Failed to unfollow the user"
	}
	if err := tx.Model(&userToUnfollow).Association("Followers").Delete(&currentUser); err != nil {
		tx.Rollback()
		return "Failed to unfollow the user"
	}
	if err := tx.Commit().Error; err != nil {
		return "Failed to unfollow the user"
	}
	currentUser.FollowingCount--
	userToUnfollow.FollowerCount--
	db.DB.Save(&currentUser)
	db.DB.Save(&userToUnfollow)
	return "User successfully unfollowed"
}