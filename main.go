package main

import (
	"log"
	"musicSharingAPp/db"
	"musicSharingAPp/router"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func init() {
	err := db.StartDB()
	if err != nil {
		log.Fatal("Could not connect to database")
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true, // Set AllowCredentials to true to allow credentials (cookies) in requests
		MaxAge:           300,
	}))
	r.Use(middleware.Logger)
	r.Post("/api/v1/user/signup", router.HandleSignUp)
	r.Delete("/api/v1/user/delete", router.DeleteAcc)
	r.Get("/api/v1/user", router.GetUserInfo)
	r.Get("/api/v1/user/:id", router.GetUserInfoByID)
	r.Get("/api/v1/user/all", router.GetAllUsersInfo)
	r.Post("/api/v1/user/login", router.UserLogin)
	r.Put("/api/v1/user/follow/{FollowID}", router.FollowUser)
	r.Put("/api/v1/user/unfollow/{UnfollowID}", router.UnfollowUser)
	r.Put("/api/v1/user/update/password", router.UpdatePassword)
	r.Put("/api/v1/user/update/pfp", router.UpdatePFP)

	r.Post("/api/v1/post/create", router.CreatePost)
	r.Get("/api/v1/post", router.GetCurrentUserPosts)
	r.Get("/api/v1/post/{id}", router.GetPostsViaUserID)
	r.Delete("/api/v1/post/delete/{id}", router.DeletePost)
	r.Put("/api/v1/post/like/{id}", router.LikePost)
	r.Put("/api/v1/post/unlike/{id}", router.UnlikePost)

	r.Get("/api/v1/feed", router.GetFeed)

	http.ListenAndServe(":8080", r)
}
