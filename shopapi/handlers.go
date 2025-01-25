package shopapi

import (
	"ecommercesite/util"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type Shop struct {
	DataAccess DataAccess
}

var GenderOptions = []string{"Womens", "Mens", "Unisex"}
var CategoryOptions = []string{"Coat", "Jacket", "Knitwear", "Sweatshirt", "Top", "Bottoms", "Shorts", "Shoe", "Accessories"}
var ConditionOptions = []string{"Brand New", "Like New", "Good Condition", "Fair Condition", "Poor Condition"}

var categoryQueryMap = map[string]bool{
	"%Coat%":        true,
	"%Jacket%":      true,
	"%Knitwear%":    true,
	"%Sweatshirt%":  true,
	"%Top%":         true,
	"%Bottoms%":     true,
	"%Shorts%":      true,
	"%Shoe%":        true,
	"%Accessories%": true,
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
		return
	}

	session.Values["CurrentQuery"] = ""
	session.Save(c.Request, c.Writer)

	c.Header("Content-Type", "text/html")

	err = templates.ExecuteTemplate(c.Writer, "layout.html", gin.H{
		"Title":                "Home",
		"items":                response.Items,
		"isAuthenticated":      isAuthenticated,
		"ShowCategoriesBanner": true,
		"GridTitle":            "New In",
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
	itemGalleryImage := ""
	var imageList string

	for _, image := range images {
		fileBase := filepath.Base(image.Filename)
		//validate file type
		isValid := util.ValidateImage(fileBase)
		if !isValid {
			c.String(http.StatusInternalServerError, "Could not upload item: image file extension not valid. Image name: %s", fileBase)
			return
		}

		newFileName, keyErr := util.GenerateRandomKey(10)
		if keyErr != nil {
			newFileName = fileBase
			keyErr = nil
		} else {
			newFileName += ".jpeg"
		}

		filePath := filepath.Join("images", newFileName)

		if err = c.SaveUploadedFile(image, filePath); err != nil {
			c.String(http.StatusInternalServerError, "Failed to save image. Error: %s", err.Error())
			return
		}

		if itemGalleryImage == "" {
			itemGalleryImage = newFileName
		}

		imageList += newFileName
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
		GalleryImage: itemGalleryImage,
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error Creating item"})
		return
	}

	c.Redirect(http.StatusFound, "/shopapi/")
}

func (shop *Shop) EditItemHandler(c *gin.Context, store *sessions.CookieStore) {
	itemID := c.Param("ID")
	session, err := store.Get(c.Request, "session")
	if err != nil {
		// User session may not exist, so don't return from error
		log.Println("EditItemHandler: Error getting session: %s", err.Error())
	}

	isAuthenticated := session.Values["Authenticated"].(bool)
	if !isAuthenticated {
		log.Println("EditItemHandler: User must be authenticated to edit an item")
		c.String(http.StatusBadRequest, "User must be authenticated before editing an item. Please log in.")
		return
	}

	itemName := c.PostForm("name-input")
	itemGender := c.PostForm("gender-input")
	itemDesc := c.PostForm("description-input")
	itemPrice := c.PostForm("price-input")
	itemSize := c.PostForm("size-input")
	itemCategory := c.PostForm("category-input")
	itemCondition := c.PostForm("condition-input")

	itemDetails := Item{
		ID:          itemID,
		Name:        itemName,
		Gender:      itemGender,
		Description: itemDesc,
		Price:       itemPrice,
		IsSold:      false,
		Size:        itemSize,
		Category:    itemCategory,
		Condition:   itemCondition,
	}

	err = shop.DataAccess.UpdateItem(itemDetails)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error editing item"})
		return
	}

	c.Redirect(http.StatusFound, "/shopapi/")
}

func (shop *Shop) DeleteItemHandler(c *gin.Context) {
	itemID := c.Param("ID")
	err := shop.DataAccess.DeleteItem(itemID)
	if err != nil {
		log.Printf("DeleteItemHandler: Error deleting item: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
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
	if err != nil {
		log.Printf("ViewItemHandler: Error getting Item data: %s", err.Error())
		return
	}

	// Retrieve items by this Seller
	itemsBySeller, err := shop.DataAccess.GetItemsBySeller(itemData.SellerID, 8)
	if err != nil {
		log.Printf("ViewItemHandler: Error getting items by seller: %s", err.Error())
		return
	}

	templates, err := template.ParseFiles("templates/layout.html", "templates/navbar.html", "templates/itemview.html")
	if err != nil {
		log.Printf("ViewItemHandler: Error parsing templates: %s", err.Error())
		return
	}

	images := util.ParseImageString(itemData.ImageName)

	c.Header("Content-Type", "text/html")

	err = templates.ExecuteTemplate(c.Writer, "layout.html", gin.H{
		"isAuthenticated":      isAuthenticated,
		"ID":                   itemData.ID,
		"itemName":             itemData.Name,
		"itemDescription":      itemData.Description,
		"images":               images,
		"itemUploadDate":       itemData.DateUploaded,
		"itemPrice":            itemData.Price,
		"itemIsSold":           itemData.IsSold,
		"itemCategory":         itemData.Category,
		"itemCondition":        itemData.Condition,
		"itemSize":             itemData.Size,
		"itemSellerName":       itemData.SellerName,
		"itemSellerID":         itemData.SellerID,
		"otherItems":           itemsBySeller,
		"ShowCategoriesBanner": false,
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

	session.Values["CurrentQuery"] = query
	session.Save(c.Request, c.Writer)

	c.Header("Content-Type", "text/html")

	gridTitle := "Search Results for '" + queryTerm + "'"

	err = templates.ExecuteTemplate(c.Writer, "layout.html", gin.H{
		"Title":                "Home",
		"items":                response.Items,
		"isAuthenticated":      isAuthenticated,
		"ShowCategoriesBanner": true,
		"GridTitle":            gridTitle,
	})

	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
		return
	}
}

func (shop *Shop) SortItemsHandler(c *gin.Context, store *sessions.CookieStore) {
	session, err := store.Get(c.Request, "session")
	if err != nil {
		// User session may not exist, so don't return from error
		log.Println("SortItemsHandler: Error getting session: %s", err.Error())
	}

	isAuthenticated := session.Values["Authenticated"]

	var query string
	var ok bool
	currentQuery := session.Values["CurrentQuery"]
	if currentQuery != nil {
		query, ok = currentQuery.(string)
		if !ok {
			c.JSON(500, gin.H{"error": "Current query is invalid session type"})
			return
		}
	}

	isCategory := categoryQueryMap[query]
	category := ""
	gridTitle := ""
	if isCategory {
		category = util.RemoveQueryBrackets(query)
		gridTitle = category
	} else {
		gridTitle = query
	}

	order := c.PostForm("order-input")

	var items []Item
	var sortErr error
	switch order {
	case "price-inc":
		items, sortErr = shop.DataAccess.getSortedItemsByPriceInc(query, category)
	case "price-dec":
		items, sortErr = shop.DataAccess.getSortedItemsByPriceDec(query, category)
	}

	if sortErr != nil {
		log.Printf("SortItemsHandler: Could not sort items by %s. Error: %s", order, err)
		c.String(http.StatusInternalServerError, "Error sorting items: %v", err)
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
		log.Printf("SortItemsHandler: Error parsing templates: %v", err)
		return
	}

	c.Header("Content-Type", "text/html")

	err = templates.ExecuteTemplate(c.Writer, "layout.html", gin.H{
		"Title":                "Home",
		"items":                response.Items,
		"isAuthenticated":      isAuthenticated,
		"ShowCategoriesBanner": true,
		"GridTitle":            gridTitle,
	})

	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
		return
	}
}

func (shop *Shop) SearchByCategoryHandler(c *gin.Context, store *sessions.CookieStore) {
	session, err := store.Get(c.Request, "session")
	if err != nil {
		// User session may not exist, so don't return from error
		log.Println("SearchByCategoryHandler: Error getting session: %s", err.Error())
	}

	isAuthenticated := session.Values["Authenticated"]

	categoryTerm := c.Query("CategoryID")
	category := "%" + categoryTerm + "%"

	items, err := shop.DataAccess.getItemsByCategory(category)
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
		log.Printf("SearchByCategoryHandler: Error parsing templates: %v", err)
		return
	}

	session.Values["CurrentQuery"] = category
	session.Save(c.Request, c.Writer)

	c.Header("Content-Type", "text/html")

	err = templates.ExecuteTemplate(c.Writer, "layout.html", gin.H{
		"Title":                "Home",
		"items":                response.Items,
		"isAuthenticated":      isAuthenticated,
		"ShowCategoriesBanner": true,
		"GridTitle":            categoryTerm,
	})

	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
		return
	}
}

func (shop *Shop) UpdateCartHandler(c *gin.Context, store *sessions.CookieStore) {
	session, err := store.Get(c.Request, "session")
	if err != nil {
		// User session may not exist, so don't return from error
		log.Println("UpdateCartHandler: Error getting session: %s", err.Error())
	}

	var requestBody struct {
		ItemID string `json:"itemID"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	itemData, err := shop.DataAccess.GetItem(requestBody.ItemID)
	if err != nil {
		log.Printf("UpdateCartHandler: Error getting Item data: %s", err.Error())
		return
	}

	if itemData.IsSold {
		c.JSON(http.StatusOK, gin.H{"error": "Item no longer available"})
		return
	}

	var oldCartItems, newCartItems string
	var ok bool
	sessionCartItems := session.Values["CartItems"]
	if sessionCartItems != nil {
		// CartItems exists, get existing list and add new item id to it
		oldCartItems, ok = sessionCartItems.(string)
		if !ok {
			log.Println("UpdateCartHandler: Error converting session CartItems value to string")
		} else {
			newCartItems = oldCartItems + itemData.ID + ";"
		}

	} else {
		// CartItems doesn't yet exist on the user session, set to current Item ID
		newCartItems = itemData.ID + ";"
	}

	session.Values["CartItems"] = newCartItems
	session.Save(c.Request, c.Writer)

	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart!"})
}

func (shop *Shop) GetShoppingCartHandler(c *gin.Context, store *sessions.CookieStore) {

	// Get session values
	session, err := store.Get(c.Request, "session")
	if err != nil {
		log.Println("GetShoppingCartHandler: Error getting session: %s", err.Error())
	}

	firstname := session.Values["FirstName"]
	lastname := session.Values["LastName"]
	email := session.Values["EmailAddress"]
	dateofbirth := session.Values["DateOfBirth"]
	sessionCartItems := session.Values["CartItems"]

	var cartItemsString string
	var ok bool
	if sessionCartItems != nil {
		// Items have been added to cart
		cartItemsString, ok = sessionCartItems.(string)
		if !ok {
			log.Println("GetShoppingCartHandler: Error converting session CartItems to string")
			c.String(http.StatusInternalServerError, "Error converting session CartItems to string")
			return
		}
	} else {
		// No items in cart, return early
		c.String(http.StatusOK, "You have no items in the cart")
		return
	}

	startIndex := 0
	var cartItems []Item
	for i := 0; i < len(cartItemsString); i++ {
		if cartItemsString[i] != ';' {
			continue
		}

		itemData, err := shop.DataAccess.GetItem(cartItemsString[startIndex:i])
		if err != nil {
			log.Printf("GetShoppingCartHandler: Error getting Item data for ID %s: %s", cartItemsString[startIndex:i], err.Error())
			err = nil
			startIndex = i + 1
			continue
		}

		cartItems = append(cartItems, itemData)
		startIndex = i + 1
	}

	templates, err := template.ParseFiles("templates/layout.html", "templates/navbar.html", "templates/cart.html")
	if err != nil {
		log.Printf("GetShoppingCartHandler: Error parsing templates: %s", err.Error())
	}

	c.Header("Content-Type", "text/html")

	err = templates.ExecuteTemplate(c.Writer, "layout.html", gin.H{
		"FirstName":       firstname,
		"LastName":        lastname,
		"EmailAddress":    email,
		"DateOfBirth":     dateofbirth,
		"IsAuthenticated": true,
		"ItemsInCart":     true,
		"CartItems":       cartItems,
	})

	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}

func (shop *Shop) CheckoutCartHandler(c *gin.Context, store *sessions.CookieStore) {
	var requestBody struct {
		Items []string `json:"itemIDs"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	items := requestBody.Items

	// for each item, id, search item and update its  status to sold
	for i := 0; i < len(items); i++ {
		err := shop.DataAccess.ListAsSold(items[i])
		if err != nil {
			log.Println("CheckoutCartHandler: Could not checkout cart item. Error:", err.Error())
		}
	}

	return
}

func (shop *Shop) RemoveCartItemHandler(c *gin.Context, store *sessions.CookieStore) {
	session, err := store.Get(c.Request, "session")
	if err != nil {
		// User session may not exist, so don't return from error
		log.Println("RemoveCartItemHandler: Error getting session: %s", err.Error())
	}

	var requestBody struct {
		ItemID string `json:"itemID"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	itemID := requestBody.ItemID

	var oldCartItems, updatedCartItems string
	var ok bool
	sessionCartItems := session.Values["CartItems"]
	if sessionCartItems != nil {
		// CartItems exists, get existing list and add new item id to it
		oldCartItems, ok = sessionCartItems.(string)
		if !ok {
			log.Println("RemoveCartItemHandler: Error converting session CartItems value to string")
			return
		}

	} else {
		log.Println("RemoveCartItemHandler: Failed to get Cart Items")
	}

	start := 0
	strlen := len(oldCartItems)
	for i := 0; i < strlen; i++ {
		if oldCartItems[i] != ';' {
			continue
		} else {
			substr := oldCartItems[start:i]
			if substr == itemID {
				// how to remove ; if single item?
				updatedCartItems = oldCartItems[0:start] + oldCartItems[i+1:strlen]
				break
			}
			start = i + 1
		}
	}

	session.Values["CartItems"] = updatedCartItems
	session.Save(c.Request, c.Writer)

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart!"})
}

func (shop *Shop) GetProfilePageHandler(c *gin.Context, store *sessions.CookieStore) {

	// Get session values
	session, err := store.Get(c.Request, "session")
	if err != nil {
		log.Println("GetProfilePageHandler: Error getting session: %s", err.Error())
	}

	firstname := session.Values["FirstName"]
	lastname := session.Values["LastName"]
	email := session.Values["EmailAddress"]
	dateofbirth := session.Values["DateOfBirth"]
	sellerID := session.Values["UserID"]

	var userID int
	var ok bool
	if sellerID != nil {
		userID, ok = sellerID.(int)
		if !ok {
			log.Println("GetProfilePageHandler: Error converting sellerID to string")
			c.String(http.StatusInternalServerError, "Error converting session sellerID to string")
			return
		}
	} else {
		c.String(http.StatusOK, "User ID invalid")
		return
	}

	var dob string
	if dateofbirth != nil {
		dob, ok = dateofbirth.(string)
		if !ok {
			log.Println("GetProfilePageHandler: Error converting date of birth to string")
			c.String(http.StatusInternalServerError, "Error converting session date of birth to string")
			return
		}
	} else {
		c.String(http.StatusOK, "Date of Birth - Invalid Format")
		return
	}

	parsedDate, err := time.Parse("2006-01-02", dob)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}

	formattedDate := parsedDate.Format("02 January 2006")

	items, err := shop.DataAccess.GetItemsBySeller(userID, -1)
	if err != nil {
		log.Println("GetProfilePageHandler: Error retrieving items for seller")
		c.String(http.StatusInternalServerError, " Error retrieving items for seller")
		err = nil
	}

	templates, err := template.ParseFiles("templates/layout.html", "templates/navbar.html", "templates/profile.html")
	if err != nil {
		log.Printf("GetProfilePageHandler: Error parsing templates: %s", err.Error())
	}

	c.Header("Content-Type", "text/html")

	err = templates.ExecuteTemplate(c.Writer, "layout.html", gin.H{
		"FirstName":       firstname,
		"LastName":        lastname,
		"EmailAddress":    email,
		"DateOfBirth":     formattedDate,
		"isAuthenticated": true,
		"Items":           items,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}

func (shop *Shop) GetAddItemPageHandler(c *gin.Context, store *sessions.CookieStore) {

	templates, err := template.ParseFiles("templates/layout.html", "templates/navbar.html", "templates/additem.html")
	if err != nil {
		log.Printf("GetAddItemPageHandler: Error parsing templates: %s", err.Error())
	}

	c.Header("Content-Type", "text/html")

	err = templates.ExecuteTemplate(c.Writer, "layout.html", nil)

	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}

func (shop *Shop) GetEditItemPageHandler(c *gin.Context, store *sessions.CookieStore) {

	session, err := store.Get(c.Request, "session")
	if err != nil {
		// User session may not exist, so don't return from error
		log.Println("GetEditItemPageHandler: Error getting session: %s", err.Error())
	}

	isAuthenticated := session.Values["Authenticated"]

	itemID := c.Param("ID")

	// Retrieve item from database
	itemData, err := shop.DataAccess.GetItem(itemID)
	if err != nil {
		log.Printf("GetEditItemPageHandler: Error getting Item data: %s", err.Error())
		return
	}

	// Clothing size options
	var sizeOptions []string
	bottomSizeOptions := []string{"28", "30", "32", "34", "36"}
	shoeSizeOptions := []string{"3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}
	regularSizeOptions := []string{"XS", "S", "M", "L", "XL"}

	switch itemData.Category {
	case "Bottoms":
		sizeOptions = bottomSizeOptions
	case "Shoes":
		sizeOptions = shoeSizeOptions
	default:
		sizeOptions = regularSizeOptions
	}

	templates, err := template.ParseFiles("templates/layout.html", "templates/navbar.html", "templates/edititem.html")
	if err != nil {
		log.Printf("GetEditItemPageHandler: Error parsing templates: %s", err.Error())
	}

	c.Header("Content-Type", "text/html")

	err = templates.ExecuteTemplate(c.Writer, "layout.html", gin.H{
		"isAuthenticated":      isAuthenticated,
		"ID":                   itemData.ID,
		"itemName":             itemData.Name,
		"itemGender":           itemData.Gender,
		"itemDescription":      itemData.Description,
		"itemPrice":            itemData.Price,
		"itemCategory":         itemData.Category,
		"itemCondition":        itemData.Condition,
		"itemSize":             itemData.Size,
		"genderOptions":        GenderOptions,
		"categoryOptions":      CategoryOptions,
		"conditionOptions":     ConditionOptions,
		"sizeOptions":          sizeOptions,
		"ShowCategoriesBanner": false,
	})

	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}
