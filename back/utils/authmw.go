package utils

import (
	"backend/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type BasicAuthMiddleware struct {
	db *gorm.DB
}

func NewBasicAuthMiddleware(db *gorm.DB) *BasicAuthMiddleware {
	return &BasicAuthMiddleware{
		db: db,
	}
}

func (mw *BasicAuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization required"})
			c.Abort()
			return
		}

		var user model.User
		if err := mw.db.Where("email = ?", username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email"})
			c.Abort()
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid password"})
			c.Abort()
			return
		}

		c.Set("userID", user.Id)
		c.Next()
	}
}
