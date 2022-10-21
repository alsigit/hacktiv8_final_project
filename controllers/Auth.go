package controllers

import (
	"encoding/json"
	"fmt"
	"hacktiv8_final_project/helpers"
	"hacktiv8_final_project/middlewares"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type AuthLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"status":   "ok",
		"base_url": viper.GetString("base_url"),
	})
}

func Auth_Login(c *gin.Context, db *sqlx.DB) {
	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	var postdata AuthLogin
	json.Unmarshal([]byte(jsonData), &postdata)

	if postdata.Email == "" || postdata.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Username password cannot be empty.",
			"data":    nil,
		})
		return
	}

	get_userdata := helpers.DatabaseQuerySingleRow(db, `SELECT * FROM public."User" WHERE email = $1`, postdata.Email)

	if len(get_userdata) < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User data not found.",
			"data":    nil,
		})
		return
	}

	// _, enc_password := helpers.BcryptEncode(postdata.Password)
	passValid := bcrypt.CompareHashAndPassword([]byte(cast.ToString(get_userdata["password"])), []byte(postdata.Password))
	if passValid != nil {
		fmt.Println(passValid)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Wrong Password.",
			"data":    nil,
		})
		return
	}

	secret := viper.GetString("jwt.token_secret")
	expired := viper.GetInt("jwt.expired_duration")

	token, errCJT := middlewares.GenToken(cast.ToString(get_userdata["username"]), cast.ToString(get_userdata["id"]), cast.ToString(get_userdata["email"]), secret, cast.ToInt(get_userdata["age"]), expired)
	if errCJT != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occured (JWT) : " + errCJT.Error(),
			"status":  "error",
			"data":    nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}

func Register(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"status":   "ok",
		"base_url": viper.GetString("base_url"),
	})
}
