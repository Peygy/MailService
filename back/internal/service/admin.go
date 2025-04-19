package service

import (
	"backend/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	AdminService interface {
		GetAllUsers(c *gin.Context)
		DeleteUser(c *gin.Context)
		GetAllMails(c *gin.Context)
		DeleteMail(c *gin.Context)
	}

	adminService struct {
		db model.MailDB
	}
)

func NewAdminService(db model.MailDB) AdminService {
	return &adminService{
		db: db,
	}
}

func (as *adminService) GetAllUsers(c *gin.Context) {
	var users []model.User
	if err := as.db.Where("role <> ?", model.RoleAdmin).Find(&users).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (as *adminService) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	if err := as.db.Where("id = ?", userID).Delete(&model.User{}).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (as *adminService) GetAllMails(c *gin.Context) {
	var mails []model.Mail
	if err := as.db.Find(&mails).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching mails"})
		return
	}
	c.JSON(http.StatusOK, mails)
}

func (as *adminService) DeleteMail(c *gin.Context) {
	mailID := c.Param("id")

	if err := as.db.Where("id = ?", mailID).Delete(&model.Mail{}).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting mail"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
