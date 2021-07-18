package worktype

import "github.com/labstack/echo"

func PrivateRoute(g *echo.Group) {
	g.GET("/work-type/get", getWorkType)
	g.GET("/work-type/get-all", getWorkTypeAll)
	g.POST("/work-type/create", createWorkType)
	g.PUT("/work-type/edit", editWorkType)

	g.POST("/work-sub-type/create", createWorkSubType)
	g.PUT("/work-sub-type/edit", editWorkSubType)
	g.DELETE("/work-sub-type/delete/:work-sub-type-id", deleteWorkSubType)

	g.DELETE("/work-type/delete/:user-id", deleteWorkType)
	g.GET("/work-sub-type/getall", getWorkSubTypeAll)
	g.GET("/work-type/get/work-sub-type/:sla-level/:seach", getAllWorkTypeAndWorkSubType)
}
