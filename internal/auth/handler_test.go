package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/palSagnik/uriel/internal/config"
	"github.com/palSagnik/uriel/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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
	var res models.FailedResponse
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
	var res models.FailedResponse
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, "email already exists", res.Error)

	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_Failure(t *testing.T) {
	mockRepo := new(MockAuthRepository)

	// Call the service
	mockRepo.On("GetUserByUsername", mock.Anything, "user").Return(nil, nil)
	mockRepo.On("GetUserByEmail", mock.Anything, "user@example.com").Return(nil, nil)
	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("models.User")).Return(errors.New("internal server error"))

	service := NewService(mockRepo, []byte("test_jwt_here"))
	handler := NewHandler(service)

	router := gin.New()
	router.POST("/auth/register", handler.RegisterUser)

	payload := models.RegisterRequest{
		Username: "user",
		Email:    "user@example.com",
		Password: "password",
		Confirm:  "password",
	}
	jsonPayload, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var res models.FailedResponse
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, "Failed to register user", res.Error)
	
	mockRepo.AssertExpectations(t)
}

func TestLoginPlayer_Success(t *testing.T) {
	mockRepo := new(MockAuthRepository)

	// Create a user model with a pre-hashed password and objectId for the repo to return
	hashed_password, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	testId := "6592008029c8c3e4dc76256c"
    parsedID, _ := primitive.ObjectIDFromHex(testId)

	mockUser := &models.User{
		ID: parsedID,
		Username: "test",
		Email: "test@test.com",
		Password: string(hashed_password),
		Role: config.USER,
		IsOnline: false,
	}
	mockRepo.On("GetUserByUsername", mock.Anything, "test").Return(mockUser, nil)

	service := NewService(mockRepo, []byte("test_jwt_here"))
	handler := NewHandler(service)

	router := gin.New()
	router.POST("/auth/login", handler.LoginUser)

	payload := models.LoginRequest{
		Username: "test",
		Password: "correctpassword",
	}
	jsonPayload, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res models.LoginResponse
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, "User login successful", res.Message)
	assert.NotEmpty(t, res.Token)
	assert.Equal(t, res.UserID, parsedID.Hex())
	
	mockRepo.AssertExpectations(t)
}

func TestLoginPlayer_WrongPassword(t *testing.T) {
	mockRepo := new(MockAuthRepository)

	// Create a user model with a pre-hashed password and objectId for the repo to return
	hashed_password, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	testId := "6592008029c8c3e4dc76256c"
    parsedID, _ := primitive.ObjectIDFromHex(testId)

	mockUser := &models.User{
		ID: parsedID,
		Username: "test",
		Email: "test@test.com",
		Password: string(hashed_password),
		Role: config.USER,
		IsOnline: false,
	}
	mockRepo.On("GetUserByUsername", mock.Anything, "test").Return(mockUser, nil)

	service := NewService(mockRepo, []byte("test_jwt_here"))
	handler := NewHandler(service)

	router := gin.New()
	router.POST("/auth/login", handler.LoginUser)

	payload := models.LoginRequest{
		Username: "test",
		Password: "wrongpassword",
	}
	jsonPayload, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var res models.FailedResponse
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, "invalid username or password", res.Error)
	
	mockRepo.AssertExpectations(t)
}

func TestLoginPlayer_UserDoesnotExist(t *testing.T) {
	mockRepo := new(MockAuthRepository)

	mockRepo.On("GetUserByUsername", mock.Anything, "test").Return(nil, mongo.ErrNoDocuments)

	service := NewService(mockRepo, []byte("test_jwt_here"))
	handler := NewHandler(service)

	router := gin.New()
	router.POST("/auth/login", handler.LoginUser)

	payload := models.LoginRequest{
		Username: "test",
		Password: "password",
	}
	jsonPayload, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var res models.FailedResponse
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, "invalid username or password", res.Error)
	
	mockRepo.AssertExpectations(t)
}
