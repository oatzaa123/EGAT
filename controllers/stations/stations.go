package stations

import (
	"go-callcenter/models/stations"

	"github.com/labstack/echo"
)

// PrivateRoute comment
func PrivateRoute(g *echo.Group) {
	g.GET("/station/getStations", getStations)
}

func getStations(c echo.Context) error {
	var (
		payload stations.Stations
	)

	result, err := payload.GetStations()
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
