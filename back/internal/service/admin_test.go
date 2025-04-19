package service

import (
	"backend/internal/model"
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

type MockMailDB struct {
	mock.Mock
}

func (m *MockMailDB) Model(value interface{}) (tx model.MailDB) {
	args := m.Called(value)
	return args.Get(0).(model.MailDB)
}

func (m *MockMailDB) Select(query interface{}, args ...interface{}) (tx model.MailDB) {
	callArgs := make([]interface{}, 0)
	callArgs = append(callArgs, query)
	callArgs = append(callArgs, args...)
	return m.Called(callArgs...).Get(0).(model.MailDB)
}

func (m *MockMailDB) Create(value interface{}) (tx model.MailDB) {
	return m.Called(value).Get(0).(model.MailDB)
}

func (m *MockMailDB) Update(column string, value interface{}) (tx model.MailDB) {
	return m.Called(column, value).Get(0).(model.MailDB)
}

func (m *MockMailDB) Delete(value interface{}, conds ...interface{}) (tx model.MailDB) {
	callArgs := make([]interface{}, 0)
	callArgs = append(callArgs, value)
	callArgs = append(callArgs, conds...)
	return m.Called(callArgs...).Get(0).(model.MailDB)
}

func (m *MockMailDB) Where(query interface{}, args ...interface{}) (tx model.MailDB) {
	callArgs := make([]interface{}, 0)
	callArgs = append(callArgs, query)
	callArgs = append(callArgs, args...)
	return m.Called(callArgs...).Get(0).(model.MailDB)
}

func (m *MockMailDB) Find(dest interface{}, conds ...interface{}) (tx model.MailDB) {
	callArgs := make([]interface{}, 0)
	callArgs = append(callArgs, dest)
	callArgs = append(callArgs, conds...)
	return m.Called(callArgs...).Get(0).(model.MailDB)
}

func (m *MockMailDB) First(dest interface{}, conds ...interface{}) (tx model.MailDB) {
	callArgs := make([]interface{}, 0)
	callArgs = append(callArgs, dest)
	callArgs = append(callArgs, conds...)
	return m.Called(callArgs...).Get(0).(model.MailDB)
}

func (m *MockMailDB) Error() error {
	return m.Called().Error(0)
}

func FuzzAdminService_DeleteUser(f *testing.F) {
	f.Add("123")
	f.Add("999")
	f.Add("0")
	f.Add("abc")

	f.Fuzz(func(t *testing.T, userID string) {
		mockDB := new(MockMailDB)
		service := NewAdminService(mockDB)

		mockDB.On("Where", "id = ?", userID).Return(mockDB)
		mockDB.On("Delete", &model.User{}).Return(mockDB)

		if rand.Intn(2) == 0 {
			mockDB.On("Error").Return(nil)
		} else {
			mockDB.On("Error").Return(assert.AnError)
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: userID}}

		service.DeleteUser(c)

		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
		mockDB.AssertExpectations(t)
	})
}

func FuzzAdminService_DeleteMail(f *testing.F) {
	rand.Seed(time.Now().UnixNano())
	f.Add("1")
	f.Add("100000")
	f.Add("invalid_id")
	f.Add(strconv.Itoa(rand.Int()))

	f.Fuzz(func(t *testing.T, mailID string) {
		mockDB := new(MockMailDB)
		service := NewAdminService(mockDB)

		mockDB.On("Where", "id = ?", mailID).Return(mockDB)
		mockDB.On("Delete", &model.Mail{}).Return(mockDB)

		if rand.Intn(2) == 0 {
			mockDB.On("Error").Return(nil)
		} else {
			mockDB.On("Error").Return(assert.AnError)
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: mailID}}

		service.DeleteMail(c)

		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
		mockDB.AssertExpectations(t)
	})
}
