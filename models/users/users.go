package users

import (
	"context"
	"errors"
	"fmt"
	"go-callcenter/database"
	"go-callcenter/keys"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Users comment
type Users struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id" form:"-" query:"_id"`
	Division    string             `bson:"division,omitempty" json:"division" form:"-" query:"division" `
	UserGroupID primitive.ObjectID `bson:"user_group_id,omitempty" json:"user_group_id" form:"-" query:"user_group_id" validate:"required"`
	Image       string             `bson:"image" json:"image" form:"-" query:"image"`
	FirstName   string             `bson:"first_name" json:"first_name"  form:"first_name" query:"first_name" validate:"required"`
	LastName    string             `bson:"last_name" json:"last_name" form:"last_name" query:"last_name" validate:"required"`
	DateOfBirth time.Time          `bson:"date_of_birth" json:"date_of_birth" form:"-" query:"date_of_birth" validate:"required"`
	Gender      string             `bson:"gender" json:"gender" form:"gender" query:"gender" validate:"required"`
	Mobile      string             `bson:"mobile" json:"mobile" form:"mobile" query:"mobile" validate:"required"`
	UserName    string             `bson:"username" json:"username" form:"username" query:"username" validate:"required"`
	Password    string             `bson:"password" json:"password" form:"password" query:"password" validate:"required"`
	CreateDate  time.Time          `bson:"created_date" json:"created_date" form:"created_date" query:"created_date" validate:"required"`
	Active      bool               `bson:"active" json:"active" form:"active" query:"active" validate:"required"`
	LastUpdated time.Time          `bson:"last_updated" json:"last_updated" form:"last_updated" query:"last_updated" validate:"required"`
	StaffID     string             `bson:"staff_id" json:"staff_id" form:"staff_id" query:"staff_id"`
}

var (
	ctx = context.Background()
)

// CreateUser comment
func (u *Users) CreateUser() error {
	collection := database.MongoClient.Database(keys.Database).Collection("users")
	if _, err := collection.InsertOne(ctx, u); err != nil {
		fmt.Println("model err", err)
		return err
	}
	return nil
}

// FindUserName for username is already exists.
func (u *Users) FindUserName() error {
	collection := database.MongoClient.Database(keys.Database).Collection("users")
	filter := bson.M{
		"username": u.UserName,
		"active":   true,
	}

	var userNameInDataBase Users
	if err := collection.FindOne(ctx, filter).Decode(&userNameInDataBase); err != nil {
		fmt.Println(err)
	}

	compareUserName := userNameInDataBase.UserName == u.UserName
	if compareUserName {
		newError := errors.New("Username is already exists")
		return newError
	}
	return nil
}

// HashPassword comment
func (u *Users) HashPassword() (hashPassword string, err error) {
	bytePassword := []byte(u.Password)
	result, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(result), nil
}

// UploadImage comment
func UploadImage(src io.Reader) (imagePath string, err error) {

	genFileName := time.Now().UnixNano()
	filename := strconv.Itoa(int(genFileName))

	dst, err := os.Create(keys.ImagePath + filename + ".jpg")
	if err != nil {
		return "", err
	}
	// fmt.Println("*** Path of image : ", keys.ImagePath+filename+".jpg")

	imagePath = keys.ImagePath + filename + ".jpg"

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return imagePath, err
}

// Validate require comment
func (u *Users) Validate() error {
	v := validator.New()
	err := v.Struct(u)
	if err != nil {
		var text string = "field that's not correct : "
		for _, e := range err.(validator.ValidationErrors) {
			fmt.Println(e.Field())
			text = text + e.Field() + ", "
			// println("length of empty filed", len(emptyFiled))
		}
		text = strings.TrimSuffix(text, ", ")
		return errors.New(text)
	}
	return nil
}

// GetProfileUsers comment
func (u *Users) GetProfileUsers() ([]map[string]interface{}, error) {
	collection := database.MongoClient.Database(keys.Database).Collection("users")
	aggregateStage := []bson.M{
		{
			"$match": bson.M{
				"":       u.ID,
				"active": true,
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

// GetUsers comment
func (u *Users) GetUsers() ([]map[string]interface{}, error) {
	collection := database.MongoClient.Database(keys.Database).Collection("users")
	aggregateStage := []bson.M{
		{
			"$match": bson.M{
				"active": true,
			},
		},
		{
			"$lookup": bson.M{
				"from":         "user_groups",
				"localField":   "user_groups",
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
		{"$sort": bson.M{"created_date": -1}},
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

// UpdateUserGroup comment
func (u *Users) UpdateUserGroup() error {
	collection := database.MongoClient.Database(keys.Database).Collection("users")
	filter := bson.M{"_id": u.ID}
	update := bson.M{
		"$set": bson.M{
			"user_group_id": u.UserGroupID,
			"last_updated":  u.LastUpdated,
		},
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
