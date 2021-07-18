package database

import (
	"context"
	"fmt"
	"time"

	"go-callcenter/keys"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient comment
var (
	MongoClient *mongo.Client
)

// InitDatabase connect to mongo database
func InitDatabase() {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
	)
	defer cancel()
	fmt.Println("Init Database")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(keys.URI))
	if err != nil {
		fmt.Println("DB Error :", err)
	}
	MongoClient = client
}
