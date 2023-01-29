package helpers

import (
	"context"
	"fmt"
	"ima-svc-management/config"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

type TokenDetail struct {
	AccessToken         string
	RefreshToken        string
	AccessUuid          string
	RefreshUuid         string
	ActiveTokenExpires  int64
	RefreshTokenExpires int64
}

type AccessDetail struct {
	AccessUUID string
	Email      string
}

type Token struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

const ACCESS_TOKEN_EXPIRATION = time.Minute * 15
const REFRESH_TOKEN_EXPIRATION = time.Hour * 24

type Auth struct{}

func (auth Auth) CreateToken(email string) (*TokenDetail, error) {

	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	tokenDetail := &TokenDetail{}
	tokenDetail.ActiveTokenExpires = time.Now().Add(ACCESS_TOKEN_EXPIRATION).Unix()
	tokenDetail.AccessUuid = uuid.New().String()

	tokenDetail.RefreshTokenExpires = time.Now().Add(REFRESH_TOKEN_EXPIRATION).Unix()
	tokenDetail.RefreshUuid = uuid.New().String()

	activeTokenClaims := jwt.MapClaims{}
	activeTokenClaims["authorized"] = true
	activeTokenClaims["access_uuid"] = tokenDetail.AccessUuid
	activeTokenClaims["email"] = email
	activeTokenClaims["exp"] = tokenDetail.ActiveTokenExpires
	activeToken := jwt.NewWithClaims(jwt.SigningMethodHS256, activeTokenClaims)
	tokenDetail.AccessToken, err = activeToken.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}

	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["refresh_uuid"] = tokenDetail.RefreshUuid
	refreshTokenClaims["email"] = email
	refreshTokenClaims["exp"] = tokenDetail.RefreshTokenExpires
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	tokenDetail.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}
	return tokenDetail, nil
}

func (auth Auth) CreateAuth(ctx context.Context, email string, tokenDetail *TokenDetail) error {
	activeToken := time.Unix(tokenDetail.ActiveTokenExpires, 0)
	refreshToken := time.Unix(tokenDetail.RefreshTokenExpires, 0)
	now := time.Now()

	err := config.RedisClient.Set(ctx, tokenDetail.AccessUuid, email, activeToken.Sub(now)).Err()
	if err != nil {
		return err
	}
	err = config.RedisClient.Set(ctx, tokenDetail.RefreshUuid, email, refreshToken.Sub(now)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (auth Auth) ExtractToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	splitToken := strings.Split(token, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}

func (auth Auth) VerifyToken(c *gin.Context) (*jwt.Token, error) {
	tokenString := auth.ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (auth Auth) VerifyRefreshToken(c *gin.Context, refreshToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_TOKEN_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (auth Auth) TokenValid(c *gin.Context) error {
	token, err := auth.VerifyToken(c)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return err
	}
	return nil
}

func (auth Auth) ExtractTokenMetadata(c *gin.Context) (*AccessDetail, error) {
	token, err := auth.VerifyToken(c)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		return &AccessDetail{
			AccessUUID: accessUUID,
			Email:      claims["email"].(string),
		}, nil
	}
	return nil, err
}

func (auth Auth) FetchAuth(ctx context.Context, accessDetail *AccessDetail) (int64, error) {
	userid, err := config.RedisClient.Get(ctx, accessDetail.AccessUUID).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseInt(userid, 10, 64)
	return userID, nil
}

func (auth Auth) DeleteAuth(ctx context.Context, uuid string) (int64, error) {
	deleted, err := config.RedisClient.Del(ctx, uuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
