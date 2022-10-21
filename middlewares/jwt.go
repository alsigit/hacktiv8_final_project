package middlewares

import (
	"errors"
	"fmt"
	"hacktiv8_final_project/core"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type MyClaims struct {
	UserName string `json:"username"`
	ID       string `json:"id"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func GenToken(username, id, email, secret string, age, expired int) (string, error) {
	// Create our own statement
	JWTTokenSecret := []byte(secret)
	JWTExp := expired
	TokenExpireDuration := time.Now().Add(time.Second * time.Duration(JWTExp)).Unix()
	c := MyClaims{
		username, // Custom field
		id,
		age,
		email,
		jwt.StandardClaims{
			ExpiresAt: TokenExpireDuration, // Expiration time
			Issuer:    "BCK-SERVICE",       // Issuer
		},
	}
	// Creates a signed object using the specified signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// Use the specified secret signature and obtain the complete encoded string token
	return token.SignedString(JWTTokenSecret)
}

func VerifyJWT(c *gin.Context) {
	authorization := strings.Split(c.GetHeader("Authorization"), " ")
	isBearer := authorization[0]
	// if isBearer not equal with'bearer'
	if strings.ToLower(isBearer) != "bearer" {
		c.Abort()
		core.InvalidJWTRes(c)
		return
	}
	JWT := authorization[1]

	// check is jwt have valid token secret key
	token, err := core.VerifyTokenSecretKey(JWT)
	if err != nil {
		c.Abort()
		core.InvalidJWTRes(c)
		return
	}

	// if token exp already expired
	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		c.Abort()
		core.InvalidJWTRes(c)
		return
	}

	// if token exp still remains
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Next()
		return
	}
}

func GetJWTClaims(header []string, key string) (string, error) {
	if header[0] == "" {
		return "", errors.New("header empty, request ignored")
	}
	jwtsplit := strings.Split(header[0], " ")
	token, _, err := new(jwt.Parser).ParseUnverified(jwtsplit[1], jwt.MapClaims{})
	if err != nil {
		fmt.Printf("Error %s", err)
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// obtains claims
		ret := fmt.Sprint(claims[key])
		return ret, nil
	}

	return "", nil
}
