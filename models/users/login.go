package users

import (
	"errors"
	"fmt"
	"go-callcenter/database"
	"go-callcenter/keys"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Login input struct
type Login struct {
	Username string `json:"username" bson:"username" form:"username"`
	Password string `json:"password" bson:"password" form:"password"`
}

// LoginHistories struct
type LoginHistories struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" form:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id" form:"user_id"`
	IPAddress   string             `json:"ip_address" bson:"ip_address" form:"ip_address"`
	Device      string             `json:"device" bson:"device" form:"device"`
	LastUpdated time.Time          `json:"last_updated" bson:"last_updated" form:"last_updated"`
}

// FindUsername comment
func (g *Users) FindUsername(username string) error {
	collection := database.MongoClient.Database(keys.Database).Collection("users")
	filter := bson.M{
		"username": username,
		"active":   true,
	}
	if err := collection.FindOne(ctx, filter).Decode(&g); err != nil {
		fmt.Println(err)
		return errors.New("Don't have username")
	}
	return nil
}

// CreateLoginHistories comment
func (g *LoginHistories) CreateLoginHistories() error {
	collection := database.MongoClient.Database(keys.Database).Collection("login_histories")
	if _, err := collection.InsertOne(ctx, g); err != nil {
		fmt.Println(err)
		return errors.New("Can't create login history list.")
	}
	return nil
}

// GetLoginHistories comment
func (g *LoginHistories) GetLoginHistories() ([]LoginHistories, error) {
	collection := database.MongoClient.Database(keys.Database).Collection("login_histories")
	filter := bson.M{
		"user_id": g.ID,
	}
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"_id": -1})
	// findOptions.SetLimit(5)

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		fmt.Println(err)
	}
	var loginUsersHistories []LoginHistories
	if err := cursor.All(ctx, &loginUsersHistories); err != nil {
		fmt.Println(err)
	}
	return loginUsersHistories, err
}
