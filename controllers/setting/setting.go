package setting

import (
	"fmt"
	"go-callcenter/models/setting"

	"github.com/labstack/echo"
)

// PrivateRoute comment
func PrivateRoute(g *echo.Group) {
	g.GET("/general-setting/get", getGeneralSetting)
	g.POST("/general-setting/update", updateGeneralSetting)
}

func getGeneralSetting(c echo.Context) error {
	var (
		payload setting.Setting
	)

	// if err := c.Bind(&payload); err != nil {
	// 	return c.JSON(400, echo.Map{
	// 		"status": false,
	// 	})
	// }

	if err := payload.GetSetting(); err != nil {
		fmt.Println(payload)
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"status": true,
		"result": payload,
	})
}

func updateGeneralSetting(c echo.Context) error {
	var (
		payload setting.Setting
	)

	// payload.OpeningWord = c.QueryParam("something")
	// payload.OpeningWord = c.FormValue("someelse")
	if err := c.Bind(&payload); err != nil {
		fmt.Println(err)
		fmt.Println(payload)
		return c.JSON(422, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	if err := payload.UpdateSetting(); err != nil {
		return c.JSON(500, echo.Map{
			"status": false,
			"result": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"status": true,
		"result": payload,
	})
}
