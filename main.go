package main

import (
	controller "./controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3005", "http://localhost:3000"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
	/*	ExposeHeaders:    []string{"Content-Length"},*/
	}))
	router.Use(cors.New(config))
	client := router.Group("/api")
	{
		client.POST("/clothes/create", controller.CreateClothes)
		client.GET("/clothesList", controller.GetList2)
		client.GET("/category", controller.GetAllCategory)
		client.DELETE("/clothes/delete/:id", controller.DeleteClothes)
		client.POST("/clothes/uploadImage", controller.UploadImage)
		client.GET("/clothes/image/:id", controller.GetImageByClothesId)
		client.DELETE("/clothes/image/delete/:id", controller.DeleteClothes)
		client.GET("/clothes/countPage", controller.CountPageClothes)
		client.POST("/login", controller.Login)
		client.GET("/logout", controller.Logout)
		client.GET("/getSession", controller.GetSession)
	}

	return router
}

func main() {
	router := setupRouter()
	router.Run(":8080") // Ứng dụng chạy tại cổng 8080
}