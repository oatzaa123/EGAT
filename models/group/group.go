package division

import (
	"context"
	"errors"
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

// UserGroups comment
type UserGroups struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id" form:"_id" query:"_id"`
	Name        string             `bson:"name" json:"name" form:"name" query:"name"`
	Description string             `bson:"description" json:"description" form:"description" query:"description"`
	Active      bool               `bson:"active" json:"active" form:"active" query:"active"`
	LastUpdated time.Time          `bson:"last_updated" json:"last_updated" form:"last_updated" query:"last_updated"`
}

// GetUserGroups comment
func (g *UserGroups) GetCallCenterGroups() ([]UserGroups, error) {
	collection := database.MongoClient.Database(keys.Database).Collection("user_groups")
	filter := bson.M{"active": true}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	result := []UserGroups{}
	if err := cursor.All(ctx, &result); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}

// GetMembersInUserGroups comment
func (g *UserGroups) GetMembersInUserGroups() ([]map[string]interface{}, error) {
	collection := database.MongoClient.Database(keys.Database).Collection("user_groups")
	lookupStage := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "_id",
				"foreignField": "user_group_id",
				"as":           "members",
			},
		}}
	cursor, err := collection.Aggregate(ctx, lookupStage)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	result := []map[string]interface{}{}
	if err := cursor.All(ctx, &result); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return result, nil
}

// GetMembersForChangeUserGroups comment
func (g *UserGroups) GetMembersForChangeUserGroups() ([]map[string]interface{}, error) {
	collection := database.MongoClient.Database(keys.Database).Collection("user_groups")
	aggregateStage := []bson.M{
		{
			"$match": bson.M{
				"active":        true,
				"user_group_id": bson.M{"$ne": g.ID},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "user_groups",
				"localField":   "user_group_id",
				"foreignField": "_id",
				"as":           "groups",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$groups",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{"$sort": bson.M{"_id": -1}},
		{
			"$project": bson.M{
				"_id":           1,
				"image":         1,
				"first_name":    1,
				"last_name":     1,
				"date_of_birth": 1,
				"gender":        1,
				"username":      1,
				"last_updated":  1,
				"mobile":        1,
				"division":      1,
				"groups":        1,
				"staff_id":      1,
			},
		},
	}
	cursor, err := collection.Aggregate(ctx, aggregateStage)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	result := []map[string]interface{}{}
	if err := cursor.All(ctx, &result); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}

// MembersChangeGroups comment
func (g *UserGroups) MembersChangeGroups(members []string) error {
	collection := database.MongoClient.Database(keys.Database).Collection("users")
	var objectIDMembers []primitive.ObjectID
	for _, v := range members {
		objectID, err := primitive.ObjectIDFromHex(v)
		if err == nil {
			objectIDMembers = append(objectIDMembers, objectID)
		}
	}
	filter := bson.M{
		"active": true,
		"_id":    bson.M{"$in": objectIDMembers},
	}
	update := bson.M{
		"$set": bson.M{"user_group_id": g.ID},
	}
	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		newError := errors.New("Cannot Updated")
		return newError
	}
	return nil
}

// UpdateUserGroups comment
func (g *UserGroups) UpdateUserGroups() error {
	collection := database.MongoClient.Database(keys.Database).Collection("user_groups")
	filter := bson.M{"_id": g.ID}
	update := bson.M{
		"$set": bson.M{
			"name":         g.Name,
			"description":  g.Description,
			"last_updated": g.LastUpdated,
		},
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// GetGroupDetail comment
func (g *UserGroups) GetGroupDetail() ([]map[string]interface{}, error) {
	collection := database.MongoClient.Database(keys.Database).Collection("user_groups")
	aggregateStage := []bson.M{
		{
			"$match": bson.M{
				"_id": g.ID,
			},
		},
		{
			"$lookup": bson.M{
				"from": "users",
				"let": bson.M{
					"group_id": "$_id"},
				"pipeline": []bson.M{
					{"$match": bson.M{
						"active": true,
						"$expr":  bson.M{"$eq": bson.A{"$user_group_id", "$$group_id"}},
					}},
				},
				"as": "user",
			},
		},
		{
			"$project": bson.M{
				"_id":               1,
				"name":              1,
				"description":       1,
				"default":           1,
				"user._id":          1,
				"user.image":        1,
				"user.first_name":   1,
				"user.last_name":    1,
				"user.active":       1,
				"user.last_updated": 1,
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, aggregateStage)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	result := []map[string]interface{}{}
	if err := cursor.All(ctx, &result); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}
