package main

import (
	. "example/web-service-gin/api"
	. "example/web-service-gin/database"
	"fmt"
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

	v1 := router.Group("/api/v1")
	{
		v1.GET("/searchRestaurant", SearchRestaurant)
		v1.GET("/searchDish", SearchDish)
		v1.POST("/placeOrder", PlaceOrder)
		v1.DELETE("/deleteOrder", DeleteOrder)
	}

	router.Run("localhost:8080")
}
