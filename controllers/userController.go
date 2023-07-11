package controllers

import (
	"fmt"
	"log"
	"musicSharingAPp/db"
	"musicSharingAPp/models"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
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
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil{
		return err
	}
	usernameTaken := db.DB.Where("Username = ?", user.Username).First(&models.User{}).Error
	if usernameTaken != nil{
		return usernameTaken
	}
	user.Password = string(password)
	fmt.Println(password)
	err2 := db.DB.Create(&user).Error
	if err2 != nil {
		return err2
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

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
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
func UpdatePassword(currentID string, newPassword string) string{
	var currentUser models.User
	err := db.DB.First(&currentUser,currentID).Error
	if err != nil{
		return "Could not find account"
	}
	password, err := bcrypt.GenerateFromPassword([]byte(currentUser.Password), bcrypt.DefaultCost)
	if err != nil{
		log.Fatal(err)
	}
	currentUser.Password = string(password)
	db.DB.Save(&currentUser)
	return "Successsfully updated password"
}
func UpdatePFP(currentID string, newImgURL string) string{
	var currentUser models.User
	err := db.DB.First(&currentUser,currentID).Error
	if err != nil{
		return "Could not find account"
	}
	currentUser.PfpURL = newImgURL
	db.DB.Save(&currentUser)
	return "Successsfully updated pfp"
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
