package worktype

import (
	"fmt"
	"go-callcenter/models/worktype"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getWorkType(c echo.Context) error {
	fmt.Println(("/work-type/get"))
	var (
		payload worktype.WorkType
	)
	result, err := payload.GetWorkType()
	if err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status":  false,
			"error":   err.Error(),
			"message": "Can't get work type.",
		})
	}

	// if there is no work sub type, it will return
	// work_sub_type : [{
	// call_center_groups: [],
	// other_groups: []
	// }]
	// don't know correct query
	// have to check on website again
	return c.JSON(200, echo.Map{
		"status": true,
		"result": result,
	})
}

func createWorkType(c echo.Context) error {
	fmt.Println("/work-type/create")
	var (
		payload worktype.WorkType
		t       = time.Now()
	)
	if err := c.Bind(&payload); err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status":  false,
			"messege": "Failed at bind function.",
			"error":   err.Error(),
		})
	}

	payload.Active = true
	payload.LastUpdated = t

	if err := payload.CreateWorkType(); err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status":  false,
			"messege": "Failed at create work type.",
			"error":   err.Error(),
		})
	}
	return c.JSON(201, echo.Map{
		"status":  true,
		"message": "Add new work type successful.",
		// "result":  payload,
	})
}

func editWorkType(c echo.Context) error {
	fmt.Println("/work-type/edit")
	var (
		payload worktype.WorkType
		t       = time.Now()
	)

	userObjID, _ := primitive.ObjectIDFromHex(c.FormValue("_id"))
	if err := c.Bind(&payload); err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Failed in bind function.",
			"error":   err.Error(),
		})
	}

	payload.ID = userObjID
	payload.Active = true
	payload.LastUpdated = t

	if err := payload.EditWorkType(); err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Failed in bind function.",
			"error":   err.Error(),
		})
	}
	return c.JSON(200, echo.Map{
		"status":  true,
		"message": "Edit work type successful.",
		// "result": payload,
	})
}

func getWorkTypeAll(c echo.Context) error {
	fmt.Println("/work-type/get-all")
	var (
		payload worktype.WorkType
	)
	result, err := payload.GetWorkTypeAll()
	if err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Failed in get work type all.",
			"error":   err.Error(),
		})
	}
	return c.JSON(200, echo.Map{
		"status": true,
		"result": result,
	})
}

//DelectWorkType function
func deleteWorkType(c echo.Context) error {
	var (
		payload worktype.WorkType
	)
	userID := c.Param("user-id")
	objectUserID, _ := primitive.ObjectIDFromHex(userID)
	payload.ID = objectUserID
	if err := payload.DeleteWorkType(); err != nil {
		return c.JSON(500, echo.Map{
			"status":  false,
			"message": "Failed in delete function",
			"error":   err.Error(),
		})
	}
	return c.JSON(200, echo.Map{
		"status": true,
		"result": payload,
	})
}

func getAllWorkTypeAndWorkSubType(c echo.Context) error {
	var (
		payload worktype.WorkType
	)
	slaLevelInput, _ := strconv.Atoi(c.Param("sla-level"))
	var slaLevel int
	if slaLevelInput == 0 {
		result, err := payload.GetAllWorkTypeAndWorkSubTypeAll()
		if err != nil {
			fmt.Println("error in GetAllWorkTypeAndWorkSubTypeAll", err)
			return c.JSON(500, echo.Map{
				"status":  false,
				"error":   err,
				"message": "error in GetAllWorkTypeAndWorkSubTypeAll",
			})
		}
		return c.JSON(200, echo.Map{
			"status": true,
			"result": result,
		})
	}
	if slaLevelInput != 0 {
		slaLevel = slaLevelInput
	}
	fmt.Println(slaLevel)
	seach := c.Param("seach")
	fmt.Println(seach)
	resutl, err := payload.GetAllWorkTypeAndWorkSubType(slaLevel, seach)
	if err != nil {
		return c.JSON(500, echo.Map{
			"status": false,
			"error":  err.Error(),
		})
	}
	return c.JSON(200, echo.Map{
		"status": true,
		"result": resutl,
	})
}
