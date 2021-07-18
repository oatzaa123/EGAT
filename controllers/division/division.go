package division

// import (
// 	division "go-callcenter/models/division"

// 	"github.com/labstack/echo"
// )

// // PrivateRoute comment
// func PrivateRoute(g *echo.Group) {
// 	g.GET("/division/get", getDivision)
// }

// func getDivision(c echo.Context) error {
// 	var (
// 		payload division.Division
// 	)

// 	result, err := payload.GetDivision()
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
