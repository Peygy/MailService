package service

import (
	"backend/internal/model"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

type (
	AuthService interface {
		RegisterUser(c *gin.Context)
		Login(c *gin.Context)
	}

	authService struct {
		db *gorm.DB
	}
)

func NewAuthService(db *gorm.DB) AuthService {
	return &authService{
		db: db,
	}
}

func (as *authService) RegisterUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if err := utils.IsCorporateEmail(input.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error hashing password"})
		return
	}
	input.Password = string(hashedPassword)

	user := model.User{
		Email:    input.Email,
		Password: input.Password,
		Role:     model.RoleUser,
	}

	if err := as.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (as *authService) Login(c *gin.Context) {
	var entryUser model.User
	if err := c.ShouldBindJSON(&entryUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	var user model.User
	if err := as.db.Where("email = ?", entryUser.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(entryUser.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
