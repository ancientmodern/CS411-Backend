package api

import (
	. "example/web-service-gin/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SearchRestaurant(c *gin.Context) {
	//var req SearchRestaurantRequest
	//if c.Bind(&req) != nil {
	//	c.String(http.StatusBadRequest, "Query parameters are not correct")
	//	return
	//}

	//var whereStr string
	//if req.RestaurantName != "" && req.ZipCode != 0 {
	//	whereStr = fmt.Sprintf("WHERE RestaurantName LIKE '%%%s%%' AND ZipCode = %d", req.RestaurantName, req.ZipCode)
	//} else if req.RestaurantName != "" {
	//	whereStr = fmt.Sprintf("WHERE RestaurantName LIKE '%%%s%%'", req.RestaurantName)
	//} else if req.ZipCode != 0 {
	//	whereStr = fmt.Sprintf("WHERE ZipCode = %d", req.ZipCode)
	//} else {
	//	whereStr = ""
	//}

	name := c.DefaultQuery("restaurantName", "")
	minCode := c.DefaultQuery("zipCode", "0")
	maxCode := c.DefaultQuery("zipCode", "100000")
	order := c.DefaultQuery("orderBy", "RestaurantID")
	ascend := c.DefaultQuery("ascend", "ASC")
	if ascend == "false" {
		ascend = "DESC"
	}

	sqlStr := fmt.Sprintf(
		"SELECT RestaurantID, RestaurantName, ZipCode, RestaurantAddr "+
			"FROM Restaurants "+
			"WHERE RestaurantName LIKE '%%%s%%' AND ZipCode >= %s AND ZipCode <= %s "+
			"ORDER BY %s %s",
		name, minCode, maxCode, order, ascend,
	)
	fmt.Println(sqlStr)

	rows, err := DBPool.Query(sqlStr)
	if err != nil {
		fmt.Printf("query failed, err: %v\n", err)
		c.String(http.StatusBadRequest, "query failed, err: %v\n", err)
		return
	}
	defer rows.Close()

	res := make(searchRestaurantResponse, 0)
	for rows.Next() {
		var row searchRestaurantResponseItem
		err := rows.Scan(&row.RestaurantID, &row.RestaurantName, &row.ZipCode, &row.RestaurantAddr)
		if err != nil {
			fmt.Printf("scan failed, err: %v\n", err)
			c.String(http.StatusBadRequest, "scan failed, err: %v\n", err)
			return
		}
		res = append(res, row)
	}

	c.IndentedJSON(http.StatusOK, res)
}

func SearchDish(c *gin.Context) {
	rid := c.Query("restaurantID")
	if rid == "" {
		fmt.Println("Missing query string 'restaurantID'")
		c.String(http.StatusBadRequest, "Missing query string 'restaurantID'")
		return
	}
	order := c.DefaultQuery("orderBy", "DishID")
	ascend := c.DefaultQuery("ascend", "ASC")
	if ascend == "false" {
		ascend = "DESC"
	}

	sqlStr := fmt.Sprintf(
		"SELECT DishID, DishName, Price, Category "+
			"FROM Dishes "+
			"WHERE RestaurantID = %s "+
			"ORDER BY %s %s",
		rid, order, ascend,
	)
	fmt.Println(sqlStr)

	rows, err := DBPool.Query(sqlStr)
	if err != nil {
		fmt.Printf("Query failed, err: %v\n", err)
		c.String(http.StatusBadRequest, "Query failed, err: %v\n", err)
		return
	}
	defer rows.Close()

	res := make(searchDishResponse, 0)
	for rows.Next() {
		var row searchDishResponseItem
		err := rows.Scan(&row.DishID, &row.DishName, &row.Price, &row.Category)
		if err != nil {
			fmt.Printf("Scan failed, err: %v\n", err)
			c.String(http.StatusBadRequest, "Scan failed, err: %v\n", err)
			return
		}
		res = append(res, row)
	}

	c.IndentedJSON(http.StatusOK, res)
}
