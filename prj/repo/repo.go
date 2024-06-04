package repo

import (
	"filmlib/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	GetUser(username, password string) (model.User, error)
	CreateUser(model.User) (int, error)
	AddToFavorites(userId, movId primitive.ObjectID) error
	RemoveFavorite(userId, movId primitive.ObjectID) error
	UserFavorites(userId primitive.ObjectID) ([]primitive.ObjectID, error)
}
