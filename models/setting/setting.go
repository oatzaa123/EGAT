package setting

import (
	"context"
	"fmt"
	"go-callcenter/database"
	"go-callcenter/keys"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ctx = context.Background()
)

// Setting comment
type Setting struct {
	ID               primitive.ObjectID `bson:"_id" json:"_id" form:"_id" query:"_id"`
	KwhToTHB      float64             `bson:"kwh_to_thb" json:"kwh_to_thb" form:"kwh_to_thb" query:"kwh_to_thb"`
	LastUpdated      time.Time          `bson:"last_updated" json:"last_updated" form:"last_updated" query:"last_updated"`	
}

// GetSetting comment
func (g *Setting) GetSetting() error {
	collection := database.MongoClient.Database(keys.Database).Collection("setting")
	filter := bson.M{}
	if err := collection.FindOne(ctx, filter).Decode(&g); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// UpdateSetting comment
func (g *Setting) UpdateSetting() error {
	g.LastUpdated = time.Now()
	collection := database.MongoClient.Database(keys.Database).Collection("setting")
	filter := bson.M{"_id": g.ID}
	update := bson.M{
		"$set": g,
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
