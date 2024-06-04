package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty"`
	Username       string               `json:"username" validate:"required" bson:"username"`
	Password       string               `json:"password" validate:"required" bson:"password"`
	IsAdmin        bool                 `bson:"isAdmin"`
	FavoriteMovies []primitive.ObjectID `bson:"favoriteMovies"`
}

type Identity struct {
	Id      string `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
}
