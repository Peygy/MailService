package service

import (
	"backend/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	AdminService interface {
		GetAllUsers(c *gin.Context)
		DeleteUser(c *gin.Context)
	}

	adminService struct {
		db *gorm.DB
	}
)

func NewAdminService(db *gorm.DB) AdminService {
	return &adminService{
		db: db,
	}
}

func (as *adminService) GetAllUsers(c *gin.Context) {
	var users []model.User
	if err := as.db.Where("role <> ?", model.RoleAdmin).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (as *adminService) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	if err := as.db.Where("id = ?", userID).Delete(&model.User{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
