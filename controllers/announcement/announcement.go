package announcement

import (
	"fmt"
	announcement "go-callcenter/models/announcement"
	"time"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PrivateRoute comment
func PrivateRoute(g *echo.Group) {
	g.POST("/announcement/add", addAnnouncement)
	g.GET("/announcement/get", getAnnouncement)
	g.PUT("/announcement/edit-detail", updataAnnouncement)
	g.DELETE("/announcement/delete/:anm_id", deleteAnnouncement)
	g.PUT("/announcement/edit-order", rearrangeOrder)
}

func getAnnouncement(c echo.Context) error {
	fmt.Println("/announcement/get")
	var (
		payload announcement.Announcement
	)
	importanceAnnouncement, otherAnnouncement, err := payload.GetAnnouncement()
	if err != nil {
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}
	result := echo.Map{
		"status":                  true,
		"other_announcement":      otherAnnouncement,
		"importance_announcement": importanceAnnouncement,
	}
	return c.JSON(200, result)
}

func addAnnouncement(c echo.Context) error {
	fmt.Println("/announcement/add")
	var (
		payload announcement.Announcement
	)
	if err := c.Bind(&payload); err != nil {
		fmt.Println(err)
		return c.JSON(422, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	if err := payload.CountOtherAnnouncement(); err != nil {
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	if err := payload.ArrangeOrderType(); err != nil {
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	payload.LastUpdated = time.Now()
	payload.Order = 1
	payload.Active = true
	payload.Type = "OTHERS"

	if err := payload.AddAnnouncement(); err != nil {
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	return c.JSON(201, echo.Map{
		"status":  true,
		"result":  payload,
		"message": "Add new announcement successful.",
	})
}

func updataAnnouncement(c echo.Context) error {
	fmt.Println("/announcement/edit-detail")
	var (
		payload announcement.Announcement
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
	if err := payload.UpdateAnnouncement(); err != nil {
		fmt.Println(err)
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}
	return c.JSON(200, echo.Map{
		"status":  true,
		"message": "Edit announcement successful.",
	})
}

func deleteAnnouncement(c echo.Context) error {
	fmt.Println("/announcement/delete")
	var (
		payload announcement.Announcement
	)

	anmID := c.Param("anm_id")
	fmt.Println("Announcement ID :", anmID)

	objectAnmID, errAmnID := primitive.ObjectIDFromHex(anmID)
	if errAmnID != nil {
		fmt.Println(errAmnID)
	}

	payload.ID = objectAnmID

	anmOrder, err := payload.FindAnnouncementOrderByID()
	if err != nil {
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	if err := payload.DeleteAnnouncement(); err != nil {
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	fmt.Println("Announcement order :", anmOrder)
	payload.Order = anmOrder

	if err := payload.ReplaceOrderAfterDeleteAnm(); err != nil {
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"status": true,
		// "result": payload,
		"message": "Delete announcement successful.",
	})
}

func rearrangeOrder(c echo.Context) error {
	fmt.Println("/announcement/edit-order")
	var (
		payload announcement.Announcement
	)

	// 1. find announcement by id that want to change order
	// 2. find announcement by order that must to swap order with 1.

	if err := c.Bind(&payload); err != nil {
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	updateAnm := payload.ID
	swapOrder := payload.Order

	currentOrder, err := payload.FindAnnouncementOrderByID()
	if err != nil {
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	swapAnmID, err := payload.FindAnnouncementIDByOrder(swapOrder)
	if err != nil {
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	fmt.Println("Update order ID :", updateAnm)
	fmt.Println("Current order   :", currentOrder)
	fmt.Println("Swap order ID   :", swapAnmID)
	fmt.Println("Swap order      :", swapOrder)

	errUpdateAnm := payload.RearrangeOrder(updateAnm, swapOrder)
	if errUpdateAnm != nil {
		return c.JSON(500, echo.Map{
			"status": false,
			"result": errUpdateAnm.Error(),
		})
	}

	errSwapAnm := payload.RearrangeOrder(swapAnmID, currentOrder)
	if errSwapAnm != nil {
		return c.JSON(500, echo.Map{
			"status": false,
			"result": errSwapAnm.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"status":  true,
		"message": "Update order announcement successful.",
	})
}
