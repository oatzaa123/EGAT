package stations

import (
	"context"
	"fmt"
	"go-callcenter/database"
	"go-callcenter/keys"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ctx = context.Background()
)

type Stations struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name_en            string             `bson:"name_en,omitempty" json:"name_en"`
	Name_th            string             `bson:"name_th,omitempty" json:"name_th"`
	Area               string             `bson:"area" json:"area"`
	Coordinate         coordinate         `bson:"coordinate" json:"coordinate"`
	Network            primitive.ObjectID `bson:"network" json:"network"`
	Address            string             `bson:"address" json:"address"`
	Landmark           string             `bson:"landmark" json:"landmark"`
	Contact            contact            `bson:"contact" json:"contact"`
	Construction_date  string             `bson:"construction_date" json:"construction_date"`
	Cod_date           string             `bson:"cod_date" json:"cod_date"`
	Open_houer_weekday string             `bson:"open_houer_weekday" json:"open_houer_weekday"`
	Open_houer_weekend string             `bson:"open_houer_weekend" json:"open_houer_weekend"`
	Contact_person     string             `bson:"contact_person" json:"contact_person"`
	Contact_tel        string             `bson:"contact_tel" json:"contact_tel"`
	Parking_lots       []parking_lots     `bson:"parking_lots" json:"parking_lots"`
}

type coordinate struct {
	Lattitude float64 `bson:"lattitude" json:"lattitude" form:"-" query:"lattitude"`
	Longitude float64 `bson:"longitude" json:"longitude" form:"-" query:"longitude"`
}

type parking_lots struct {
	ID           int32 `bson:"id" json:"id"`
	Connector_id int32 `bson:"connector_id" json:"connector_id"`
}

type contact struct {
	Tel       string `bson:"tel" json:"tel"`
	Fax       string `bson:"fax" json:"fax"`
	Email     string `bson:"email" json:"email"`
	Facebook  string `bson:"facebook" json:"facebook"`
	Instagram string `bson:"instagram" json:"instagram"`
	Twtter    string `bson:"twtter" json:"twtter"`
	Line      string `bson:"line" json:"line"`
}

func (g *Stations) GetStations() ([]map[string]interface{}, error) {
	collection := database.MongoClient.Database(keys.Database).Collection("stations")
	aggregateStage := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "networks",
				"localField":   "network",
				"foreignField": "_id",
				"as":           "data",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$data",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$project": bson.M{
				"name_en":            1,
				"name_th":            1,
				"area":               1,
				"coordinate":         1,
				"address":            1,
				"landmark":           1,
				"contact":            1,
				"construction_date":  1,
				"cod_date":           1,
				"open_houer_weekday": 1,
				"open_houer_weekend": 1,
				"contact_person":     1,
				"contact_tel":        1,
				"parking_lots":       1,
				"network":            "$data.name_th",
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
