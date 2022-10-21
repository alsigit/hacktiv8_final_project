package core

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InvalidJWTRes(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"data":    nil,
		"message": "invalid token",
	})
}

func InternalServerErrorRes(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"data":    nil,
		"message": err.Error(),
	})
}

func CreateJWTToken(claims jwt.MapClaims) (string, error) {
	JWTTokenSecret := viper.GetString("jwt.token_secret")
	JWTExp := viper.GetInt("jwt.expaired_duration")

	var err error
	// extend claims data
	// claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Second * time.Duration(JWTExp)).Unix()
	JWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := JWT.SignedString([]byte(JWTTokenSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func VerifyTokenSecretKey(JWTToken string) (*jwt.Token, error) {
	JWTTokenSecret := viper.GetString("jwt.token_secret")

	token, err := jwt.Parse(JWTToken, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWTTokenSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
