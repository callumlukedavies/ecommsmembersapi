package application

import (
	"ecommercesite/shopapi"

	"github.com/gin-gonic/gin"
)

func (a *App) loadRoutes() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static/")

	shop := shopapi.Shop{
		DataAccess: shopapi.DataAccess{
			DB: a.db,
		},
	}

	shopRoute := router.Group("/shopapi")
	{
		shopRoute.GET("/", shop.GetAllProductsHandler)
		shopRoute.PUT("/:ID/:Name", shop.UpdatePriceHandler)
		shopRoute.POST("/", shop.CreateItemHandler)
		shopRoute.DELETE("/:ID", shop.DeleteItemHandler)
	}

	router.Run(":8080")
}
