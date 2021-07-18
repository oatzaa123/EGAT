package users

import (
	"fmt"
	"go-callcenter/models/users"
	"strings"
	"time"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func editProfileUsers(c echo.Context) error {
	var (
		payloadUser users.Users
		payload     users.EditProfileUsers
		t           = time.Now()
	)
	if err := c.Bind(&payload); err != nil {
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Can not bind data.",
		})
	}
	payload.LastUpdated = t
	payload.Active = true

	// recieve callcenteruser id string convert to primitive objectID
	_ID, err := primitive.ObjectIDFromHex(c.FormValue("_id"))
	if err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}
	payload.ID = _ID
	//check userID
	if payload.ID.Hex() == "" {
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Don't have user ID.",
		})
	}

	// recieve division id string covert to primitive objectid
	// if c.FormValue("division_id") != "" {
	// 	division, err := primitive.ObjectIDFromHex(c.FormValue("division"))
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return c.JSON(500, echo.Map{
	// 			"status": false,
	// 			"result": err.Error(),
	// 		})
	// 	}
	// 	payload.Division = divisionID
	// }
	payload.Division = c.FormValue("division")

	// recieve call_center_group_id string covert to primitive objectid
	if c.FormValue("call_center_group_id") != "" {
		callCenterGroupID, err := primitive.ObjectIDFromHex(c.FormValue("call_center_group_id"))
		if err != nil {
			fmt.Println(err)
			return c.JSON(500, echo.Map{
				"status": false,
				"result": err.Error(),
			})
		}
		payload.UserGroupID = callCenterGroupID
	}
	// recieve date_of_birth
	if c.FormValue("date_of_birth") != "" {
		layOutDate := "2006-01-02"
		dateOfBirth, err := time.Parse(layOutDate, c.FormValue("date_of_birth"))
		if err != nil {
			return c.JSON(500, echo.Map{
				"status": false,
				"result": err.Error(),
			})
		}
		payload.DateOfBirth = dateOfBirth
	}

	// Check NewUser is already exists.
	payloadUser.UserName = payload.UserName
	fmt.Println(payload.UserName == "")
	if payload.UserName != "" {
		if err := payloadUser.FindUserName(); err != nil {
			return c.JSON(500, echo.Map{
				"status": false,
				"result": err.Error(),
			})
		}
	}

	// check gender male female to database
	if payload.Gender != "" {
		genderToUpper := strings.ToUpper(payload.Gender)
		compareGender := genderToUpper == "MALE" || genderToUpper == "FEMALE"
		if !compareGender {
			return c.JSON(500, echo.Map{
				"status":  false,
				"message": "Please enter gender MALE or FEMALE.",
			})
		}
		payload.Gender = genderToUpper
	}

	// // hash password
	// payloadUser.Password = payload.Password
	// hashPassword, err := payloadUser.HashPassword()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return c.JSON(500, echo.Map{
	// 		"status": false,
	// 		"result": err.Error(),
	// 	})
	// }
	// payload.Password = hashPassword

	// upload image
	file, err := c.FormFile("image")
	if file != nil {
		if err != nil {
			return c.JSON(500, echo.Map{
				"status": false,
				"result": err.Error(),
			})
		}
		src, err := file.Open()
		if err != nil {
			return c.JSON(500, echo.Map{
				"status": false,
				"result": err.Error(),
			})
		}
		defer src.Close()
		imagePath, err := users.UploadImage(src)
		if err != nil {
			fmt.Println(err)
			return c.JSON(500, echo.Map{
				"status": false,
				"result": err.Error(),
			})
		}
		payload.Image = imagePath
	}

	// update to database
	if err := payload.EditProfileUser(); err != nil {
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"status":  true,
		"message": "Edit profile successful.",
	})
}
