package service

import (
	"backend/internal/model"
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func FuzzAuthService_RegisterUser(f *testing.F) {
	rand.Seed(time.Now().UnixNano())

	f.Add("valid@example.com", "strongPassword123")
	f.Add("invalid-email", "short")
	f.Add("", "")
	f.Add("another@test.com", "password with spaces")
	f.Add(generateRandomString(10)+"@test.com", generateRandomString(20))

	f.Fuzz(func(t *testing.T, email, password string) {
		mockDB := new(MockMailDB)
		service := NewAuthService(mockDB)

		mockDB.On("Model", mock.AnythingOfType("*model.User")).Return(mockDB)
		mockDB.On("Select", mock.Anything, mock.Anything).Return(mockDB)
		mockDB.On("Where", "email = ?", email).Return(mockDB)

		if rand.Intn(2) == 0 {
			mockDB.On("Find", mock.Anything, mock.Anything).Return(mockDB).Run(func(args mock.Arguments) {
				exists := args.Get(0).(*bool)
				*exists = false
			})
		} else {
			mockDB.On("Find", mock.Anything, mock.Anything).Return(mockDB).Run(func(args mock.Arguments) {
				exists := args.Get(0).(*bool)
				*exists = true
			})
		}

		if rand.Intn(2) == 0 {
			mockDB.On("Create", mock.AnythingOfType("*model.User")).Return(mockDB)
			mockDB.On("Create", mock.AnythingOfType("*model.Trash")).Return(mockDB)
			mockDB.On("Error").Return(nil)
		} else {
			mockDB.On("Create", mock.AnythingOfType("*model.User")).Return(mockDB)
			mockDB.On("Error").Return(assert.AnError)
		}

		input := map[string]string{
			"email":    email,
			"password": password,
		}
		jsonData, _ := json.Marshal(input)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		service.RegisterUser(c)

		validCodes := []int{
			http.StatusCreated,
			http.StatusBadRequest,
			http.StatusUnauthorized,
			http.StatusInternalServerError,
		}
		assert.Contains(t, validCodes, w.Code, "unexpected status code")
	})
}

func FuzzAuthService_Login(f *testing.F) {
	rand.Seed(time.Now().UnixNano())

	f.Add("user@example.com", "correct_password")
	f.Add("", "")
	f.Add("unknown@user.com", "wrong_password")
	f.Add(generateRandomString(5)+"@test.com", generateRandomString(15))

	f.Fuzz(func(t *testing.T, email, password string) {
		mockDB := new(MockMailDB)
		service := NewAuthService(mockDB)

		mockDB.On("Where", "email = ?", email).Return(mockDB)

		if rand.Intn(2) == 0 {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), bcrypt.DefaultCost)
			mockDB.On("First", mock.Anything, mock.Anything).Return(mockDB).Run(func(args mock.Arguments) {
				user := args.Get(0).(*model.User)
				*user = model.User{
					Email:    email,
					Password: string(hashedPassword),
				}
			})
			mockDB.On("Error").Return(nil)
		} else {
			mockDB.On("First", mock.Anything, mock.Anything).Return(mockDB)
			mockDB.On("Error").Return(assert.AnError)
		}

		input := map[string]string{
			"email":    email,
			"password": password,
		}
		jsonData, _ := json.Marshal(input)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		service.Login(c)

		validCodes := []int{
			http.StatusOK,
			http.StatusBadRequest,
			http.StatusUnauthorized,
		}
		assert.Contains(t, validCodes, w.Code, "unexpected status code")
		mockDB.AssertExpectations(t)
	})
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
