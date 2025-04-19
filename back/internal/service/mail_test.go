package service

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func FuzzMailService_GetInboxMails(f *testing.F) {
	rand.Seed(time.Now().UnixNano())
	f.Add(uint(1))
	f.Add(uint(100000))
	f.Add(uint(rand.Uint32()))

	f.Fuzz(func(t *testing.T, userID uint) {
		mockDB := new(MockMailDB)
		service := NewMailService(mockDB)

		mockDB.On("Where", "id = ?", userID).Return(mockDB)
		mockDB.On("First", mock.AnythingOfType("*model.User")).Return(mockDB)

		if rand.Intn(2) == 0 {
			mockDB.On("Error").Return(nil)
			mockDB.On("Find", mock.AnythingOfType("*[]model.Mail")).Return(mockDB)
			mockDB.On("Error").Return(nil)
		} else {
			mockDB.On("Error").Return(assert.AnError)
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)

		service.GetInboxMails(c)

		validCodes := []int{
			http.StatusOK,
			http.StatusUnauthorized,
			http.StatusInternalServerError,
		}
		assert.Contains(t, validCodes, w.Code, "unexpected status code")
		mockDB.AssertExpectations(t)
	})
}

func FuzzMailService_SendMail(f *testing.F) {
	rand.Seed(time.Now().UnixNano())
	f.Add(uint(1), "test@example.com", "subject", "body")
	f.Add(uint(0), "", "", "")
	f.Add(uint(rand.Uint32()), generateRandomString(10), generateRandomString(20), generateRandomString(50))

	f.Fuzz(func(t *testing.T, userID uint, receiver, subject, body string) {
		mockDB := new(MockMailDB)
		service := NewMailService(mockDB)

		mockDB.On("Where", "id = ?", userID).Return(mockDB)
		mockDB.On("First", mock.AnythingOfType("*model.User")).Return(mockDB)

		if rand.Intn(2) == 0 {
			mockDB.On("Error").Return(nil)
			mockDB.On("Create", mock.AnythingOfType("*model.Mail")).Return(mockDB)
			mockDB.On("Error").Return(nil)
		} else {
			mockDB.On("Error").Return(assert.AnError)
		}

		input := map[string]interface{}{
			"receivers": []string{receiver},
			"subject":   subject,
			"body":      body,
		}
		jsonData, _ := json.Marshal(input)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Request = httptest.NewRequest(http.MethodPost, "/send", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		service.SendMail(c)

		validCodes := []int{
			http.StatusCreated,
			http.StatusBadRequest,
			http.StatusUnauthorized,
			http.StatusInternalServerError,
		}
		assert.Contains(t, validCodes, w.Code, "unexpected status code")
	})
}

func FuzzMailService_ArchiveMail(f *testing.F) {
	rand.Seed(time.Now().UnixNano())
	f.Add(uint(1), "1")
	f.Add(uint(0), "invalid_id")
	f.Add(uint(rand.Uint32()), strconv.Itoa(rand.Int()))

	f.Fuzz(func(t *testing.T, userID uint, mailID string) {
		mockDB := new(MockMailDB)
		service := NewMailService(mockDB)

		mockDB.On("Where", "user_id = ?", userID).Return(mockDB)
		mockDB.On("First", mock.AnythingOfType("*model.Trash")).Return(mockDB)

		if rand.Intn(2) == 0 {
			mockDB.On("Error").Return(nil)
			mockDB.On("Model", mock.AnythingOfType("*model.Trash")).Return(mockDB)
			mockDB.On("Where", "user_id = ?", userID).Return(mockDB)
			mockDB.On("Update", "archived", mock.Anything).Return(mockDB)
			mockDB.On("Error").Return(nil)
		} else {
			mockDB.On("Error").Return(assert.AnError)
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: mailID}}

		service.ArchiveMail(c)

		validCodes := []int{
			http.StatusOK,
			http.StatusBadRequest,
			http.StatusInternalServerError,
		}
		assert.Contains(t, validCodes, w.Code, "unexpected status code")
		mockDB.AssertExpectations(t)
	})
}

// func FuzzMailService_DeleteMail(f *testing.F) {
// 	rand.Seed(time.Now().UnixNano())
// 	f.Add(uint(1), "1")
// 	f.Add(uint(0), "invalid_id")
// 	f.Add(uint(rand.Uint32()), strconv.Itoa(rand.Int()))

// 	f.Fuzz(func(t *testing.T, userID uint, mailID string) {
// 		mockDB := new(MockMailDB)
// 		service := NewMailService(mockDB)

// 		mockDB.On("Where", "user_id = ?", userID).Return(mockDB)

// 		tr := model.Trash{
// 			Archived: []int64{},
// 			Deleted:  []int64{},
// 		}

// 		if rand.Intn(2) == 0 {
// 			id, _ := strconv.ParseInt(mailID, 10, 64)
// 			tr.Archived = append(tr.Archived, id)
// 		}

// 		mockDB.On("First", mock.AnythingOfType("*model.Trash")).Run(func(args mock.Arguments) {
// 			arg := args.Get(0).(*model.Trash)
// 			*arg = tr
// 		}).Return(mockDB)

// 		if rand.Intn(2) == 0 {
// 			mockDB.On("Error").Return(nil)
// 			mockDB.On("Model", mock.AnythingOfType("*model.Trash")).Return(mockDB)
// 			mockDB.On("Where", "user_id = ?", userID).Return(mockDB)
// 			mockDB.On("Update", "archived", mock.Anything).Return(mockDB)
// 			mockDB.On("Update", "deleted", mock.Anything).Return(mockDB)
// 			mockDB.On("Error").Return(nil)
// 		} else {
// 			mockDB.On("Error").Return(assert.AnError)
// 		}

// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)
// 		c.Set("userID", userID)
// 		c.Params = gin.Params{gin.Param{Key: "id", Value: mailID}}

// 		service.DeleteMail(c)

// 		validCodes := []int{
// 			http.StatusOK,
// 			http.StatusBadRequest,
// 			http.StatusInternalServerError,
// 		}
// 		assert.Contains(t, validCodes, w.Code, "unexpected status code")
// 		mockDB.AssertExpectations(t)
// 	})
// }
