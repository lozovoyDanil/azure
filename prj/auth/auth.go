package auth

import (
	"context"
	"crypto/sha1"
	"errors"
	"filmlib/model"
	"filmlib/repo"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/golang-jwt/jwt"
)

const (
	salt       = "euv45675kdfjd458dhg43"
	signingKey = "dfjh47ty34hfd89wofdhf"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId  string `json:"user_id"`
	IsAdmin bool   `json:"isAdmin"`
}

type AuthService struct {
	repo repo.Repository
}

func NewService(repo repo.Repository) Service {
	return &AuthService{repo: repo}
}

func (s *AuthService) Healthz(_ context.Context) error {
	//TODO: add actual check
	return nil
}

func (s *AuthService) CreateUser(user model.User) (int, error) {
	user.Password = s.generatePassHash(user.Password)

	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, pass string) (string, error) {
	user, err := s.repo.GetUser(username, s.generatePassHash(pass))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID.Hex(),
		user.IsAdmin,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(token string) (model.Identity, error) {
	var u model.Identity

	userToken, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid jwt signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return u, err
	}

	claims, ok := userToken.Claims.(*tokenClaims)
	if !ok {
		return u, errors.New("wrong claims type")
	}

	u.Id = claims.UserId
	u.IsAdmin = claims.IsAdmin

	return u, nil
}

func (s *AuthService) AddToFavorites(userId, movId string) error {
	userObjectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	movObjectID, err := primitive.ObjectIDFromHex(movId)
	if err != nil {
		return err
	}

	return s.repo.AddToFavorites(userObjectID, movObjectID)
}

func (s *AuthService) RemoveFavorite(userId, movId string) error {
	userObjectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	movObjectID, err := primitive.ObjectIDFromHex(movId)
	if err != nil {
		return err
	}

	return s.repo.RemoveFavorite(userObjectID, movObjectID)
}

func (s *AuthService) UserFavorites(userId string) ([]primitive.ObjectID, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	return s.repo.UserFavorites(userObjectID)
}

func (s *AuthService) generatePassHash(pass string) string {
	hash := sha1.New()
	hash.Write([]byte(pass))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
