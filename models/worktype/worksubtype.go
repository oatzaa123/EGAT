package worktype

import (
	"fmt"
	"go-callcenter/database"
	"go-callcenter/keys"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkSubType struct {
	ID                      primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty" form:"-"`
	WorkTypeID              primitive.ObjectID   `json:"work_type_id" bson:"work_type_id" form:"-"`
	CallCenterGroups        []primitive.ObjectID `json:"call_center_groups" bson:"call_center_groups" form:"-"`
	OtherGroups             []primitive.ObjectID `json:"other_groups" bson:"other_groups" form:"-"`
	Name                    string               `json:"name" bson:"name" form:"name"`
	ServiceLevel            int                  `json:"service_level" bson:"service_level" form:"service_level"`
	BasicTroubleshooting    string               `json:"basic_troubleshooting" bson:"basic_troubleshooting" form:"basic_troubleshooting"`
	AdvancedTroubleshooting string               `json:"advanced_troubleshooting" bson:"advanced_troubleshooting" form:"advanced_troubleshooting"`
	Active                  bool                 `json:"active" bson:"active" form:"active"`
	LastUpdated             time.Time            `json:"last_updated" bson:"last_updated" form:"last_updated"`
}

type EditWorkSubType struct {
	ID                      primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty" form:"-"`
	WorkTypeID              primitive.ObjectID   `json:"work_type_id,omitempty" bson:"work_type_id,omitempty" form:"-,omitempty"`
	CallCenterGroups        []primitive.ObjectID `json:"call_center_groups,omitempty" bson:"call_center_groups,omitempty" form:"-,omitempty "`
	OtherGroups             []primitive.ObjectID `json:"other_groups,omitempty" bson:"other_groups,omitempty" form:"-,omitempty"`
	Name                    string               `json:"name,omitempty" bson:"name,omitempty" form:"name"`
	ServiceLevel            int                  `json:"service_level,omitempty" bson:"service_level,omitempty" form:"service_level"`
	BasicTroubleshooting    string               `json:"basic_troubleshooting,omitempty" bson:"basic_troubleshooting,omitempty" form:"basic_troubleshooting"`
	AdvancedTroubleshooting string               `json:"advanced_troubleshooting,omitempty" bson:"advanced_troubleshooting,omitempty" form:"advanced_troubleshooting"`
	Active                  bool                 `json:"active,omitempty" bson:"active,omitempty" form:"active"`
	LastUpdated             time.Time            `json:"last_updated,omitempty" bson:"last_updated,omitempty" form:"last_updated"`
}

func (g *WorkSubType) CreateWorkSubType() error {
	collection := database.MongoClient.Database(keys.Database).Collection("WorkSubType")
	if _, err := collection.InsertOne(ctx, g); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (g *WorkSubType) GetWorkSubTypeAll() ([]WorkSubType, error) {
	collection := database.MongoClient.Database(keys.Database).Collection("WorkSubType")
	var workSubType []WorkSubType
	filter := bson.M{
		"active": true,
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err := cursor.All(ctx, &workSubType); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return workSubType, nil
}

func (g *EditWorkSubType) PutWorkSubType() error {
	collection := database.MongoClient.Database(keys.Database).Collection("WorkSubType")
	filter := bson.M{
		"_id":    g.ID,
		"active": g.Active,
	}
	update := bson.M{
		"$set": g,
	}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (g *EditWorkSubType) DeleteWorkSubType() error {
	collection := database.MongoClient.Database(keys.Database).Collection("WorkSubType")
	filter := bson.M{
		"_id":    g.ID,
		"active": true,
	}
	update := bson.M{
		"$set": bson.M{
			"active":       g.Active,
			"last_updated": g.LastUpdated,
		},
	}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
