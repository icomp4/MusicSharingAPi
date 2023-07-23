package router

import (
	"encoding/json"
	"errors"
	"musicSharingAPp/controllers"
	"musicSharingAPp/models"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var secretKey string
var store *sessions.CookieStore

func init() {
	err := godotenv.Load()
	if err != nil {
		return
	}
	secretKey = os.Getenv("SECRET_KEY")
	store = sessions.NewCookieStore([]byte(os.Getenv(secretKey)))
}

func HandleSignUp(w http.ResponseWriter, r *http.Request) {
    var user *models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Failed to decode request body", http.StatusBadRequest)
        return
    }

    if user.Username == "" || user.Password == "" {
        http.Error(w, "Fields must not be blank", http.StatusBadRequest)
        return
    }

	err3 := validatePassword(user.Password)
	if err3 != nil {
		http.Error(w, "Please enter a strong password", http.StatusBadRequest)
		return
	}
	
    err2 := controllers.SignUp(user)
    if err2 != nil {
        http.Error(w, "Failed to create user", http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
func UserLogin(w http.ResponseWriter, r *http.Request) {
	var login controllers.LoginStruct
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Could not get session", http.StatusBadGateway)
		return
	}
	err2 := json.NewDecoder(r.Body).Decode(&login)
	if err2 != nil {
		http.Error(w, "Could not decode response body", http.StatusBadGateway)
		return
	}
	userID, err3 := controllers.Login(login)
	if err3 != nil {
		http.Error(w, "Failed to login", http.StatusBadRequest)
		return
	}
	if userID == "password incorrect" {
		http.Error(w, "Incorrect password", http.StatusBadRequest)
		return
	}
	session.Values["userID"] = userID
	session.Values["Authorized"] = true
	session.Save(r, w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully logged in"))
}
func UserLogout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Could not get session", http.StatusBadGateway)
		return
	}
	session.Values["Authorized"] = false
	session.Save(r, w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully logged out"))
}

func DeleteAcc(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Could not get session", http.StatusBadGateway)
		return
	}
	isAuth := session.Values["Authorized"]
	if isAuth != true {
		http.Error(w, "User not authroized", http.StatusBadGateway)
		return
	}
	id := session.Values["userID"].(string)
	_, err2 := controllers.GetUserInfo(id)
	if err2 != nil {
		http.Error(w, "Could not find account", http.StatusBadRequest)
		return
	}
	controllers.DeleteAcc(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Successfully deleted account"))
}
func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Could not get session", http.StatusBadGateway)
		return
	}
	isAuth := session.Values["Authorized"]
	if isAuth != true {
		http.Error(w, "User not authroized", http.StatusBadGateway)
		return
	}
	id := session.Values["userID"].(string)
	user, err := controllers.GetUserInfo(id)
	if err != nil {
		http.Error(w, "Could not find account", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(user)
}
func GetUserInfoByID(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Could not get session", http.StatusBadGateway)
		return
	}
	isAuth := session.Values["Authorized"]
	if isAuth != true {
		http.Error(w, "User not authroized", http.StatusBadGateway)
		return
	}
	id := chi.URLParam(r, "id")
	user, err2 := controllers.GetUserInfo(id)
	if err2 != nil {
		http.Error(w, "Could not find account", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(user)
}
func GetAllUsersInfo(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Could not get session", http.StatusBadGateway)
		return
	}
	isAuth := session.Values["Authorized"]
	if isAuth != true {
		http.Error(w, "User not authroized", http.StatusBadGateway)
		return
	}
	users, err2 := controllers.GetAllUsersInfo()
	if err2 != nil {
		http.Error(w, "Could not find accounts", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(users)
}
func FollowUser(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Could not get session", http.StatusBadGateway)
		return
	}
	isAuth := session.Values["Authorized"]
	if isAuth != true {
		http.Error(w, "User not authroized", http.StatusBadGateway)
		return
	}
	id := session.Values["userID"].(string)
	userToFollowID := chi.URLParam(r, "FollowID")
	response := controllers.FollowAccount(id, userToFollowID)
	if response == "Users not found" {
		http.Error(w, "Users not found", http.StatusBadRequest)
		return
	}
	if response == "Already following" {
		http.Error(w, "User is already following that user", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("User successfully followed!"))
}
func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Could not get session", http.StatusBadGateway)
		return
	}
	isAuth := session.Values["Authorized"]
	if isAuth != true {
		http.Error(w, "User not authroized", http.StatusBadGateway)
		return
	}
	id := session.Values["userID"].(string)
	userToUnfollowID := chi.URLParam(r, "UnfollowID")
	response := controllers.UnfollowAccount(id, userToUnfollowID)
	if response == "Users not found" {
		http.Error(w, "Users not found", http.StatusBadRequest)
		return
	}
	if response == "Not following this user" {
		http.Error(w, "User is not following that user", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("User successfully unfollowed!"))
}
func CreatePost(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	var post *models.Post
	if err != nil {
		http.Error(w, "Could not get session", http.StatusBadRequest)
		return
	}
	isAuth := session.Values["Authorized"]
	if isAuth != true {
		http.Error(w, "User not authroized", http.StatusBadGateway)
		return
	}
	id := session.Values["userID"].(string)
	err2 := json.NewDecoder(r.Body).Decode(&post)
	if err2 != nil {
		http.Error(w, "Could not decode response body", http.StatusBadGateway)
		return
	}
	if IsStringEmpty(post.Title) || IsStringEmpty(post.Content) {
		http.Error(w, "Fields must not be blank", http.StatusBadGateway)
		return
	}
	resp := controllers.CreatePost(id, post)
	if resp == "User not found" {
		http.Error(w, "User not found", http.StatusBadGateway)
		return
	}
	if resp == "Failed to create post" {
		http.Error(w, "Failed to create post", http.StatusBadGateway)
		return
	}
	if resp == "Failed to update user's posts" {
		http.Error(w, "Failed to update user's posts", http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Post successfully created!"))
}
func GetCurrentUserPosts(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Could not get session", http.StatusBadRequest)
		return
	}
	isAuth := session.Values["Authorized"]
	if isAuth != true {
		http.Error(w, "User not authroized", http.StatusBadRequest)
		return
	}
	id := session.Values["userID"].(string)
	posts, err2 := controllers.GetPostsByUserID(id)
	if err2 != nil {
		http.Error(w, "Failed to get posts with specified id", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}
func GetPostsViaUserID(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Could not get session", http.StatusBadRequest)
		return
	}
	isAuth := session.Values["Authorized"]
	if isAuth != true {
		http.Error(w, "User not authroized", http.StatusBadRequest)
		return
	}
	id := chi.URLParam(r, "id")
	posts, err2 := controllers.GetPostsByUserID(id)
	if err2 != nil {
		http.Error(w, "Failed to get posts with specified id", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}
func DeletePost(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Could not get session", http.StatusBadRequest)
		return
	}
	isAuth := session.Values["Authorized"]
	if isAuth != true {
		http.Error(w, "User not authroized", http.StatusBadRequest)
		return
	}
	id := session.Values["userID"].(string)
	postid := chi.URLParam(r, "postID")
	resp := controllers.DeletePost(id, postid)
	if resp == "Post not found" {
		http.Error(w, "Post not found", http.StatusBadRequest)
		return
	}
	if resp == "User not found" {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}
	if resp == "Failed to delete post" {
		http.Error(w, "Failed to delete post", http.StatusBadRequest)
		return
	}
	if resp == "Failed to update user's posts" {
		http.Error(w, "Failed to update user's posts", http.StatusBadRequest)
		return
	}
	if resp == "Incorrect user" {
		http.Error(w, "Current user is not authorized to delete that post", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post successfully deleted"))
}
func LikePost(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Could not get session", http.StatusBadRequest)
		return
	}
	isAuth := session.Values["Authorized"]
	if isAuth != true {
		http.Error(w, "User not authroized", http.StatusBadRequest)
		return
	}
	id := session.Values["userID"].(string)
	postid := chi.URLParam(r, "postID")
	err2 := controllers.LikePost(id, postid)
	if err2 == "User not found" {
		http.Error(w, "User not authroized", http.StatusBadRequest)
		return
	}
	if err2 == "Post not found" {
		http.Error(w, "Post not found", http.StatusBadRequest)
		return
	}
	if err2 == "Post already liked" {
		http.Error(w, "Post already liked", http.StatusBadRequest)
		return
	}
	if err2 == "Failed to like the post" {
		http.Error(w, "Failed to like the post", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post successfully liked"))
}
func UnlikePost(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Could not get session", http.StatusBadRequest)
		return
	}
	isAuth := session.Values["Authorized"]
	if isAuth != true {
		http.Error(w, "User not authroized", http.StatusBadRequest)
		return
	}
	id := session.Values["userID"].(string)
	postid := chi.URLParam(r, "postID")
	err2 := controllers.UnlikePost(id, postid)
	if err2 == "User and/or post not found" {
		http.Error(w, "User and/or post not found", http.StatusBadRequest)
		return
	}
	if err2 == "Failed to unlike the post" {
		http.Error(w, "Failed to unlike the post", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post successfully unliked"))
}

func validatePassword(password string) error {
	hasUppercase := `[A-Z]`
	hasLowercase := `[a-z]`
	hasDigit := `[0-9]`
	hasSpecialChar := `[@$!%*#?&]`

	match, err := regexp.MatchString(hasUppercase, password)
	if err != nil {
		return err
	}
	if !match {
		return errors.New("password must contain at least one uppercase letter")
	}

	match, err = regexp.MatchString(hasLowercase, password)
	if err != nil {
		return err
	}
	if !match {
		return errors.New("password must contain at least one lowercase letter")
	}
	match, err = regexp.MatchString(hasDigit, password)
	if err != nil {
		return err
	}
	if !match {
		return errors.New("password must contain at least one digit")
	}
	match, err = regexp.MatchString(hasSpecialChar, password)
	if err != nil {
		return err
	}
	if !match {
		return errors.New("password must contain at least one special character (@$!%*#?&)")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Could not get session", http.StatusBadRequest)
		return
	}
	isAuth := session.Values["Authorized"]
	if isAuth != true {
		http.Error(w, "User not authorized", http.StatusBadRequest)
		return
	}
	id := session.Values["userID"].(string)
	var newPassword struct {
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&newPassword)
	if validatePassword(newPassword.Password) != nil {
		http.Error(w, "Password does not meet requirements", http.StatusBadRequest)
		return
	}
	resp := controllers.UpdatePassword(id, newPassword.Password)
	if resp == "Could not find account" {
		http.Error(w, "Could not find account", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully updated password"))
}
func UpdatePFP(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Could not get session", http.StatusBadRequest)
		return
	}
	isAuth := session.Values["Authorized"]
	if isAuth != true {
		http.Error(w, "User not authorized", http.StatusBadRequest)
		return
	}
	id := session.Values["userID"].(string)
	var newPFP struct {
		Url string `json:"url"`
	}
	json.NewDecoder(r.Body).Decode(&newPFP)
	resp := controllers.UpdatePFP(id, newPFP.Url)
	if resp == "Could not find account" {
		http.Error(w, "Could not find account", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully updated profile picture"))
}
func GetFeed(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Could not get session", http.StatusBadRequest)
		return
	}
	isAuth := session.Values["Authorized"]
	if isAuth != true {
		http.Error(w, "User not authorized", http.StatusBadRequest)
		return
	}
	resp := controllers.GetFeed(session.Values["userID"].(string))
	if len(resp) == 0 {
		http.Error(w, "Unable to get feed", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func IsStringEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}
