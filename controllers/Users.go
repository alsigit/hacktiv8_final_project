package controllers

import (
	"encoding/json"
	"fmt"
	"hacktiv8_final_project/helpers"
	"hacktiv8_final_project/middlewares"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cast"
)

type UserPayload struct {
	Age      int    `json:"age"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func Users_Register(c *gin.Context, db *sqlx.DB) {
	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	var postdata UserPayload
	json.Unmarshal([]byte(jsonData), &postdata)
	skrg := time.Now()

	valid_age, message_age := helpers.User_Age(postdata.Age)
	if !valid_age {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": message_age,
			"data":    nil,
		})
		return
	}

	valid_email, message_email := helpers.User_Email(postdata.Email)
	if !valid_email {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": message_email,
			"data":    nil,
		})
		return
	}

	valid_username, message_username := helpers.User_Username(postdata.Username)
	if !valid_username {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": message_username,
			"data":    nil,
		})
		return
	}

	valid_password, message_password := helpers.User_Password(postdata.Password)
	if !valid_password {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": message_password,
			"data":    nil,
		})
		return
	}

	get_UserData := helpers.DatabaseQuerySingleRow(db, `SELECT username, email FROM public."User" WHERE username = $1 OR email = $2`, postdata.Username, postdata.Email)

	if cast.ToString(get_UserData["username"]) == postdata.Username {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Username already exist.",
			"data":    nil,
		})
		return
	}

	if cast.ToString(get_UserData["email"]) == postdata.Email {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Email already exist.",
			"data":    nil,
		})
		return
	}

	_, enc_pass := helpers.BcryptEncode(postdata.Password)
	lastInsertID := 0
	err := db.QueryRow(`INSERT INTO public."User"(username,email,password,age,created_at) VALUES ($1,$2,$3,$4,$5) RETURNING id`, postdata.Username, postdata.Email, enc_pass, postdata.Age, skrg.Format("2006-01-02")).Scan(&lastInsertID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error",
			"data":    nil,
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"age":      postdata.Age,
			"email":    postdata.Email,
			"id":       lastInsertID,
			"username": postdata.Username,
		})
	}
}

func User_update(c *gin.Context, db *sqlx.DB) {
	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	var postdata UserPayload
	json.Unmarshal([]byte(jsonData), &postdata)
	user_id := c.Param("userid")
	skrg := time.Now()

	if _, err := strconv.Atoi(user_id); err == nil {
		// fmt.Printf("%q looks like a number.\n", user_id)
	} else {
		claim, _ := middlewares.GetJWTClaims(c.Request.Header["Authorization"], "id")
		user_id = claim
	}

	valid_email, message_email := helpers.User_Email(postdata.Email)
	if !valid_email {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": message_email,
			"data":    nil,
		})
		return
	}

	valid_username, message_username := helpers.User_Username(postdata.Username)
	if !valid_username {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": message_username,
			"data":    nil,
		})
		return
	}

	currentEmail, _ := middlewares.GetJWTClaims(c.Request.Header["Authorization"], "email")
	currentUsername, _ := middlewares.GetJWTClaims(c.Request.Header["Authorization"], "username")
	currentAge, _ := middlewares.GetJWTClaims(c.Request.Header["Authorization"], "age")

	if postdata.Username != currentUsername {
		check_username := helpers.DatabaseQuerySingleRow(db, `SELECT username, email FROM public."User" WHERE username = $1`, postdata.Username)
		if len(check_username) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Username already exist.",
				"data":    nil,
			})
			return
		}
	}

	if postdata.Email != currentEmail {
		check_email := helpers.DatabaseQuerySingleRow(db, `SELECT username, email FROM public."User" WHERE email = $1`, postdata.Email)
		if len(check_email) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Email already exist.",
				"data":    nil,
			})
			return
		}
	}

	_, err := db.Exec(`UPDATE public."User" SET email=$1, username=$2, updated_at=$3 WHERE id=$4`, postdata.Email, postdata.Username, skrg.Format("2006-01-02"), user_id)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error.",
			"data":    nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"id":         user_id,
			"email":      postdata.Email,
			"username":   postdata.Username,
			"age":        cast.ToInt(currentAge),
			"updated_at": skrg.Format("2006-01-02"),
		})
	}
}

func User_Delete(c *gin.Context, db *sqlx.DB) {
	// user_id, _ := middlewares.GetJWTClaims(c.Request.Header["Authorization"], "id")
	user_id := c.Param("userid")
	_, err := db.Exec(`DELETE FROM public."User" WHERE id = $1`, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error.",
			"data":    nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Your account has been successfully deleted",
		})
	}
}
