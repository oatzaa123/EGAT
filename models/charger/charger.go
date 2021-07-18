package charger

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

type Chargers struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Station_id        primitive.ObjectID `bson:"station_id,omitempty" json:"station_id"`
	Charger_no        int32              `bson:"charger_no,omitempty" json:"charger_no"`
	Kwh_total         int32              `bson:"kwh_total,omitempty" json:"kwh_total"`
	Max_kw_used       int32              `bson:"max_kw_used,omitempty" json:"max_kw_used"`
	Construction_date string             `bson:"construction_date,omitempty" json:"construction_date"`
	Cod_date          string             `bson:"cod_date,omitempty" json:"cod_date"`
	Transformer_id    primitive.ObjectID `bson:"transformer_id,omitempty" json:"transformer_id"`
	Charger_model_id  primitive.ObjectID `bson:"charger_model_id,omitempty" json:"charger_model_id"`
	Connector_status  []connector_status `bson:"connector_status,omitempty" json:"connector_status"`
	Commuication      commuication       `bson:"commuication,omitempty" json:"commuication"`
}

type connector_status struct {
	Connector_id   int32  `bson:"connector_id" json:"connector_id"`
	Max_current    int32  `bson:"max_current" json:"max_current"`
	Current_status string `bson:"current_status" json:"current_status"`
	Last_update    string `bson:"last_update" json:"last_update"`
}

type commuication struct {
	Mac_address string `bson:"mac_address" json:"mac_address"`
	Ip_address  string `bson:"ip_address" json:"ip_address"`
	Subnet      string `bson:"subnet" json:"subnet"`
	Gateway     string `bson:"gateway" json:"gateway"`
}

func (g *Chargers) GetChargers() ([]map[string]interface{}, error) {
	collection := database.MongoClient.Database(keys.Database).Collection("chargers")
	aggregateStage := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "stations",
				"localField":   "station_id",
				"foreignField": "_id",
				"as":           "station",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "networks",
				"localField":   "station.network",
				"foreignField": "_id",
				"as":           "network_data",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "charger_models",
				"localField":   "charger_model_id",
				"foreignField": "_id",
				"as":           "charger_model",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "transformers",
				"localField":   "transformer_id",
				"foreignField": "_id",
				"as":           "transformer",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$station",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$network_data",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$charger_model",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$transformer",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$project": bson.M{
				"charger_no":        1,
				"kwh_total":         1,
				"max_kw_used":       1,
				"construction_date": 1,
				"cod_date":          1,
				"transformer":       1,
				"connector_status":  1,
				"station": bson.M{
					"address":            1,
					"area":               1,
					"cod_date":           1,
					"construction_date":  1,
					"contact":            1,
					"contact_person":     1,
					"contact_tel":        1,
					"coordinate":         1,
					"landmark":           1,
					"name_en":            1,
					"name_th":            1,
					"open_houer_weekday": 1,
					"open_houer_weekend": 1,
					"parking_lots":       1,
					"network":            "$network_data",
				},
				"charger_model": 1,
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
