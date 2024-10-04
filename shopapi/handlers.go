package shopapi

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type Shop struct {
	DataAccess DataAccess
}

func (shop *Shop) GetAllProductsHandler(c *gin.Context, store *sessions.CookieStore) {
	session, err := store.Get(c.Request, "session")
	if err != nil {
		// User session may not exist, so don't return from error
		log.Println("GetAllProductsHandler: Error getting session: %s", err.Error())
	}

	isAuthenticated := session.Values["Authenticated"]

	items, err := shop.DataAccess.GetAllItems()
	if err != nil {
		log.Fatal(err)
		return
	}

	var response struct {
		Items []Item `json:"shopitems"`
	}

	response.Items = items

	templates, err := template.ParseFiles("templates/layout.html", "templates/navbar.html", "templates/itemsgrid.html", "templates/item.html")
	if err != nil {
		log.Printf("GetAllProductsHandler: Error parsing templates: %v", err)
	}

	c.Header("Content-Type", "text/html")

	err = templates.ExecuteTemplate(c.Writer, "layout.html", gin.H{
		"Title":           "Home",
		"items":           response.Items,
		"isAuthenticated": isAuthenticated,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}

func (shop *Shop) CreateItemHandler(c *gin.Context) {
	itemName := c.PostForm("name")
	itemPrice := c.PostForm("price")
	itemImage, err := c.FormFile("image")
	if err != nil {
		c.String(http.StatusBadRequest, "Could not upload item: Error uploading item photo")
		return
	}

	path := filepath.Join("images", itemImage.Filename)
	if err = c.SaveUploadedFile(itemImage, path); err != nil {
		c.String(http.StatusInternalServerError, "Failed to save image: %s", err.Error())
		return
	}

	//validate file type

	err = shop.DataAccess.CreateItem(itemName, itemPrice, itemImage.Filename)

	if err != nil {
		fmt.Println(err)
		return
	}

	// c.HTML(http.StatusOK, "item.html", gin.H{
	// 	"Name":  itemName,
	// 	"Price": itemPrice,
	// 	"ID":    itemID,
	// })
}

func (shop *Shop) UpdatePriceHandler(c *gin.Context) {
	itemPrice := c.PostForm("price")
	paramID := c.Param("ID")
	paramName := c.Param("Name")
	itemID, _ := strconv.ParseInt(paramID, 10, 64)
	err := shop.DataAccess.UpdatePrice(itemID, itemPrice)

	if err != nil {
		fmt.Println(err)
		return
	}

	c.HTML(http.StatusOK, "item.html", gin.H{
		"Name":  paramName,
		"Price": itemPrice,
		"ID":    itemID,
	})
}

func (shop *Shop) DeleteItemHandler(c *gin.Context) {
	param := c.Param("ID")
	id, _ := strconv.ParseInt(param, 10, 64)
	shop.DataAccess.DeleteItem(id)

	c.HTML(http.StatusOK, "deleteditem.html", nil)
}
