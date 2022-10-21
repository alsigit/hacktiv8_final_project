package controllers

import (
	"encoding/json"
	"hacktiv8_final_project/helpers"
	"hacktiv8_final_project/middlewares"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type PhotosPayload struct {
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhtotURL string `json:"photo_url"`
}

func Photos_Create(c *gin.Context, db *sqlx.DB) {
	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	var postdata PhotosPayload
	json.Unmarshal([]byte(jsonData), &postdata)

	user_id, _ := middlewares.GetJWTClaims(c.Request.Header["Authorization"], "id")
	skrg := time.Now()

	valid_title, message_title := helpers.Photos_Title(postdata.Title)
	if !valid_title {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": message_title,
			"data":    nil,
		})
		return
	}

	valid_url, message_url := helpers.Photos_Url(postdata.PhtotURL)
	if !valid_url {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": message_url,
			"data":    nil,
		})
		return
	}

	newfilename := ""

	if postdata.PhtotURL != "" {
		asset_image := postdata.PhtotURL
		split_str := strings.Split(asset_image, ",")
		ftype := ""
		filename := ""
		coI := strings.Index(string(asset_image), ",")
		fileType := strings.TrimSuffix(asset_image[5:coI], ";base64")
		if fileType == "image/png" {
			ftype = "png"
		} else if fileType == "image/jpeg" {
			ftype = "jpeg"
		} else if fileType == "image/jpg" {
			ftype = "jpg"
		}
		skrg := time.Now().UTC()
		filename = user_id + "-" + skrg.Format("20060102150405") + "." + ftype
		image_path := viper.GetString("asset_image_upload_path")

		save_file := helpers.B64Decode(split_str[1], image_path+filename)
		if save_file != nil {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": save_file.Error(),
				"data":    nil,
			})
			return
		}
		newfilename = filename
	}

	lastInsertID := 0
	err := db.QueryRow(`INSERT INTO public."Photo"(title,caption,photo_url,user_id,created_at) VALUES ($1,$2,$3,$4,$5) RETURNING id`, postdata.Title, postdata.Caption, newfilename, user_id, skrg.Format("2006-01-02")).Scan(&lastInsertID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error.",
			"data":    nil,
		})
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"id":        lastInsertID,
			"title":     postdata.Title,
			"caption":   postdata.Caption,
			"photo_url": newfilename,
			"user_id":   user_id,
			"create_at": skrg.Format("2006-01-02"),
		})
		return
	}
}

func Photos_List(c *gin.Context, db *sqlx.DB) {
	list := helpers.DatabaseQueryRows(db, `
		SELECT
			a.*,
			b.email,
			b.username
		FROM
			(
				SELECT
					*
				FROM
					public."Photo"
			) a
		LEFT OUTER JOIN public."User" b ON a.user_id=b.id
	`)

	resp := []map[string]interface{}{}
	format := "2006-01-02 15:04:05"
	for _, data := range list {
		cd, _ := time.Parse(format, cast.ToString(data["created_at"])[0:19])
		ud := "-"
		if cast.ToString(data["updated_at"]) != "" {
			ddd, _ := time.Parse(format, cast.ToString(data["updated_at"])[0:19])
			ud = ddd.Format("01 Jan 2006")
		}
		// fmt.Println(err)
		holder := data
		holder["User"] = map[string]interface{}{
			"email":    holder["email"],
			"username": holder["username"],
		}
		holder["created_at"] = cd.Format("01 Jan 2006")
		holder["updated_at"] = ud
		holder["photo_url"] = viper.GetString("base_url") + "user_photo/" + cast.ToString(data["photo_url"])
		delete(holder, "email")
		delete(holder, "username")

		resp = append(resp, holder)
	}
	res, _ := json.Marshal(resp)

	c.Data(http.StatusOK, "json", []byte(res))
}

func Photos_Update(c *gin.Context, db *sqlx.DB) {
	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	var postdata PhotosPayload
	json.Unmarshal([]byte(jsonData), &postdata)

	user_id, _ := middlewares.GetJWTClaims(c.Request.Header["Authorization"], "id")
	photo_id := c.Param("photoid")
	skrg := time.Now()

	valid_title, message_title := helpers.Photos_Title(postdata.Title)
	if !valid_title {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": message_title,
			"data":    nil,
		})
		return
	}

	valid_url, message_url := helpers.Photos_Url(postdata.PhtotURL)
	if !valid_url {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": message_url,
			"data":    nil,
		})
		return
	}

	get_photo := helpers.DatabaseQuerySingleRow(db, `SELECT * FROM public."Photo" WHERE id = $1`, photo_id)

	if cast.ToInt(get_photo["user_id"]) != cast.ToInt(user_id) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Only your photo(s) are editable.",
			"data":    nil,
		})
		return
	}

	newfilename := ""

	if postdata.PhtotURL != "" {
		asset_image := postdata.PhtotURL
		split_str := strings.Split(asset_image, ",")
		ftype := ""
		filename := ""
		coI := strings.Index(string(asset_image), ",")
		fileType := strings.TrimSuffix(asset_image[5:coI], ";base64")
		if fileType == "image/png" {
			ftype = "png"
		} else if fileType == "image/jpeg" {
			ftype = "jpeg"
		} else if fileType == "image/jpg" {
			ftype = "jpg"
		}
		skrg := time.Now().UTC()
		filename = user_id + "-" + skrg.Format("20060102150405") + "." + ftype
		image_path := viper.GetString("asset_image_upload_path")

		save_file := helpers.B64Decode(split_str[1], image_path+"/"+filename)
		if save_file != nil {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": save_file.Error(),
				"data":    nil,
			})
			return
		}
		newfilename = filename

		// delete old photo
		e := os.Remove(image_path + "/" + cast.ToString(get_photo["photo_url"]))
		if e != nil {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": save_file.Error(),
				"data":    nil,
			})
			return
		}
	}

	_, err := db.Exec(`UPDATE public."Photo" SET title=$1, caption=$2, photo_url=$3, updated_at=$4 WHERE id=$5`, postdata.Title, postdata.Caption, newfilename, skrg.Format("2006-01-02"), photo_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error.",
			"data":    nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"id":         cast.ToInt(photo_id),
			"title":      postdata.Title,
			"caption":    postdata.Caption,
			"photo_url":  newfilename,
			"user_id":    cast.ToInt(user_id),
			"create_at":  skrg.Format("2006-01-02"),
			"updated_at": skrg.Format("2006-01-02"),
		})
	}
}

func Photos_Delete(c *gin.Context, db *sqlx.DB) {
	photo_id := c.Param("photoid")
	user_id, _ := middlewares.GetJWTClaims(c.Request.Header["Authorization"], "id")
	get_photo := helpers.DatabaseQuerySingleRow(db, `SELECT * FROM public."Photo" WHERE id = $1`, photo_id)

	if cast.ToInt(get_photo["user_id"]) != cast.ToInt(user_id) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Only your photo(s) are deleteable.",
			"data":    nil,
		})
		return
	}
	_, err := db.Exec(`DELETE FROM public."Photo" WHERE id = $1`, photo_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error.",
			"data":    nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Your photo has been successfully deleted",
		})
	}
}
