package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/palSagnik/uriel/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func mockAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", "user-player-id-123")
		c.Set("username", "user-player")
		c.Next()
	}
}

func TestGetAvatars_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAvatarRepo := new(MockAvatarRepository)

	// calling the repo
	testId1 := "6592008029c8c3e4dc76256c"
	testId2 := "66931199abc8c3e4dc2feb6c"
    parsedID1, _ := primitive.ObjectIDFromHex(testId1)
	parsedID2, _ := primitive.ObjectIDFromHex(testId2)

	mockAvatars := []models.Avatar{
		{
			ID:  parsedID1,
			AvatarUrl: "https://uriel.com/avatars/1.png",
		},
		{
			ID:  parsedID2,
			AvatarUrl: "https://uriel.com/avatars/2.png",
		},
	}

	mockAvatarRepo.On("GetAvatars", mock.Anything).Return(mockAvatars, nil)

	service := NewService(mockUserRepo, mockAvatarRepo)
	handler := NewHandler(service)

	router := gin.New()
	router.GET("/users/avatar", mockAuthMiddleware(), handler.GetAllAvatars)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users/avatar", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res models.GetAvatarsResponse
	_ = json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, mockAvatars, res.Avatars)
	mockAvatarRepo.AssertExpectations(t)
}

func TestGetAvatars_InternalServerError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAvatarRepo := new(MockAvatarRepository)

	mockAvatarRepo.On("GetAvatars", mock.Anything).Return(nil, errors.New("avatar list not found"))

	service := NewService(mockUserRepo, mockAvatarRepo)
	handler := NewHandler(service)

	router := gin.New()
	router.GET("/users/avatar", mockAuthMiddleware(), handler.GetAllAvatars)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users/avatar", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var res map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, "failed to get avatars", res["error"])
	mockAvatarRepo.AssertExpectations(t)
}