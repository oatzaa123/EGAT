package users

import (
	"fmt"
	"go-callcenter/keys"
	"go-callcenter/models/users"
	"net"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

//Login function comment
func login(c echo.Context) error {
	var (
		payload        users.Users
		login          users.Login
		loginHistories users.LoginHistories
		t              = time.Now()
	)

	// recieve username and password
	if err := c.Bind(&login); err != nil {
		fmt.Println(err)
	}

	// find username in database
	if err := payload.FindUsername(login.Username); err != nil {
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Username or password is incorrect.",
		})
	}

	// compare hash and password
	if err := bcrypt.CompareHashAndPassword([]byte(payload.Password), []byte(login.Password)); err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Username or password is incorrect.",
		})
	}

	// check ip address
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				loginHistories.IPAddress = ipnet.IP.String()
			}
		}
	}
	// Trial
	// addrs, err := net.InterfaceAddrs()
	// if err != nil {
	// 	panic(err)
	// }
	// for i, addr := range addrs {
	// 	fmt.Printf("%d %v\n", i, addr)
	// }

	// create token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	claims["_id"] = payload.ID
	claims["division"] = payload.Division
	claims["user_group_id"] = payload.UserGroupID
	claims["image"] = payload.Image
	claims["first_name"] = payload.FirstName
	claims["last_name"] = payload.LastName
	claims["date_of_birth"] = payload.DateOfBirth
	claims["gender"] = payload.Gender
	claims["username"] = payload.UserName
	claims["created_date"] = payload.CreateDate
	claims["active"] = payload.Active
	claims["last_updated"] = payload.LastUpdated
	claims["staff_id"] = payload.StaffID

	tokenString, err := token.SignedString([]byte(keys.Secret))
	if err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Can't create token.",
		})
	}

	// insert login history to database
	loginHistories.UserID = payload.ID
	loginHistories.LastUpdated = t
	if err := loginHistories.CreateLoginHistories(); err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": err,
		})
	}

	return c.JSON(200, echo.Map{
		"status": true,
		"result": tokenString,
		"name":   payload.FirstName,
		"image":  payload.Image,
	})
}

func getLoginHistories(c echo.Context) error {
	var (
		payload users.LoginHistories
	)
	usersID := c.Param("user-id")
	objectUsersID, err := primitive.ObjectIDFromHex(usersID)
	if err != nil {
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Invalid user ID to get login hostories.",
		})
	}
	payload.ID = objectUsersID
	loginUsersHistories, err := payload.GetLoginHistories()
	if err != nil {
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Failed to get login hostories.",
		})
	}
	return c.JSON(200, echo.Map{
		"status": true,
		"result": loginUsersHistories,
	})
}
