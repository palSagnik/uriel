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

func (s *Service) RegisterPlayerService(ctx context.Context, req *models.RegisterRequest) (*models.Player, error) {

	// check if this username already exists
	existingPlayerByUsername, err := s.repo.GetPlayerByUsername(ctx, req.Username)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("service: error checking existing username %v", err)
	}
	if existingPlayerByUsername != nil {
		return nil, errors.New("username already exists")
	}

	// check if this email already exists
	existingPlayerByEmail, err := s.repo.GetPlayerByEmail(ctx, req.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("service: error checking existing email %v", err)
	}
	if existingPlayerByEmail != nil {
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
	newPlayer := models.Player{
		ID: primitive.NewObjectID(),
		Username: req.Username,
		Email: req.Email,
		Password: string(hashedPassword),
		Role: config.PLAYER,
		IsOnline: false,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.repo.CreatePlayer(ctx, newPlayer); err != nil {
		return nil, fmt.Errorf("service: error in creating new player %v", err)
	}

	return &newPlayer, nil
}

func (s *Service) LoginPlayerService(ctx context.Context, username string, password string) (string, string, error) {
	
	// retrieve player
	player, err := s.repo.GetPlayerByUsername(ctx, username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", "", errors.New("invalid username or password")
		}
		return "", "", fmt.Errorf("service: error retrieving player %v", err)
	}
	if player == nil {
		return "", "", errors.New("invalid username or password")
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(password)); err != nil {
		return "", "", errors.New("invalid username or password") 
	}

	// update player online status
	if err := s.repo.UpdatePlayerStatus(ctx, player.ID.Hex()); err != nil {
		return "", "", fmt.Errorf("service: error in updating player status %v", err)
	}

	// generate Token
	token, err := s.GenerateToken(player.ID.Hex(), player.Username, player.Role)
	if err != nil {
		return "", "", fmt.Errorf("service: error in generating token %v", err)
	}

	return token, player.ID.Hex(), nil
}

func (s *Service) GenerateToken(playerId, username, role string) (string, error) {
	claims := models.Claims{
		PlayerID: playerId,
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

		c.Set("playerID", claims.PlayerID)
		c.Set("username", claims.Username)

		c.Next()
	}
}