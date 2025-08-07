package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/palSagnik/uriel/internal/config"
	"github.com/palSagnik/uriel/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo AuthRepository
	jwtSecretKey []byte
}

func NewService(repo AuthRepository, jwtSecretKey []byte) *Service {
	return &Service{
		repo: repo,
		jwtSecretKey: jwtSecretKey,
	}
}

func (s *Service) RegisterUserService(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {

	// check if this username already exists
	existingUserByUsername, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("service: error checking existing username %v", err)
	}
	if existingUserByUsername != nil {
		return nil, errors.New("username already exists")
	}

	// check if this email already exists
	existingUserByEmail, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("service: error checking existing email %v", err)
	}
	if existingUserByEmail != nil {
		return nil, errors.New("email already exists")
	}

	// hash password
	// does not accept more than 72 bytes
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("service: error hashing password %v", err)
	}

	// create player
	// TODO: Errors should be ENUMS
	newUser := models.User{
		ID: primitive.NewObjectID(),
		Username: req.Username,
		Email: req.Email,
		Password: string(hashedPassword),
		Role: config.USER,
		IsOnline: false,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.repo.CreateUser(ctx, newUser); err != nil {
		return nil, fmt.Errorf("service: error in creating new user %v", err)
	}

	return &newUser, nil
}

func (s *Service) LoginUserService(ctx context.Context, username string, password string) (string, string, error) {

	// retrieve user
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", "", errors.New("invalid username or password")
		}
		return "", "", fmt.Errorf("service: error retrieving player %v", err)
	}
	if user == nil {
		return "", "", errors.New("invalid username or password")
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("invalid username or password") 
	}

	// update user online status
	if err := s.repo.UpdateUserStatus(ctx, user.ID.Hex()); err != nil {
		return "", "", fmt.Errorf("service: error in updating user status %v", err)
	}

	// generate Token
	token, err := s.GenerateToken(user.ID.Hex(), user.Username, user.Role)
	if err != nil {
		return "", "", fmt.Errorf("service: error in generating token %v", err)
	}

	return token, user.ID.Hex(), nil
}

func (s *Service) GenerateToken(userId, username, role string) (string, error) {
	claims := models.Claims{
		UserID: userId,
		Username: username,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims {
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * config.TOKEN_DURATION)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer: "uriel",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token %v", err)
	}

	return tokenString, nil
}

func (s *Service) ValidateToken(tokenString string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
            return nil, errors.New("token is malformed")
        } else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
            return nil, errors.New("token has expired or is not yet valid")
        }
        return nil, fmt.Errorf("token parsing failed: %w", err)
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims or token is not valid")
	}

	return claims, nil
}

func (s *Service) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get the token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorisation header",
			})
			c.Abort()
			return
		}
		
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorisation header",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// validate token
		claims, err := s.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("playerID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}