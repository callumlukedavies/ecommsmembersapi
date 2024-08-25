package membersapi

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-crypt/crypt"
	"github.com/go-crypt/crypt/algorithm"
	"github.com/go-crypt/crypt/algorithm/argon2"
)

type UserDatabase struct {
	DataAccess DataAccess
}

func (userDatabase *UserDatabase) CreateUserHandler(c *gin.Context) {
	firstname := c.PostForm("firstname")
	lastname := c.PostForm("lastname")
	email := c.PostForm("emailaddress")
	dob := c.PostForm("dateofbirth")
	password := c.PostForm("password")

	var errors map[string]string
	errors = make(map[string]string, 5)
	if len(firstname) < 3 {
		errors["firstname"] = "First Name must be at least 3 characters long."
	}

	if len(lastname) < 3 {
		errors["lastname"] = "Last Name must be at least 3 characters long."
	}

	userExists, err := userDatabase.DataAccess.CheckUserExists(email)
	if err != nil {
		fmt.Printf("CreateUserHandler: An error occurred checking if email exists\n")
	}

	if userExists {
		errors["email"] = "Email address already exists."
	}

	if len(errors) > 0 {
		// Return validation errors
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	var hasher *argon2.Hasher
	var digest algorithm.Digest

	if hasher, err = argon2.New(
		argon2.WithProfileRFC9106LowMemory(),
	); err != nil {
		panic(err)
	}

	if digest, err = hasher.Hash(password); err != nil {
		panic(err)
	}

	encodedPass := digest.Encode()

	err = userDatabase.DataAccess.CreateUser(firstname, lastname, email, dob, encodedPass)
	if err != nil {
		fmt.Printf("CreateUserHandler: Error occurred adding new user to database: %s\n", err.Error())
		errors["general"] = "An unexpected error occurred. Please try again later."
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	c.Redirect(http.StatusSeeOther, "/shopapi/")
}

func (userDatabase *UserDatabase) UpdateUserHandler(c *gin.Context) {
	userKey := c.PostForm("price")
	paramID := c.Param("ID")
	userValue := c.Param("Name")
	userID, _ := strconv.ParseInt(paramID, 10, 64)
	err := userDatabase.DataAccess.UpdateUserData(userID, userKey, userValue)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func (userDatabase *UserDatabase) GetUserHandler(c *gin.Context) {
	emailAddr := c.PostForm("EmailAddress")
	password := c.PostForm("Password")

	var (
		decoder *crypt.Decoder
		err     error
		digest  algorithm.Digest
	)

	userData, err := userDatabase.DataAccess.GetUser(emailAddr)
	if err != nil {
		fmt.Printf("Something went wrong with getting the user: %s", err)
	}

	if decoder, err = crypt.NewDefaultDecoder(); err != nil {
		panic(err)
	}

	if digest, err = decoder.Decode(userData.Hash); err != nil {
		panic(err)
	}

	if digest.Match(password) {
		// correct password, redirect to profile page
	} else {
		// incorrect password, update form
	}

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"user": userData,
	})
}

func (userDatabase *UserDatabase) DeleteUserHandler(c *gin.Context) {

	param := c.Param("ID")
	id, _ := strconv.ParseInt(param, 10, 64)
	userDatabase.DataAccess.DeleteUser(id)

	c.HTML(http.StatusOK, "deleteditem.html", nil)
}

func (UserDatabase *UserDatabase) GetSignUpPageHandler(c *gin.Context) {

	// tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/navbar.html", "templates/signup.html"))
	// c.Header("Content-Type", "text/html")
	// tmpl.ExecuteTemplate(c.Writer, "layout", gin.H{
	// 	"Title": "Sign Up",
	// })

	templates, err := template.ParseFiles("templates/layout.html", "templates/navbar.html", "templates/signup.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	c.Header("Content-Type", "text/html")

	// Execute the main layout template with the "signup" content embedded
	err = templates.ExecuteTemplate(c.Writer, "layout.html", gin.H{
		"Title": "Sign Up",
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}
