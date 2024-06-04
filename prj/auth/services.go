package auth

import (
	"context"
	"filmlib/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (model.Identity, error)
	Healthz(c context.Context) error
	AddToFavorites(userId, movId string) error
	RemoveFavorite(userId, movId string) error
	UserFavorites(userId string) ([]primitive.ObjectID, error)
}
