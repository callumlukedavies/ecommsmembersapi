package application

import (
	"ecommercesite/membersapi"
	"ecommercesite/middleware"
	"ecommercesite/shopapi"
	"ecommercesite/util"

	"github.com/gin-gonic/gin"
)

func (a *App) loadRoutes() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static/")
	store := util.InitializeStore()

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
		userDatabaseRoute.POST("/login", func(c *gin.Context) { userDatabase.LoginHandler(c, store) })
		userDatabaseRoute.GET("/signup", userDatabase.GetSignUpPageHandler)
		userDatabaseRoute.GET("/profile", middleware.AuthorizeUser(store), func(c *gin.Context) { userDatabase.GetProfilePageHandler(c, store) })
		userDatabaseRoute.POST("/createuser", func(c *gin.Context) { userDatabase.CreateUserHandler(c, store) })
		userDatabaseRoute.DELETE("/:ID", userDatabase.DeleteUserHandler)
	}

	router.Run(":8080")
}
