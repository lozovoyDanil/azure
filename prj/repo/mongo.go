package repo

import (
	"context"
	"errors"
	"filmlib/model"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

type AuthMongo struct {
	collection *mongo.Collection
}

func NewMongoRepo(uri string) (*AuthMongo, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	collection := client.Database("auth").Collection("users")
	return &AuthMongo{collection: collection}, nil
}

func (r *AuthMongo) MongoShutdown() {
	err := r.collection.Database().Client().Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func (r *AuthMongo) CreateUser(user model.User) (int, error) {
	user.FavoriteMovies = []primitive.ObjectID{}

	logger.Println(user)

	result, err := r.collection.InsertOne(context.Background(), user)
	if err != nil {
		return 0, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return 0, errors.New("failed to convert InsertedID to ObjectID")
	}
	return int(id.Timestamp().Unix()), nil
}

func (r *AuthMongo) GetUser(username, password string) (model.User, error) {
	var user model.User
	ctx := context.Background()
	filter := bson.M{"username": username, "password": password}

	err := r.collection.FindOne(ctx, filter).Decode(&user)

	return user, err
}

func (r *AuthMongo) AddToFavorites(userId, movId primitive.ObjectID) error {
	ctx := context.Background()
	filter := bson.M{"_id": userId}
	update := bson.M{"$addToSet": bson.M{"favoriteMovies": movId}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("no documents were updated, result.ModifiedCount = 0")
	}

	return nil
}

func (r *AuthMongo) RemoveFavorite(userId, movId primitive.ObjectID) error {
	ctx := context.Background()
	filter := bson.M{"_id": userId}
	update := bson.M{"$pull": bson.M{"favoriteMovies": movId}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("no documents were updated, result.ModifiedCount = 0")
	}

	return nil
}

func (r *AuthMongo) UserFavorites(userId primitive.ObjectID) ([]primitive.ObjectID, error) {
	ctx := context.Background()
	filter := bson.M{"_id": userId}
	projection := bson.M{"favoriteMovies": 1}

	var result struct {
		FavoriteMovies []primitive.ObjectID `bson:"favoriteMovies"`
	}

	err := r.collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("no documents found")
		}
		return nil, err
	}

	return result.FavoriteMovies, nil
}
