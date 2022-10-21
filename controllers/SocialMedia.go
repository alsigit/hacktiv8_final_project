package controllers

import (
	"encoding/json"
	"hacktiv8_final_project/helpers"
	"hacktiv8_final_project/middlewares"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cast"
)

type SocmedPayload struct {
	Name string `json:"name"`
	URL  string `json:"social_media_url"`
}

func Socmed_Create(c *gin.Context, db *sqlx.DB) {
	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	var postdata SocmedPayload
	json.Unmarshal([]byte(jsonData), &postdata)

	user_id, _ := middlewares.GetJWTClaims(c.Request.Header["Authorization"], "id")
	skrg := time.Now()

	valid_name, message_name := helpers.SocMed_Name(postdata.Name)
	if !valid_name {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": message_name,
			"data":    nil,
		})
		return
	}

	valid_url, message_url := helpers.SocMed_Name(postdata.URL)
	if !valid_url {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": message_url,
			"data":    nil,
		})
		return
	}

	lastInsertID := 0
	err := db.QueryRow(`INSERT INTO public."SocialMedia"(name,social_media_url,user_id,created_at) VALUES ($1,$2,$3,$4) RETURNING id`, postdata.Name, postdata.URL, user_id, skrg.Format("2006-01-02")).Scan(&lastInsertID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error.",
			"data":    nil,
		})
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"id":               lastInsertID,
			"name":             postdata.Name,
			"social_media_url": postdata.URL,
			"user_id":          cast.ToInt(user_id),
			"create_at":        skrg.Format("2006-01-02"),
		})
		return
	}
}

func Socmed_List(c *gin.Context, db *sqlx.DB) {
	list := helpers.DatabaseQueryRows(db, `
		SELECT
			a.*,
			b.id AS user_id,
			b.email,
			b.username
		FROM
			(
				SELECT
					*
				FROM
					public."SocialMedia"
			) a
		LEFT OUTER JOIN public."User" b ON a.user_id=b.id
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
		delete(holder, "email")
		delete(holder, "username")

		resp = append(resp, holder)
	}

	c.JSON(http.StatusOK, gin.H{
		"social_medias": resp,
	})
}

func Socmed_Update(c *gin.Context, db *sqlx.DB) {
	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	var postdata SocmedPayload
	json.Unmarshal([]byte(jsonData), &postdata)

	socmed_id := c.Param("socialMediaId")
	skrg := time.Now()

	valid_name, message_name := helpers.SocMed_Name(postdata.Name)
	if !valid_name {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": message_name,
			"data":    nil,
		})
		return
	}

	valid_url, message_url := helpers.SocMed_Name(postdata.URL)
	if !valid_url {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": message_url,
			"data":    nil,
		})
		return
	}

	get_SocMed := helpers.DatabaseQuerySingleRow(db, `SELECT * FROM public."SocialMedia" WHERE id = $1`, socmed_id)
	user_id, _ := middlewares.GetJWTClaims(c.Request.Header["Authorization"], "id")

	if cast.ToInt(get_SocMed["user_id"]) != cast.ToInt(user_id) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Only your social media(s) are editable.",
			"data":    nil,
		})
		return
	}

	_, err := db.Exec(`UPDATE public."SocialMedia" SET name=$1, social_media_url=$2, updated_at=$3 WHERE id=$4`, postdata.Name, postdata.URL, skrg.Format("2006-01-02"), socmed_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error.",
			"data":    nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"id":               cast.ToInt(socmed_id),
			"name":             postdata.Name,
			"social_media_url": postdata.URL,
			"user_id":          cast.ToInt(user_id),
			"updated_at":       skrg.Format("2006-01-02"),
		})
	}
}

func Socmed_Delete(c *gin.Context, db *sqlx.DB) {
	comment_id := c.Param("socialMediaId")
	get_socmed := helpers.DatabaseQuerySingleRow(db, `SELECT * FROM public."SocialMedia" WHERE id = $1`, comment_id)
	user_id, _ := middlewares.GetJWTClaims(c.Request.Header["Authorization"], "id")

	if cast.ToInt(get_socmed["user_id"]) != cast.ToInt(user_id) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Only your social media(s) are deleteable.",
			"data":    nil,
		})
		return
	}

	_, err := db.Exec(`DELETE FROM public."SocialMedia" WHERE id = $1`, comment_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error.",
			"data":    nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Your social media has been deleted",
		})
	}
}
