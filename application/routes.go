package application

import (
	"ecommercesite/membersapi"
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

	userDatabase := membersapi.UserDatabase{
		DataAccess: membersapi.DataAccess{
			DB: a.db,
		},
	}

	shopRoute := router.Group("/shopapi")
	{
		shopRoute.GET("/", shop.GetAllProductsHandler)
		shopRoute.POST("/", shop.CreateItemHandler)
		shopRoute.PUT("/:ID/:Name", shop.UpdatePriceHandler)
		shopRoute.DELETE("/:ID", shop.DeleteItemHandler)
	}

	userDatabaseRoute := router.Group("/membersapi")
	{
		userDatabaseRoute.POST("/signin", userDatabase.GetUserHandler)
		userDatabaseRoute.GET("/signup", userDatabase.GetSignUpPageHandler)
		userDatabaseRoute.POST("/createuser", userDatabase.CreateUserHandler)
		userDatabaseRoute.PUT("/:ID", userDatabase.GetUserHandler)
		userDatabaseRoute.DELETE("/:ID", userDatabase.DeleteUserHandler)
	}

	router.Run(":8080")
}
