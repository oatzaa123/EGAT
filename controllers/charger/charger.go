package charger

import (
	"go-callcenter/models/charger"

	"github.com/labstack/echo"
)

func PrivateRoute(g *echo.Group) {
	g.GET("/charger/getChargers", getChargers)
}

func getChargers(c echo.Context) error {
	var (
		payload charger.Chargers
	)

	result, err := payload.GetChargers()
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
