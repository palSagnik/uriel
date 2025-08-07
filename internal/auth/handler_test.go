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
	testReq := models.RegisterRequest{
		Username: "testuser",
		Email: "test@example.com",
		Password: "password@123",
		Confirm: "password@123",
	}
	jsonBody, _ := json.Marshal(testReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// parsing the response
	var resp models.RegisterResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "User registered succesfully", resp.Message)
	
	mockRepo.AssertExpectations(t)
}	