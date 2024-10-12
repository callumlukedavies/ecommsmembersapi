package shopapi

import (
	"ecommercesite/util"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

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

	for i := 0; i < len(response.Items); i++ {
		response.Items[i].ImageName = util.GetFirstImageFromString(items[i].ImageName)
	}

	templates, err := template.ParseFiles("templates/layout.html", "templates/navbar.html", "templates/itemsgrid.html", "templates/item.html")
	if err != nil {
		log.Printf("GetAllProductsHandler: Error parsing templates: %v", err)
		return
	}

	c.Header("Content-Type", "text/html")

	err = templates.ExecuteTemplate(c.Writer, "layout.html", gin.H{
		"Title":           "Home",
		"items":           response.Items,
		"isAuthenticated": isAuthenticated,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
		return
	}
}

func (shop *Shop) CreateItemHandler(c *gin.Context, store *sessions.CookieStore) {
	session, err := store.Get(c.Request, "session")
	if err != nil {
		// User session may not exist, so don't return from error
		log.Println("CreateItemHandler: Error getting session: %s", err.Error())
	}

	isAuthenticated := session.Values["Authenticated"].(bool)
	if !isAuthenticated {
		log.Println("CreateItemHandler: User must be authenticated to create an item")
		c.String(http.StatusBadRequest, "User must be authenticated before uploading an item. Please log in.")
		return
	}

	form, formErr := c.MultipartForm()
	if formErr != nil {
		log.Println("CreateItemHandler: Error getting multipart form. Error: %s", formErr.Error())
		c.String(http.StatusBadRequest, "There has been an issue getting the multipart form. Please try again.")
		return
	}

	itemName := c.PostForm("name-input")
	itemGender := c.PostForm("gender-input")
	itemDesc := c.PostForm("description-input")
	images := form.File["image-input"]
	itemPrice := c.PostForm("price-input")
	itemSize := c.PostForm("size-input")
	itemCategory := c.PostForm("category-input")
	itemSellerName := session.Values["FirstName"].(string) + " " + session.Values["LastName"].(string)
	itemSellerID := session.Values["UserID"].(int)
	itemCondition := c.PostForm("condition-input")
	var imageList string

	for _, image := range images {
		fileName := filepath.Base(image.Filename)
		filePath := filepath.Join("images", fileName)

		//validate file type
		isValid := util.ValidateImage(filePath)
		if !isValid {
			c.String(http.StatusInternalServerError, "Could not upload item: image file extension not valid. Image name: %s", filePath)
			return
		}

		if err = c.SaveUploadedFile(image, filePath); err != nil {
			c.String(http.StatusInternalServerError, "Failed to save image. Error: %s", err.Error())
			return
		}

		imageList += fileName
		imageList += ";"
	}

	uploadDate := time.Now()
	itemKey, err := util.GenerateRandomKey(10)
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not generate random item Key. Error: %s", err.Error())
		return
	}

	item := Item{
		ID:           itemKey,
		Name:         itemName,
		Gender:       itemGender,
		Description:  itemDesc,
		ImageName:    imageList,
		DateUploaded: uploadDate.Format(time.DateOnly),
		Price:        itemPrice,
		IsSold:       false,
		Size:         itemSize,
		Category:     itemCategory,
		Condition:    itemCondition,
		SellerName:   itemSellerName,
		SellerID:     itemSellerID,
	}

	err = shop.DataAccess.CreateItem(item)

	if err != nil {
		fmt.Println(err)
		return
	}
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

func (shop *Shop) ViewItemHandler(c *gin.Context, store *sessions.CookieStore) {

	session, err := store.Get(c.Request, "session")
	if err != nil {
		// User session may not exist, so don't return from error
		log.Println("ViewItemHandler: Error getting session: %s", err.Error())
	}

	isAuthenticated := session.Values["Authenticated"]

	itemID := c.Param("ID")

	// Retrieve item from database
	itemData, err := shop.DataAccess.GetItem(itemID)
	templates, err := template.ParseFiles("templates/layout.html", "templates/navbar.html", "templates/itemview.html")
	if err != nil {
		log.Printf("GetProfilePageHandler: Error parsing templates: %s", err.Error())
	}

	images := util.ParseImageString(itemData.ImageName)

	c.Header("Content-Type", "text/html")

	// Execute the main layout template with the "signup" content embedded
	err = templates.ExecuteTemplate(c.Writer, "layout.html", gin.H{
		"isAuthenticated": isAuthenticated,
		"ID":              itemData.ID,
		"itemName":        itemData.Name,
		"itemDescription": itemData.Description,
		"images":          images,
		"itemUploadDate":  itemData.DateUploaded,
		"itemPrice":       itemData.Price,
		"itemIsSold":      itemData.IsSold,
		"itemCategory":    itemData.Category,
		"itemCondition":   itemData.Condition,
		"itemSize":        itemData.Size,
		"itemSellerName":  itemData.SellerName,
		"itemSellerID":    itemData.SellerID,
	})

	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}

func (shop *Shop) SearchHandler(c *gin.Context, store *sessions.CookieStore) {
	session, err := store.Get(c.Request, "session")
	if err != nil {
		// User session may not exist, so don't return from error
		log.Println("SearchHandler: Error getting session: %s", err.Error())
	}

	isAuthenticated := session.Values["Authenticated"]

	queryTerm := c.PostForm("search-input")
	query := "%" + queryTerm + "%"

	items, err := shop.DataAccess.GetItemsByQueryTerm(query)
	if err != nil {
		log.Fatal(err)
		return
	}

	var response struct {
		Items []Item `json:"shopitems"`
	}

	response.Items = items

	for i := 0; i < len(response.Items); i++ {
		response.Items[i].ImageName = util.GetFirstImageFromString(items[i].ImageName)
	}

	templates, err := template.ParseFiles("templates/layout.html", "templates/navbar.html", "templates/itemsgrid.html", "templates/item.html")
	if err != nil {
		log.Printf("SearchHandler: Error parsing templates: %v", err)
		return
	}

	c.Header("Content-Type", "text/html")

	err = templates.ExecuteTemplate(c.Writer, "layout.html", gin.H{
		"Title":           "Home",
		"items":           response.Items,
		"isAuthenticated": isAuthenticated,
	})

	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
		return
	}
}
