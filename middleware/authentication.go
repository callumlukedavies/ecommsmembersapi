package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

func AuthorizeUser(store *sessions.CookieStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request, "session")
		if err != nil {
			log.Printf("AuthRequired: Error in getting new cookie session. Error: %s", err.Error())
			c.String(http.StatusUnauthorized, "Please login to view this page.")
			c.Abort()
			return
		}

		if auth, ok := session.Values["Authenticated"].(bool); !ok || !auth {
			c.String(http.StatusUnauthorized, "Please login to view this page.")
			c.Abort()
			return
		}

		c.Next()
	}
}
