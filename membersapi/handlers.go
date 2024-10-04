package membersapi

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-crypt/crypt"
	"github.com/go-crypt/crypt/algorithm"
	"github.com/go-crypt/crypt/algorithm/argon2"
	"github.com/gorilla/sessions"
)

type UserDatabase struct {
	DataAccess DataAccess
}

func (userDatabase *UserDatabase) CreateUserHandler(c *gin.Context, store *sessions.CookieStore) {
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

	c.JSON(http.StatusOK, nil)
}

func (userDatabase *UserDatabase) LoginHandler(c *gin.Context, store *sessions.CookieStore) {
	emailAddr := c.PostForm("EmailAddress")
	password := c.PostForm("Password")

	userExists, err := userDatabase.DataAccess.CheckUserExists(emailAddr)
	if err != nil {
		fmt.Printf("LoginHandler: An error occurred checking if email exists\n")
	}

	var errors map[string]string
	errors = make(map[string]string, 2)
	if !userExists {
		errors["email"] = "Email address doesn't exist"
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	var (
		decoder *crypt.Decoder
		digest  algorithm.Digest
	)

	userData, err := userDatabase.DataAccess.GetUser(emailAddr)
	if err != nil {
		fmt.Printf("LoginHandler: Something went wrong with getting the user: %s", err.Error())
	}

	if decoder, err = crypt.NewDefaultDecoder(); err != nil {
		panic(err)
	}

	if digest, err = decoder.Decode(userData.Hash); err != nil {
		panic(err)
	}

	if !digest.Match(password) {
		errors["password"] = "Wrong password. Please try again."
		c.JSON(http.StatusBadRequest, gin.H{
			"errors":          errors,
			"isAuthenticated": false,
		})
		return
	}

	session, _ := store.Get(c.Request, "session")
	session.Values["Authenticated"] = true
	session.Values["UserID"] = userData.UserID
	session.Values["FirstName"] = userData.FirstName
	session.Values["LastName"] = userData.LastName
	session.Values["EmailAddress"] = userData.EmailAddress
	session.Values["DateOfBirth"] = userData.DateOfBirth
	session.Save(c.Request, c.Writer)

	c.JSON(http.StatusOK, gin.H{
		"isAuthenticated": true,
	})
}

func (UserDatabase *UserDatabase) LogoutHandler(c *gin.Context, store *sessions.CookieStore) {
	session, _ := store.Get(c.Request, "session")
	for k := range session.Values {
		delete(session.Values, k)
	}

	session.Save(c.Request, c.Writer)
	c.Redirect(http.StatusSeeOther, "/shopapi/")
}

func (userDatabase *UserDatabase) EditUserFirstNameHandler(c *gin.Context, store *sessions.CookieStore) {

	session, err := store.Get(c.Request, "session")
	if err != nil {
		fmt.Println("EditUserFirstNameHandler error: %s", err)
		return
	}

	userID := session.Values["UserID"].(int)

	firstName := c.PostForm("firstname-input")
	if !validateName(firstName) {
		errorMessage := "First Name is not valid. Please try again."
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": errorMessage})
		return
	}

	err = userDatabase.DataAccess.UpdateUserFirstName(userID, firstName)

	if err != nil {
		fmt.Println(err)
		return
	}

	session.Values["FirstName"] = firstName
	session.Save(c.Request, c.Writer)
}

func (userDatabase *UserDatabase) EditUserLastNameHandler(c *gin.Context, store *sessions.CookieStore) {

	session, err := store.Get(c.Request, "session")
	if err != nil {
		fmt.Println("EditUserLastNameHandler error: %s", err)
		return
	}

	userID := session.Values["UserID"].(int)

	lastName := c.PostForm("lastname-input")
	if !validateName(lastName) {
		errorMessage := "Last Name is not valid. Please try again."
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": errorMessage})
		return
	}

	err = userDatabase.DataAccess.UpdateUserLastName(userID, lastName)

	if err != nil {
		fmt.Println(err)
		return
	}

	session.Values["LastName"] = lastName
	session.Save(c.Request, c.Writer)
}

func (userDatabase *UserDatabase) EditUserDateOfBirthHandler(c *gin.Context, store *sessions.CookieStore) {

	session, err := store.Get(c.Request, "session")
	if err != nil {
		fmt.Println("EditUserDateOfBirthHandler error: %s", err)
		return
	}

	userID := session.Values["UserID"].(int)

	dateOfBirth := c.PostForm("dateofbirth-input")
	err = userDatabase.DataAccess.UpdateUserDateOfBirth(userID, dateOfBirth)

	if err != nil {
		fmt.Println(err)
		return
	}

	session.Values["DateOfBirth"] = dateOfBirth
	session.Save(c.Request, c.Writer)
}

func (userDatabase *UserDatabase) EditUserEmailHandler(c *gin.Context, store *sessions.CookieStore) {

	session, err := store.Get(c.Request, "session")
	if err != nil {
		fmt.Println("EditUserEmailHandler error: %s", err)
		return
	}

	userID := session.Values["UserID"].(int)

	email := c.PostForm("emailaddress-input")
	if len(email) < 3 {
		errorMessage := "Email address must be at least 3 characters long."
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": errorMessage})
		return
	}

	err = userDatabase.DataAccess.UpdateUserEmail(userID, email)

	if err != nil {
		fmt.Println(err)
		return
	}

	session.Values["EmailAddress"] = email
	session.Save(c.Request, c.Writer)
}

func (userDatabase *UserDatabase) EditUserPasswordHandler(c *gin.Context, store *sessions.CookieStore) {

	session, err := store.Get(c.Request, "session")
	if err != nil {
		fmt.Println("EditUserPasswordHandler error: %s", err)
		return
	}

	userID := session.Values["UserID"].(int)

	password := c.PostForm("password-input")
	if !validatePassword(password) {
		errorMessage := "Password incorrect format."
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": errorMessage})
		return
	}

	err = userDatabase.DataAccess.UpdateUserPassword(userID, password)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func (userDatabase *UserDatabase) DeleteUserHandler(c *gin.Context, store *sessions.CookieStore) {

	session, err := store.Get(c.Request, "session")
	if err != nil {
		fmt.Println("EditUserPasswordHandler error: %s", err)
		return
	}

	userID := session.Values["UserID"].(int)

	userDatabase.DataAccess.DeleteUser(userID)

	c.HTML(http.StatusOK, "deleteditem.html", nil)
}

// Get templated files
func (UserDatabase *UserDatabase) GetSignUpPageHandler(c *gin.Context) {

	templates, err := template.ParseFiles("templates/layout.html", "templates/navbar.html", "templates/signup.html")
	if err != nil {
		log.Printf("GetSignUpPageHandler: Error parsing templates: %s", err.Error())
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

func (UserDatabase *UserDatabase) GetProfilePageHandler(c *gin.Context, store *sessions.CookieStore) {

	// Get session values
	session, err := store.Get(c.Request, "session")
	if err != nil {
		log.Println("GetProfilePageHandler: Error getting session: %s", err.Error())
	}

	firstname := session.Values["FirstName"]
	lastname := session.Values["LastName"]
	email := session.Values["EmailAddress"]
	dateofbirth := session.Values["DateOfBirth"]

	templates, err := template.ParseFiles("templates/layout.html", "templates/navbar.html", "templates/profile.html")
	if err != nil {
		log.Printf("GetProfilePageHandler: Error parsing templates: %s", err.Error())
	}

	c.Header("Content-Type", "text/html")

	// Execute the main layout template with the "signup" content embedded
	err = templates.ExecuteTemplate(c.Writer, "layout.html", gin.H{
		"FirstName":       firstname,
		"LastName":        lastname,
		"EmailAddress":    email,
		"DateOfBirth":     dateofbirth,
		"isAuthenticated": true,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}

func (UserDatabase *UserDatabase) GetEditPageHandler(c *gin.Context, store *sessions.CookieStore) {

	// Get session values
	session, err := store.Get(c.Request, "session")
	if err != nil {
		log.Println("GetProfilePageHandler: Error getting session: %s", err.Error())
	}

	firstname := session.Values["FirstName"]
	lastname := session.Values["LastName"]
	email := session.Values["EmailAddress"]
	dateofbirth := session.Values["DateOfBirth"]

	templates, err := template.ParseFiles("templates/layout.html", "templates/navbar.html", "templates/editprofile.html")
	if err != nil {
		log.Printf("GetProfilePageHandler: Error parsing templates: %s", err.Error())
	}

	c.Header("Content-Type", "text/html")

	// Execute the main layout template with the "signup" content embedded
	err = templates.ExecuteTemplate(c.Writer, "layout.html", gin.H{
		"FirstName":       firstname,
		"LastName":        lastname,
		"EmailAddress":    email,
		"DateOfBirth":     dateofbirth,
		"isAuthenticated": true,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}

func validatePassword(password string) bool {
	passLen := len(password)
	var containsDigits bool
	var containsCaps bool

	if passLen < 6 {
		return false
	}

	for i := 0; i < passLen; i++ {
		if password[i] >= '0' && password[i] <= '9' {
			containsDigits = true
		}

		if password[i] >= 'A' && password[i] <= 'Z' {
			containsCaps = true
		}
	}

	if !containsDigits || !containsCaps {
		return false
	}

	return true
}

func validateName(name string) bool {
	nameLen := len(name)

	if nameLen < 2 {
		return false
	}

	for i := 0; i < nameLen; i++ {
		if !(name[i] >= 'a' && name[i] <= 'z') && !(name[i] >= 'A' && name[i] <= 'Z') {
			return false
		}

	}

	return true
}
