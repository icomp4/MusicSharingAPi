package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string
	Password       string
	FavGenres      []Genre `gorm:"many2many:user_genres;"`
	Following      []User  `gorm:"many2many:user_following;"`
	Followers      []User  `gorm:"many2many:user_followers;"`
	FollowerCount  int
	FollowingCount int
	PostCount      int
	Blocked        []User     `gorm:"many2many:user_blocked;"`
	Playlists      []Playlist `gorm:"foreignKey:UserID;"`
	Posts          []Post     `gorm:"foreignKey:UserID"`
	LikedPosts     []Post     `gorm:"many2many:user_likes_posts;"`
	PfpURL        string
}

type Genre struct {
	gorm.Model
	Name string
}

type Playlist struct {
	gorm.Model
	Name   string
	UserID uint
	User   User
	Songs  []Song `gorm:"many2many:playlist_songs;"`
}

type Song struct {
	gorm.Model
	Name     string
	Artist   string
	GenreID  uint
	Genre    Genre
	Duration uint // Duration in seconds
}
type Post struct {
	gorm.Model
	Title   string
	Likes   int
	Content string
	UserID  uint
	Author  string  `gorm:"column:author"`
	Genres  []Genre `gorm:"many2many:post_genres;"`
}
