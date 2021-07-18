package users

import (
	"fmt"
	"go-callcenter/models/users"
	"strings"
	"time"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PrivateRoute comment
func PrivateRoute(g *echo.Group) {
	g.POST("/register", register)
	g.POST("/user/check-username", checkUsername)
	g.GET("/user/profile/:user-id", getProfileUsers)
	g.GET("/user/login-histories/:user-id", getLoginHistories)
	g.PUT("/user/edit-profile", editProfileUsers)
	g.GET("/user/get", getUsers)
	g.PUT("/user/delete-from-group", deleteUserFromGroup)
}


// PublicRoute comment
func PublicRoute(g *echo.Group) {
	g.POST("/login", login)

}

func register(c echo.Context) error {
	fmt.Println("/register")
	var (
		payload   = users.Users{}
		t         = time.Now()
		imagePath = ""
	)

	if err := c.Bind(&payload); err != nil {
		fmt.Println("/register err bind payload :", err)
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}
	payload.CreateDate = t
	payload.LastUpdated = t
	payload.Active = true

	// recieve division id string covert to primitive objectid
	// divisionID, err := primitive.ObjectIDFromHex(c.FormValue("division_id"))
	// if err != nil {
	// 	fmt.Println("/register err division id :", err)
	// 	return c.JSON(500, echo.Map{
	// 		"status": false,
	// 		"result": err.Error(),
	// 	})
	// }
	// payload.DivisionID = divisionID
	payload.Division = c.FormValue("division")

	// recieve call_center_group_id string covert to primitive objectid
	userGroupID, err := primitive.ObjectIDFromHex(c.FormValue("user_group_id"))
	if err != nil {
		fmt.Println("/register err call center group id :", err)
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}
	payload.UserGroupID = userGroupID
	// recieve date_of_birth
	layOutDate := "2006-01-02"
	dateOfBirth, err := time.Parse(layOutDate, c.FormValue("date_of_birth"))
	if err != nil {
		fmt.Println("/register err date of birth :", err)
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}
	payload.DateOfBirth = dateOfBirth

	// check new user is already created
	if err := payload.FindUserName(); err != nil {
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	// check gender male female to database
	genderToUpper := strings.ToUpper(payload.Gender)
	compareGender := genderToUpper == "MALE" || genderToUpper == "FEMALE"
	if !compareGender {
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Not MALE or FEMALE",
		})
	}
	payload.Gender = genderToUpper

	// check fill data
	if err := payload.Validate(); err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	// hash password
	hashPassword, err := payload.HashPassword()
	if err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}
	payload.Password = hashPassword

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
		imagePathFromDB, err := users.UploadImage(src)
		if err != nil {
			fmt.Println(err)
			return c.JSON(500, echo.Map{
				"status": false,
				"result": err.Error(),
			})
		}
		imagePath = imagePathFromDB
	}

	payload.Image = imagePath

	// insert register to database
	if err := payload.CreateUser(); err != nil {
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	return c.JSON(201, echo.Map{
		"status": true,
	})
}

// check user name comment
func checkUsername(c echo.Context) error {
	var (
		payload users.Users
	)
	if err := c.Bind(&payload); err != nil {
		fmt.Println(err)
	}

	if err := payload.FindUserName(); err != nil {
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"status":  true,
		"message": "",
	})
}

func getProfileUsers(c echo.Context) error {
	var (
		payload users.Users
	)
	userID := c.Param("user-id")
	fmt.Println(userID)
	objectUsersID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Cannot convert to objectID",
		})
	}
	payload.ID = objectUsersID

	result, err := payload.GetProfileUsers()
	if err != nil {
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"status": true,
		"result": result[0],
	})
}

func getUsers(c echo.Context) error {
	var (
		payload users.Users
	)

	result, err := payload.GetUsers()
	if err != nil {
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"status": true,
		"result": result,
	})
}

func deleteUserFromGroup(c echo.Context) error {
	var (
		payload users.Users
	)

	// convert string to ObjectID
	groupID, errFixedGroupID := primitive.ObjectIDFromHex("6094d237ea050fcfbe6f1dbf") // Unknown Group
	if errFixedGroupID != nil {
		fmt.Println(errFixedGroupID)
		return c.JSON(500, echo.Map{
			"status": false,
			"result": errFixedGroupID.Error(),
		})
	}

	payload.UserGroupID = groupID
	payload.LastUpdated = time.Now()

	err := c.Bind(&payload)
	fmt.Println(payload)
	if err != nil {
		fmt.Println(err)
		return c.JSON(422, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}
	if err := payload.UpdateUserGroup(); err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}
	return c.JSON(200, echo.Map{
		"status":  true,
		"message": "Update user's group successful.",
	})
}
