package worktype

import (
	"context"
	"fmt"
	"go-callcenter/database"
	"go-callcenter/keys"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkType struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" form:"_id,omitempty" query:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty" form:"name"`
	Active      bool               `json:"active,omitempty" bsom:"active,omitempty" form:"active,omitempty"`
	LastUpdated time.Time          `json:"last_updated,omitempty" bson:"last_updated,omitempty" form:"last_updated,omitempty"`
}

var (
	ctx = context.Background()
)

// GetWorkType comment
func (g *WorkType) GetWorkType() ([]map[string]interface{}, error) {
	fmt.Println("Function GetWorkType")
	collection := database.MongoClient.Database(keys.Database).Collection("WorkType")
	aggregateStage := []bson.M{
		{"$match": bson.M{"active": true}},
		{
			"$lookup": bson.M{
				"from": "WorkSubType",
				"let":  bson.M{"work_type_id": "$_id"},
				"pipeline": []bson.M{
					{"$match": bson.M{
						"active": true,
						"$expr":  bson.M{"$eq": bson.A{"$work_type_id", "$$work_type_id"}},
					}},
					{"$sort": bson.M{"_id": -1}},
				},
				"as": "work_sub_type",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$work_sub_type",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from":         "CallCenterGroups",
				"localField":   "work_sub_type.call_center_groups",
				"foreignField": "_id",
				"as":           "work_sub_type.call_center_groups",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "OthersGroup",
				"localField":   "work_sub_type.other_groups",
				"foreignField": "_id",
				"as":           "work_sub_type.other_groups",
			},
		},
		{
			"$group": bson.M{
				"_id":           "$_id",
				"name":          bson.M{"$first": "$name"},
				"work_sub_type": bson.M{"$addToSet": "$work_sub_type"},
			},
		},
		{"$sort": bson.M{"_id": -1}},
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

// CreateWorkType comment
func (g *WorkType) CreateWorkType() error {
	collection := database.MongoClient.Database(keys.Database).Collection("WorkType")
	if _, err := collection.InsertOne(ctx, g); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

//GetAllWorkType comment
func (g *WorkType) GetWorkTypeAll() ([]WorkType, error) {
	collection := database.MongoClient.Database(keys.Database).Collection("WorkType")
	filter := bson.M{
		"active": true,
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var worktype []WorkType
	if err := cursor.All(ctx, &worktype); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return worktype, nil
}

//EditWorkType comment ---> not finish cannot input _id
func (g *WorkType) EditWorkType() error {
	// fmt.Println("InEditWorkType")
	// fmt.Println("InEditWorkType : ", g.ID)
	// fmt.Println("InEditWorkType : ", g.Name)
	// fmt.Println("InEditWorkType : ", g.Active)
	collection := database.MongoClient.Database(keys.Database).Collection("WorkType")
	filter := bson.M{
		"_id": g.ID,
	}
	update := bson.M{
		"$set": bson.M{
			"name":         g.Name,
			"active":       g.Active,
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

func (g *WorkType) DeleteWorkType() error {
	collection := database.MongoClient.Database(keys.Database).Collection("WorkType")
	filter := bson.M{
		"_id": g.ID,
	}
	update := bson.M{
		"$set": bson.M{
			"active": false,
		},
	}
	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (g *WorkType) GetAllWorkTypeAndWorkSubType(serviceLevel int, name string) ([]map[string]interface{}, error) {
	collection := database.MongoClient.Database(keys.Database).Collection("WorkType")
	// collection := database.MongoClient.Database(keys.Database).Collection("WorkSubType")
	// lookupStage := []bson.M{
	// 	{
	// 		"$match": bson.M{
	// 			"$eq": bson.M{"$name": name},
	// 		},
	// 	},
	// }
	fmt.Println("name name : ", name)
	lookupStage := []bson.M{
		{
			"$lookup": bson.M{
				"from": "WorkSubType",
				"let": bson.M{
					"workTypeID": "$_id",
				},
				"as": "work-sub-type",
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": bson.A{
									bson.M{
										"$eq": bson.A{"$work_type_id", "$$workTypeID"},
									},
									bson.M{
										"$eq": bson.A{"$service_level", serviceLevel},
									},
									bson.M{
										"$regexMatch": bson.M{
											"input": "$name",
											"regex": name,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	cursor, err := collection.Aggregate(ctx, lookupStage)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var result []map[string]interface{}
	if err := cursor.All(ctx, &result); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}

func (g *WorkType) GetAllWorkTypeAndWorkSubTypeAll() ([]map[string]interface{}, error) {
	collection := database.MongoClient.Database(keys.Database).Collection("WorkType")
	lookupStage := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "WorkSubType",
				"localField":   "_id",
				"foreignField": "work_type_id",
				"as":           "work-sub-type",
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, lookupStage)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var result []map[string]interface{}
	if err := cursor.All(ctx, &result); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}
