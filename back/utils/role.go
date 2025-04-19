package utils

import (
	"backend/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleMiddleware struct {
	db model.MailDB
}

func NewRoleMiddleware(db model.MailDB) *RoleMiddleware {
	return &RoleMiddleware{
		db: db,
	}
}

func (rr *RoleMiddleware) Middleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)

		var user model.User
		if err := rr.db.First(&user, userID).Error(); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "User not found"})
			c.Abort()
			return
		}

		if user.Role != role {
			c.JSON(http.StatusForbidden, gin.H{"message": "Access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
