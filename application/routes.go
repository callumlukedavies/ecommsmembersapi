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
	router.Static("/images", "./images/")
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
		shopRoute.GET("/", func(c *gin.Context) { shop.GetAllProductsHandler(c, store) })
		shopRoute.POST("/create", func(c *gin.Context) { shop.CreateItemHandler(c, store) })
		shopRoute.PUT("/:ID/:Name", shop.UpdatePriceHandler)
		shopRoute.DELETE("/:ID", shop.DeleteItemHandler)
		shopRoute.GET("/view/:ID", func(c *gin.Context) { shop.ViewItemHandler(c, store) })
	}

	userDatabaseRoute := router.Group("/membersapi")
	{
		userDatabaseRoute.POST("/login", func(c *gin.Context) { userDatabase.LoginHandler(c, store) })
		userDatabaseRoute.GET("/logout", func(c *gin.Context) { userDatabase.LogoutHandler(c, store) })
		userDatabaseRoute.GET("/signup", userDatabase.GetSignUpPageHandler)
		userDatabaseRoute.GET("/profile", middleware.AuthorizeUser(store), func(c *gin.Context) { userDatabase.GetProfilePageHandler(c, store) })
		userDatabaseRoute.GET("/editpage", middleware.AuthorizeUser(store), func(c *gin.Context) { userDatabase.GetEditPageHandler(c, store) })
		userDatabaseRoute.POST("/createuser", func(c *gin.Context) { userDatabase.CreateUserHandler(c, store) })
		userDatabaseRoute.POST("/edit-user-firstname", func(c *gin.Context) { userDatabase.EditUserFirstNameHandler(c, store) })
		userDatabaseRoute.POST("/edit-user-lastname", func(c *gin.Context) { userDatabase.EditUserLastNameHandler(c, store) })
		userDatabaseRoute.POST("/edit-user-dateofbirth", func(c *gin.Context) { userDatabase.EditUserDateOfBirthHandler(c, store) })
		userDatabaseRoute.POST("/edit-user-emailaddress", func(c *gin.Context) { userDatabase.EditUserEmailHandler(c, store) })
		userDatabaseRoute.POST("/edit-user-password", func(c *gin.Context) { userDatabase.EditUserPasswordHandler(c, store) })

		userDatabaseRoute.DELETE("/:ID", func(c *gin.Context) { userDatabase.DeleteUserHandler(c, store) })
	}

	router.Run(":8080")
}
