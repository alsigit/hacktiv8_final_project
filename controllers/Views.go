package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func View_PhotoList(c *gin.Context) {
	c.HTML(http.StatusOK, "Photos.html", gin.H{
		"status":   "ok",
		"base_url": viper.GetString("base_url"),
	})
}
func SocialMedia(c *gin.Context) {
	c.HTML(http.StatusOK, "socialMedia.html", gin.H{
		"status":   "ok",
		"base_url": viper.GetString("base_url"),
	})
}
func Coment(c *gin.Context) {
	c.HTML(http.StatusOK, "coment.html", gin.H{
		"status":   "ok",
		"base_url": viper.GetString("base_url"),
	})
}
