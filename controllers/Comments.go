package controllers

import (
	"encoding/json"
	"fmt"
	"hacktiv8_final_project/helpers"
	"hacktiv8_final_project/middlewares"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cast"
)

type CommentsPayload struct {
	Message string `json:"message"`
	PhotoID int    `json:"photo_id"`
}

func Comments_Create(c *gin.Context, db *sqlx.DB) {
	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	var postdata CommentsPayload
	json.Unmarshal([]byte(jsonData), &postdata)

	user_id, _ := middlewares.GetJWTClaims(c.Request.Header["Authorization"], "id")
	skrg := time.Now()

	valid_msg, message_msg := helpers.Comments_Message(postdata.Message)
	if !valid_msg {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": message_msg,
			"data":    nil,
		})
		return
	}

	lastInsertID := 0
	err := db.QueryRow(`INSERT INTO public."Comment"(message,photo_id,user_id,created_at) VALUES ($1,$2,$3,$4) RETURNING id`, postdata.Message, postdata.PhotoID, user_id, skrg.Format("2006-01-02")).Scan(&lastInsertID)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error.",
			"data":    nil,
		})
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"id":        lastInsertID,
			"message":   postdata.Message,
			"photo_id":  postdata.PhotoID,
			"user_id":   cast.ToInt(user_id),
			"create_at": skrg.Format("2006-01-02"),
		})
		return
	}
}

func Comments_List(c *gin.Context, db *sqlx.DB) {
	list := helpers.DatabaseQueryRows(db, `
		SELECT
			a.*,
			b.id AS user_id,
			b.email,
			b.username,
			c.id AS photo_id,
			c.title,
			c.caption,
			c.photo_url,
			c.user_id AS comment_uid
		FROM
			(
				SELECT
					*
				FROM
					public."Comment"
			) a
		LEFT OUTER JOIN public."User" b ON a.user_id=b.id
		LEFT OUTER JOIN public."Photo" c ON a.photo_id=c.id
	`)

	resp := []map[string]interface{}{}

	for _, data := range list {
		holder := data
		holder["created_at"] = cast.ToString(data["created_at"])[0:10]
		if cast.ToString(data["updated_at"]) != "" {
			holder["updated_at"] = cast.ToString(data["updated_at"])[0:10]
		}
		holder["User"] = map[string]interface{}{
			"user_id":  holder["user_id"],
			"email":    holder["email"],
			"username": holder["username"],
		}
		holder["Photo"] = map[string]interface{}{
			"id":        holder["photo_id"],
			"title":     holder["title"],
			"caption":   holder["caption"],
			"photo_url": holder["photo_url"],
			"user_id":   holder["comment_uid"],
		}
		// delete(holder, "user_id")
		delete(holder, "email")
		delete(holder, "username")

		// delete(holder, "photo_id")
		delete(holder, "title")
		delete(holder, "caption")
		delete(holder, "photo_url")
		delete(holder, "comment_uid")

		resp = append(resp, holder)
	}
	res, _ := json.Marshal(resp)

	c.Data(http.StatusOK, "json", []byte(res))
}

func Comments_Update(c *gin.Context, db *sqlx.DB) {
	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	var postdata CommentsPayload
	json.Unmarshal([]byte(jsonData), &postdata)

	comment_id := c.Param("commentid")
	skrg := time.Now()

	valid_msg, message_msg := helpers.Comments_Message(postdata.Message)
	if !valid_msg {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": message_msg,
			"data":    nil,
		})
		return
	}

	get_comment := helpers.DatabaseQuerySingleRow(db, `SELECT * FROM public."Comment" WHERE id = $1`, comment_id)
	user_id, _ := middlewares.GetJWTClaims(c.Request.Header["Authorization"], "id")

	if cast.ToInt(get_comment["user_id"]) != cast.ToInt(user_id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Only your comments are editable.",
			"data":    nil,
		})
		return
	}

	_, err := db.Exec(`UPDATE public."Comment" SET message=$1, updated_at=$2 WHERE id=$3`, postdata.Message, skrg.Format("2006-01-02"), comment_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error.",
			"data":    nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"id":         cast.ToInt(comment_id),
			"message":    postdata.Message,
			"photo_id":   cast.ToInt(get_comment["photo_id"]),
			"user_id":    cast.ToInt(user_id),
			"updated_at": skrg.Format("2006-01-02"),
		})
	}
}

func Comments_Delete(c *gin.Context, db *sqlx.DB) {
	comment_id := c.Param("commentid")
	get_comment := helpers.DatabaseQuerySingleRow(db, `SELECT * FROM public."Comment" WHERE id = $1`, comment_id)
	user_id, _ := middlewares.GetJWTClaims(c.Request.Header["Authorization"], "id")

	if cast.ToInt(get_comment["user_id"]) != cast.ToInt(user_id) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Only your comments are deleteable.",
			"data":    nil,
		})
		return
	}

	_, err := db.Exec(`DELETE FROM public."Comment" WHERE id = $1`, comment_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error.",
			"data":    nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Your Comment has been successfully deleted",
		})
	}
}
