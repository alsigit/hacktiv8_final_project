package main

import (
	"fmt"
	"hacktiv8_final_project/controllers"
	"hacktiv8_final_project/core"
	"hacktiv8_final_project/middlewares"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

var MySecret = []byte("")

func main() {
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	// viper.AddConfigPath("/root/go-hactiv8-api/")
	viper.SetConfigName("app.conf")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	db := core.ConnectDB()
	defer db.Close()

	MySecret = []byte(viper.GetString("jwt.token_secret"))

	router := gin.New()
	router.Use(cors.Default())

	router.Use(gin.Recovery())

	routing(router, db)

	tmphttpreadheadertimeout, _ := time.ParseDuration(viper.GetString("server.readheadertimeout") + "s")
	tmphttpreadtimeout, _ := time.ParseDuration(viper.GetString("server.readtimeout") + "s")
	tmphttpwritetimeout, _ := time.ParseDuration(viper.GetString("server.writetimeout") + "s")
	tmphttpidletimeout, _ := time.ParseDuration(viper.GetString("server.idletimeout") + "s")

	s := &http.Server{
		Addr:              ":" + viper.GetString("server.port"),
		Handler:           router,
		ReadHeaderTimeout: tmphttpreadheadertimeout,
		ReadTimeout:       tmphttpreadtimeout,
		WriteTimeout:      tmphttpwritetimeout,
		IdleTimeout:       tmphttpidletimeout,
	}

	fmt.Println("Server running on port:", viper.GetString("server.port"))
	s.ListenAndServe()
}

func routing(router *gin.Engine, db *sqlx.DB) {

	router.Static("/assets", "./views/statics")
	router.Static("/user_photo", viper.GetString("asset_image_upload_path"))
	router.LoadHTMLGlob("views/templates/*")
	// login
	router.GET("", func(ctx *gin.Context) { controllers.Login(ctx) })
	router.GET("/register", func(ctx *gin.Context) { controllers.Register(ctx) })
	router.GET("/photo", func(ctx *gin.Context) { controllers.View_PhotoList(ctx) })
	router.GET("/social_media", func(ctx *gin.Context) { controllers.SocialMedia(ctx) })
	router.GET("/coment", func(ctx *gin.Context) { controllers.Coment(ctx) })

	// users
	router.POST("users/register", func(ctx *gin.Context) { controllers.Users_Register(ctx, db) })
	router.POST("users/login", func(ctx *gin.Context) { controllers.Auth_Login(ctx, db) })
	router.PUT("users/:userid", func(ctx *gin.Context) { controllers.User_update(ctx, db) }).Use(func(ctx *gin.Context) { middlewares.VerifyJWT(ctx) })
	router.DELETE("users/:userid", func(ctx *gin.Context) { controllers.User_Delete(ctx, db) }).Use(func(ctx *gin.Context) { middlewares.VerifyJWT(ctx) })

	// photos
	router.POST("photos", func(ctx *gin.Context) { controllers.Photos_Create(ctx, db) }).Use(func(ctx *gin.Context) { middlewares.VerifyJWT(ctx) })
	router.GET("photos", func(ctx *gin.Context) { controllers.Photos_List(ctx, db) }).Use(func(ctx *gin.Context) { middlewares.VerifyJWT(ctx) })
	router.PUT("photos/:photoid", func(ctx *gin.Context) { controllers.Photos_Update(ctx, db) }).Use(func(ctx *gin.Context) { middlewares.VerifyJWT(ctx) })
	router.DELETE("photos/:photoid", func(ctx *gin.Context) { controllers.Photos_Delete(ctx, db) }).Use(func(ctx *gin.Context) { middlewares.VerifyJWT(ctx) })

	// comment
	router.POST("comments", func(ctx *gin.Context) { controllers.Comments_Create(ctx, db) }).Use(func(ctx *gin.Context) { middlewares.VerifyJWT(ctx) })
	router.GET("comments", func(ctx *gin.Context) { controllers.Comments_List(ctx, db) }).Use(func(ctx *gin.Context) { middlewares.VerifyJWT(ctx) })
	router.PUT("comments/:commentid", func(ctx *gin.Context) { controllers.Comments_Update(ctx, db) }).Use(func(ctx *gin.Context) { middlewares.VerifyJWT(ctx) })
	router.DELETE("comments/:commentid", func(ctx *gin.Context) { controllers.Comments_Delete(ctx, db) }).Use(func(ctx *gin.Context) { middlewares.VerifyJWT(ctx) })

	// social media
	router.POST("socialmedias", func(ctx *gin.Context) { controllers.Socmed_Create(ctx, db) }).Use(func(ctx *gin.Context) { middlewares.VerifyJWT(ctx) })
	router.GET("socialmedias", func(ctx *gin.Context) { controllers.Socmed_List(ctx, db) }).Use(func(ctx *gin.Context) { middlewares.VerifyJWT(ctx) })
	router.PUT("socialmedias/:socialMediaId", func(ctx *gin.Context) { controllers.Socmed_Update(ctx, db) }).Use(func(ctx *gin.Context) { middlewares.VerifyJWT(ctx) })
	router.DELETE("socialmedias/:socialMediaId", func(ctx *gin.Context) { controllers.Socmed_Delete(ctx, db) }).Use(func(ctx *gin.Context) { middlewares.VerifyJWT(ctx) })
}
