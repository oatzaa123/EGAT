package announcement

import (
	"context"
	"errors"
	"fmt"
	"go-callcenter/database"
	"go-callcenter/keys"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx = context.Background()
)

// Announcement comment
type Announcement struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id" form:"_id" query:"_id"`
	Detail      string             `bson:"detail" json:"detail" form:"detail" query:"detail"`
	Type        string             `bson:"type" json:"type" form:"type" query:"type"`
	Order       int                `bson:"order" json:"order"`
	Active      bool               `bson:"active" json:"active"`
	LastUpdated time.Time          `bson:"last_updated" json:"last_updated"`
}

// GetAnnouncement comment
func (g *Announcement) GetAnnouncement() (Announcement, []Announcement, error) {
	fmt.Println("Function GetAnnouncement")
	collection := database.MongoClient.Database(keys.Database).Collection("Announcement")
	var announcement []Announcement
	filter := bson.M{"active": true}
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"order": 1})
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		fmt.Println(err)
	}
	err = cursor.All(ctx, &announcement)
	if err != nil {
		fmt.Println(err)
	}
	var importanceAnnouncement Announcement
	for i, v := range announcement {
		if announcement[i].Type == "IMPORTANCE" {
			importanceAnnouncement = v
			announcement = append(announcement[:i], announcement[i+1:]...)
			break
		}
	}
	return importanceAnnouncement, announcement, err
}

// AddAnnouncement comment
func (g *Announcement) AddAnnouncement() error {
	fmt.Println("Function AddAnnouncement")
	collection := database.MongoClient.Database(keys.Database).Collection("Announcement")
	if _, err := collection.InsertOne(ctx, g); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// UpdateAnnouncement comment
func (g *Announcement) UpdateAnnouncement() error {
	fmt.Println("Function UpdateAnnouncement")
	collection := database.MongoClient.Database(keys.Database).Collection("Announcement")
	filter := bson.M{"_id": g.ID}
	update := bson.M{
		"$set": bson.M{
			"detail":       g.Detail,
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

// CountOtherAnnouncement : count number of announcement type other not greater than 10
func (g *Announcement) CountOtherAnnouncement() error {
	fmt.Println("Function CountOtherAnnouncement")
	collection := database.MongoClient.Database(keys.Database).Collection("Announcement")
	filter := bson.M{"active": true, "type": "OTHERS"}
	numberOfAnnouncement, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		fmt.Println(err)
	}
	if numberOfAnnouncement >= 10 {
		return errors.New("Limit announcement is 10")
	}
	return nil
}

// ArrangeOrderType : เรียงลำดับ  order ให้เพิ่มขึ้น 1
func (g *Announcement) ArrangeOrderType() error {
	fmt.Println("Function ArrangeOrderType")
	collection := database.MongoClient.Database(keys.Database).Collection("Announcement")
	filter := bson.M{"active": true, "type": "OTHERS"}
	update := bson.M{
		"$inc": bson.M{
			"order": 1,
		},
	}
	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// DeleteAnnouncement : update status active >> false and order >> -1
func (g *Announcement) DeleteAnnouncement() error {
	fmt.Println("Function DeleteAnnouncement")
	collection := database.MongoClient.Database(keys.Database).Collection("Announcement")
	filter := bson.M{
		"_id":  g.ID,
		"type": "OTHERS",
	}
	update := bson.M{
		"$set": bson.M{
			"active": false,
			"order":  -1,
		},
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// ReplaceOrderAfterDeleteAnm : change order replace announcement that delete
func (g *Announcement) ReplaceOrderAfterDeleteAnm() error {
	fmt.Println("Function ReplaceOrderAfterDeleteAnm")
	collection := database.MongoClient.Database(keys.Database).Collection("Announcement")
	filter := bson.M{
		"type":   "OTHERS",
		"active": true,
		"order":  bson.M{"$gt": g.Order},
	}
	update := bson.M{
		"$inc": bson.M{
			"order": -1,
		},
	}
	if _, err := collection.UpdateMany(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

// FindAnnouncementOrderByID commment
func (g *Announcement) FindAnnouncementOrderByID() (currentOrder int, err error) {
	fmt.Println("Function FindAnnouncementOrderByID")
	fmt.Println("Announcement ID :", g.ID)
	collection := database.MongoClient.Database(keys.Database).Collection("Announcement")
	currentAnnoucement := g
	filter := bson.M{"_id": g.ID}
	if err = collection.FindOne(ctx, filter).Decode(&currentAnnoucement); err != nil {
		fmt.Println("err at FindAnnouncementOrderByID function :", err)
	}
	return currentAnnoucement.Order, err
}

// FindAnnouncementIDByOrder commment
func (g *Announcement) FindAnnouncementIDByOrder(swapOrder int) (currentOrder primitive.ObjectID, err error) {
	fmt.Println("Function FindAnnouncementIDByOrder")
	fmt.Println("Announcement Order :", g.Order)
	collection := database.MongoClient.Database(keys.Database).Collection("Announcement")
	swapAnm := g
	filter := bson.M{
		"order":  swapOrder,
		"active": true,
	}
	if err = collection.FindOne(ctx, filter).Decode(&swapAnm); err != nil {
		fmt.Println("err at FindAnnouncementIDByOrder function :", err)
	}
	return swapAnm.ID, err
}

// RearrangeOrder comment
func (g *Announcement) RearrangeOrder(amnID primitive.ObjectID, order int) error {
	fmt.Println("Function RearrangeOrder")
	collection := database.MongoClient.Database(keys.Database).Collection("Announcement")
	filterStage := bson.M{"_id": amnID}
	updateStage := bson.M{
		"$set": bson.M{
			"order": order,
		},
	}
	_, err := collection.UpdateOne(ctx, filterStage, updateStage)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
