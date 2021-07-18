package users

import (
	"fmt"
	"go-callcenter/database"
	"go-callcenter/keys"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EditProfileUsers comment
type EditProfileUsers struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id" form:"-" query:"_id"`
	Division    string             `bson:"division,omitempty" json:"division" form:"-" query:"division" `
	UserGroupID primitive.ObjectID `bson:"user_group_id,omitempty" json:"user_group_id" form:"-" query:"user_group_id" validate:"required"`
	Image       string             `bson:"image,omitempty" json:"image" form:"-" query:"image"`
	FirstName   string             `bson:"first_name,omitempty" json:"first_name"  form:"first_name" query:"first_name" `
	LastName    string             `bson:"last_name,omitempty" json:"last_name" form:"last_name" query:"last_name" `
	DateOfBirth time.Time          `bson:"date_of_birth,omitempty" json:"date_of_birth" form:"-" query:"date_of_birth" `
	Gender      string             `bson:"gender,omitempty" json:"gender" form:"gender" query:"gender" `
	Mobile      string             `bson:"mobile,omitempty" json:"mobile" form:"mobile" query:"mobile" `
	UserName    string             `bson:"username,omitempty" json:"username" form:"username" query:"username" `
	// Password          string             `bson:"password,omitempty" json:"password" form:"password" query:"password" `
	CreateDate  time.Time `bson:"created_date,omitempty" json:"created_date" form:"created_date" query:"created_date" `
	Active      bool      `bson:"active,omitempty" json:"active" form:"active" query:"active" `
	LastUpdated time.Time `bson:"last_updated,omitempty" json:"last_updated" form:"last_updated" query:"last_updated" `
	StaffID     string             `bson:"staff_id" json:"staff_id" form:"staff_id" query:"staff_id"`

}

// EditProfileUser comment
func (u *EditProfileUsers) EditProfileUser() error {
	collection := database.MongoClient.Database(keys.Database).Collection("users")
	filter := bson.M{
		"active": true,
		"_id":    u.ID,
	}
	update := bson.M{
		"$set": u,
	}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	return nil
}
