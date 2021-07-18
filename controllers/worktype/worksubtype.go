package worktype

import (
	"fmt"
	"go-callcenter/models/worktype"
	"strings"
	"time"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createWorkSubType(c echo.Context) error {
	fmt.Println(("/work-sub-type/create"))
	var (
		payload worktype.WorkSubType
		t       = time.Now()
	)
	if err := c.Bind(&payload); err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Failed in bind function.",
			"error":   err.Error(),
		})
	}

	othersGroupString := c.FormValue("other_groups")
	if othersGroupString == "" {
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Other groups is required.",
		})
	}

	callCenterGroupsString := c.FormValue("call_center_groups")
	if callCenterGroupsString == "" {
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Call center groups is required.",
		})
	}

	payload.Active = true
	payload.LastUpdated = t
	payload.WorkTypeID, _ = primitive.ObjectIDFromHex(c.FormValue("work_type_id"))
	otherGroups := strings.Split(othersGroupString, ",")
	callCenterGroups := strings.Split(callCenterGroupsString, ",")

	var resultCCGroups []primitive.ObjectID
	for _, v := range callCenterGroups {
		s := string(v)
		dummy, _ := primitive.ObjectIDFromHex(s)
		resultCCGroups = append(resultCCGroups, dummy)
	}
	payload.CallCenterGroups = resultCCGroups

	var resultOthersGroups []primitive.ObjectID
	for _, v := range otherGroups {
		s := string(v)
		dummy, _ := primitive.ObjectIDFromHex(s)
		resultOthersGroups = append(resultOthersGroups, dummy)
	}
	payload.OtherGroups = resultOthersGroups

	if err := payload.CreateWorkSubType(); err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Failed at create work sub type.",
			"error":   err.Error(),
		})
	}
	return c.JSON(200, echo.Map{
		"status":  true,
		"message": "Create work sub type successful.",
		// "result": payload,
	})
}

func getWorkSubTypeAll(c echo.Context) error {
	var (
		payload = worktype.WorkSubType{}
	)
	result, err := payload.GetWorkSubTypeAll()
	if err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status": false,
			"error":  err.Error(),
			// "message": "Failed in GetWorkSubTypeAll",
		})
	}
	return c.JSON(200, echo.Map{
		"status": true,
		"result": result,
	})
}

func editWorkSubType(c echo.Context) error {
	fmt.Println("/work-sub-type/edit")
	var (
		payload worktype.EditWorkSubType
		t       = time.Now()
	)

	othersGroupString := c.FormValue("other_groups")
	if othersGroupString == "" {
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Other groups is required.",
		})
	}

	callCenterGroupsString := c.FormValue("call_center_groups")
	if callCenterGroupsString == "" {
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Call center groups is required.",
		})
	}

	payload.Active = true
	payload.LastUpdated = t
	payload.ID, _ = primitive.ObjectIDFromHex(c.FormValue("_id"))
	payload.WorkTypeID, _ = primitive.ObjectIDFromHex(c.FormValue("work_type_id"))

	otherGroups := strings.Split(othersGroupString, ",")
	callCenterGroups := strings.Split(callCenterGroupsString, ",")

	var resultCCGroups []primitive.ObjectID
	for _, v := range callCenterGroups {
		s := string(v)
		dummy, _ := primitive.ObjectIDFromHex(s)
		resultCCGroups = append(resultCCGroups, dummy)
	}
	payload.CallCenterGroups = resultCCGroups

	var resultOthersGroups []primitive.ObjectID
	for _, v := range otherGroups {
		s := string(v)
		dummy, _ := primitive.ObjectIDFromHex(s)
		resultOthersGroups = append(resultOthersGroups, dummy)
	}
	payload.OtherGroups = resultOthersGroups

	if err := c.Bind(&payload); err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Failed in bind function.",
			"error":   err.Error(),
		})
	}
	if err := payload.PutWorkSubType(); err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Failed at edit work sub type.",
			"error":   err.Error(),
		})
	}
	return c.JSON(200, echo.Map{
		"status":  true,
		"message": "Edit work sub type successful.",
		// "result": payload,
	})
}

func deleteWorkSubType(c echo.Context) error {
	fmt.Println("/work-sub-type/delete")
	var (
		payload worktype.EditWorkSubType
		t       = time.Now()
	)
	workSubTypeId, err := primitive.ObjectIDFromHex(c.Param("work-sub-type-id"))
	if err != nil {
		return c.JSON(500, echo.Map{
			"status": false,
			"error ": err.Error(),
		})
	}
	payload.ID = workSubTypeId
	payload.Active = false
	payload.LastUpdated = t
	if err := payload.DeleteWorkSubType(); err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Failed at delete work sub type.",
			"error":   err.Error(),
		})
	}
	return c.JSON(200, echo.Map{
		"status":  true,
		"message": "Delete work sub type successful.",
	})
}
