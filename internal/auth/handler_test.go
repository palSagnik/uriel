package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/palSagnik/uriel/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUser_Success(t *testing.T) {
	mockRepo := new(MockAuthRepository)

	// Call the repo with a new username and email (we are testing success)
	mockRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(nil, nil)
	mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(nil, nil)
	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("models.User")).Return(nil)

	service := NewService(mockRepo, []byte("test_jwt_here"))
	handler := NewHandler(service)

	router := gin.New()
	router.POST("/auth/register", handler.RegisterUser)

	// making the post request
	payload := models.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password@123",
		Confirm:  "password@123",
	}
	jsonBody, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// parsing the response
	var res models.RegisterResponse
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, "User registered succesfully", res.Message)

	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_UsernameExists(t *testing.T) {
	mockRepo := new(MockAuthRepository)

	// Call the service
	// Expect the repository to find an existing user and return it
	mockRepo.On("GetUserByUsername", mock.Anything, "existingUser").Return(&models.User{Username: "existingUser"}, nil)

	service := NewService(mockRepo, []byte("test_jwt_here"))
	handler := NewHandler(service)

	router := gin.New()
	router.POST("/auth/register", handler.RegisterUser)

	// making the post request
	payload := models.RegisterRequest{
		Username: "existingUser",
		Email:    "newuser@example.com",
		Password: "password@123",
		Confirm:  "password@123",
	}
	jsonBody, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	// parsing the response
	var res models.RegisterResponseFailed
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, "username already exists", res.Error)

	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_EmailExists(t *testing.T) {
	mockRepo := new(MockAuthRepository)

	// Call the service
	// Expect the repository to find an existing email and return it
	mockRepo.On("GetUserByUsername", mock.Anything, "newUser").Return(nil, nil)
	mockRepo.On("GetUserByEmail", mock.Anything, "existinguser@example.com").Return(&models.User{Email: "existinguser@example.com"}, nil)

	service := NewService(mockRepo, []byte("test_jwt_here"))
	handler := NewHandler(service)

	router := gin.New()
	router.POST("/auth/register", handler.RegisterUser)

	// making the post request
	payload := models.RegisterRequest{
		Username: "newUser",
		Email:    "existinguser@example.com",
		Password: "password@123",
		Confirm:  "password@123",
	}
	jsonBody, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	// parsing the response
	var res models.RegisterResponseFailed
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, "email already exists", res.Error)

	mockRepo.AssertExpectations(t)
}
