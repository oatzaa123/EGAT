package group

import (
	"fmt"
	group "go-callcenter/models/group"
	"time"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PrivateRoute comment
func PrivateRoute(g *echo.Group) {
	g.GET("/group/get", getCallCenterGroups)
	// g.GET("/group/others-group", getOthersGroups)
	g.GET("/group/all-groups-and-members", getMembersInCallCenterGroups)
	g.GET("/group/group-members/:group_id", getMembersForChangeCallCenterGroups)
	g.PUT("/group/change-members-to-group/:group_id", membersChangeGroups)
	g.PUT("/group/edit-detail", updateGroupDetail)
	g.GET("/group/:group_id", getGroupByID)
}

func getCallCenterGroups(c echo.Context) error {
	var (
		payload group.UserGroups
	)

	result, err := payload.GetCallCenterGroups()
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

// func getOthersGroups(c echo.Context) error {
// 	var (
// 		payload group.OthersGroup
// 	)

// 	result, err := payload.GetOthersGroups()
// 	if err != nil {
// 		return c.JSON(500, echo.Map{
// 			"status": false,
// 			"result": err.Error(),
// 		})
// 	}

// 	return c.JSON(200, echo.Map{
// 		"status": true,
// 		"result": result,
// 	})
// }

func getMembersInCallCenterGroups(c echo.Context) error {
	var (
		payload group.UserGroups
	)

	result, err := payload.GetMembersInUserGroups()
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

// รับเป็น _id ของ callcentergroups
func getMembersForChangeCallCenterGroups(c echo.Context) error {
	var (
		payload group.UserGroups
	)

	// Get _id from param
	userGroupsID := c.Param("group_id")
	objectCallCenterGroupsID, err := primitive.ObjectIDFromHex(userGroupsID)
	if err != nil {
		fmt.Println(err)
	}
	payload.ID = objectCallCenterGroupsID

	callCenterUsersForChangeCallCenterGroups, err := payload.GetMembersForChangeUserGroups()
	if err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Can't get members.",
		})
	}

	return c.JSON(200, echo.Map{
		"status": true,
		"result": callCenterUsersForChangeCallCenterGroups,
	})
}

func membersChangeGroups(c echo.Context) error {
	var (
		payload group.UserGroups
	)

	callCenterGroupsID := c.Param("group_id")
	objectCallCenterGroupsID, err := primitive.ObjectIDFromHex(callCenterGroupsID)
	if err != nil {
		fmt.Println(err)
	}
	payload.ID = objectCallCenterGroupsID

	var usersIDForChangeCallCenterGroups []string
	if err := c.Bind(&usersIDForChangeCallCenterGroups); err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": err,
		})
	}

	if err := payload.MembersChangeGroups(usersIDForChangeCallCenterGroups); err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Can't update user's group.",
		})
	}

	return c.JSON(200, echo.Map{
		"status":  true,
		"message": "Update user's group complete.",
	})
}

func updateGroupDetail(c echo.Context) error {
	var (
		payload group.UserGroups
	)
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
	if err := payload.UpdateUserGroups(); err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}
	return c.JSON(200, echo.Map{
		"status":  true,
		"message": "Update group detail successful.",
	})
}

func getGroupByID(c echo.Context) error {
	var (
		payload group.UserGroups
	)

	// get _id from param
	callCenterGroupsID := c.Param("group_id")
	objectGroupsID, err := primitive.ObjectIDFromHex(callCenterGroupsID)
	if err != nil {
		fmt.Println(err)
	}
	payload.ID = objectGroupsID

	GroupDetail, err := payload.GetGroupDetail()
	if err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Can't get group detail.",
		})
	}

	return c.JSON(200, echo.Map{
		"status": true,
		"result": GroupDetail[0],
	})
}
