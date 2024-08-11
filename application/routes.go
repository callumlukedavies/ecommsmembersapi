package application

import "github.com/gin-gonic/gin"

func (a *App) loadRoutes() {
	router := gin.Default()

	shopRoute := router.Group("/shopapi")
	{
		shopRoute.GET("/")
		shopRoute.PUT("/:ID")
		shopRoute.POST("/:ID")
		shopRoute.DELETE("/:ID")
	}
}
