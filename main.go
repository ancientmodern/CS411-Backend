package main

import (
	. "example/web-service-gin/api"
	. "example/web-service-gin/database"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	err := InitDB()
	if err != nil {
		fmt.Printf("init db failed, err:%v\n", err)
		return
	}

	router := gin.Default()

	//router.Use(cors.New(cors.Config{
	//	// AllowOrigins:     []string{"https://foo.com"},
	//	AllowMethods:     []string{"GET", "POST", "DELETE"},
	//	AllowHeaders:     []string{"Origin"},
	//	ExposeHeaders:    []string{"Content-Length"},
	//	AllowCredentials: true,
	//	AllowOriginFunc: func(origin string) bool {
	//		return true
	//	},
	//	MaxAge: 12 * time.Hour,
	//}))

	router.Use(cors.Default())

	v1 := router.Group("/api/v1")
	{
		v1.GET("/searchRestaurant", SearchRestaurant)
		v1.GET("/searchDish", SearchDish)
		v1.POST("/placeOrder", PlaceOrder)
		v1.DELETE("/deleteOrder", DeleteOrder)
	}

	router.Run("localhost:80")
}
