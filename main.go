package main

import (
	"log"
	"musicSharingAPp/db"
	"musicSharingAPp/router"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func init() {
	err := db.StartDB()
	if err != nil {
		log.Fatal("Could not connect to database")
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/api/v1/user/signup", router.SignUp)
	r.Delete("/api/v1/user/delete", router.DeleteAcc)
	r.Get("/api/v1/user", router.GetUserInfo)
	r.Get("/api/v1/user/:id", router.GetUserInfoByID)
	r.Get("/api/v1/user/all", router.GetAllUsersInfo)
	r.Post("/api/v1/user/login", router.UserLogin)
	r.Put("/api/v1/user/follow/{FollowID}", router.FollowUser)
	r.Put("/api/v1/user/unfollow/{UnfollowID}", router.UnfollowUser)

	r.Post("/api/v1/post/create", router.CreatePost)
	r.Get("/api/v1/posts", router.GetCurrentUserPosts)
	r.Get("/api/v1/posts/{id}", router.GetPostsViaUserID)
	r.Delete("/api/v1/posts/delete/{id}", router.DeletePost)
	r.Put("/api/v1/posts/like/{id}", router.LikePost)
	r.Put("/api/v1/posts/unlike/{id}", router.UnlikePost)

	http.ListenAndServe(":8080", r)
}
