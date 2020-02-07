package middleware

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Auth(c *gin.Context) {
	session, err := c.Cookie("session")
	if err != nil || session == "" {
		c.String(http.StatusUnauthorized, "")
		c.Abort()
	} else {
		claims, err := models.ParseToken(session)
		if err != nil {
			c.String(http.StatusUnauthorized, "")
			c.Abort()
		}
		c.Set("isAdmin", claims.IsAdmin)
		c.Set("id", claims.ID)
		c.Next()
	}
}

func AdminOnly(c *gin.Context) {
	if !c.GetBool("isAdmin") {
		c.String(http.StatusForbidden, "")
		c.Abort()
	} else {
		c.Next()
	}
}

func CheckUserExist(c *gin.Context) {
	id := c.GetInt64("id")
	if !models.IsObjectExistsByID(&models.User{}, uint64(id)) {
		c.String(http.StatusForbidden, "")
		c.Abort()
	} else {
		c.Next()
	}
}
