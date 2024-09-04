/***
-- FILE
@file handler_test.go
@desc Contains tests for the user-related HTTP handlers in the handler package.
@auth John Smith
@v 1.0
@date 01/02/2024
*/

package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"myapp/internal/model"
	"myapp/internal/service"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetAllUsers() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(id int) (model.User, error) {
	args := m.Called(id)
	return args.Get(0).(model.User), args.Error(1)
}

func createTestHandler() *UserHandler {
	mockService := new(MockUserService)
	return &UserHandler{service: mockService}
}

func TestGetAllUsers(t *testing.T) {
	handler := createTestHandler()
	users := []model.User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
	}

	handler.service.(*MockUserService).On("GetAllUsers").Return(users, nil)

	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users", handler.GetAllUsers).Methods("GET")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var gotUsers []model.User
	if err := json.NewDecoder(rr.Body).Decode(&gotUsers); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, users, gotUsers)

	handler.service.(*MockUserService).AssertExpectations(t)
}

func TestGetUserByID(t *testing.T) {
	handler := createTestHandler()
	user := model.User{ID: 1, Name: "John Doe", Email: "john@example.com"}

	handler.service.(*MockUserService).On("GetUserByID", 1).Return(user, nil)

	req, err := http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{id:[0-9]+}", handler.GetUserByID).Methods("GET")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var gotUser model.User
	if err := json.NewDecoder(rr.Body).Decode(&gotUser); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, user, gotUser)

	handler.service.(*MockUserService).AssertExpectations(t)
}

func TestGetUserByID_InvalidID(t *testing.T) {
	handler := createTestHandler()
	handler.service.(*MockUserService).On("GetUserByID", 999).Return(model.User{}, errors.New("User not found"))

	req, err := http.NewRequest("GET", "/users/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{id:[0-9]+}", handler.GetUserByID).Methods("GET")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	assert.Equal(t, "User not found", rr.Body.String())

	handler.service.(*MockUserService).AssertExpectations(t)
}
