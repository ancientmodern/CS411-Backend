package api

import (
	"context"
	. "example/web-service-gin/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func SearchRestaurant(c *gin.Context) {
	//var req SearchRestaurantRequest
	//if c.Bind(&req) != nil {
	//	c.String(http.StatusBadRequest, "Query parameters are not correct")
	//	return
	//}

	name := c.DefaultQuery("restaurantName", "")
	minCode := c.DefaultQuery("zipCode", "0")
	maxCode := c.DefaultQuery("zipCode", "100000")
	order := c.DefaultQuery("orderBy", "RestaurantID")
	ascend := c.DefaultQuery("ascend", "ASC")
	if ascend == "false" {
		ascend = "DESC"
	}

	// TODO: SQL Prepare
	// FIXME: Potential SQL Injection
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

	// TODO: SQL Prepare
	// FIXME: Potential SQL Injection
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

func getDishPrice(dishID int) (float64, error) {
	sqlStr := "SELECT Price FROM Dishes WHERE DishID = ?"
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	stmt, err := DBPool.PrepareContext(ctx, sqlStr)
	if err != nil {
		fmt.Printf("GetDishPrice query preparing failed, err: %v\n", err)
		return -1.0, err
	}
	defer stmt.Close()

	var price float64
	row := stmt.QueryRowContext(ctx, dishID)
	if err := row.Scan(&price); err != nil {
		fmt.Printf("GetDishPrice scanning failed, err: %v\n", err)
		return -1.0, err
	}
	return price, nil
}

func getRandomRiderID() (int, error) {
	sqlStr := "SELECT RiderID FROM Riders ORDER BY RAND() LIMIT 1"
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	stmt, err := DBPool.PrepareContext(ctx, sqlStr)
	if err != nil {
		fmt.Printf("GetRandomRiderID query preparing failed, err: %v\n", err)
		return -1, err
	}
	defer stmt.Close()

	var riderID int
	row := stmt.QueryRowContext(ctx)
	if err := row.Scan(&riderID); err != nil {
		fmt.Printf("GetRandomRiderID scanning failed, err: %v\n", err)
		return -1, err
	}
	return riderID, nil
}

func PlaceOrder(c *gin.Context) {
	var req placeOrderRequest

	if err := c.BindJSON(&req); err != nil {
		fmt.Printf("BindJSON failed, err: %v\n", err)
		c.String(http.StatusBadRequest, "BindJSON failed, err: %v\n", err)
		return
	}

	sqlStr := "INSERT INTO Orders(OrderTime, DishPrice, DishID, UserID, RiderID) VALUES "
	var inserts []string
	var params []interface{}
	res := make(placeOrderResponse, 0)

	for _, v := range req.DishIDList {
		inserts = append(inserts, "(?, ?, ?, ?, ?)")

		orderTime := time.Now().Format("20060201150405")
		dishPrice, err := getDishPrice(v)
		if err != nil {
			c.String(http.StatusBadRequest, "GetDishPrice failed, err: %v\n", err)
			return
		}
		riderID, err := getRandomRiderID()
		if err != nil {
			c.String(http.StatusBadRequest, "GetRandomRiderID failed, err: %v\n", err)
			return
		}
		params = append(params, orderTime, dishPrice, v, req.UserID, riderID)
		res = append(res, placeOrderResponseItem{0, riderID})
	}
	queryVals := strings.Join(inserts, ",")
	sqlStr += queryVals
	fmt.Println(sqlStr)
	fmt.Println(params)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := DBPool.PrepareContext(ctx, sqlStr)
	if err != nil {
		fmt.Printf("PrepareContext failed, err: %v\n", err)
		c.String(http.StatusBadRequest, "PrepareContext failed, err: %v\n", err)
		return
	}
	defer stmt.Close()
	queryRes, err := stmt.ExecContext(ctx, params...)
	if err != nil {
		fmt.Printf("ExecContext failed, err: %v\n", err)
		c.String(http.StatusBadRequest, "ExecContext failed, err: %v\n", err)
		return
	}
	firstID, err := queryRes.LastInsertId()
	if err != nil {
		fmt.Printf("LastInsertId failed, err: %v\n", err)
		c.String(http.StatusBadRequest, "LastInsertId failed, err: %v\n", err)
		return
	}

	for i := 0; i < len(res); i++ {
		res[i].OrderID = int(firstID) + i
	}
	c.IndentedJSON(http.StatusOK, res)
}

func DeleteOrder(c *gin.Context) {
	oid := c.Query("orderID")
	if oid == "" {
		fmt.Println("Missing query string 'orderID'")
		c.IndentedJSON(http.StatusBadRequest, deleteOrderResponse{
			-1,
			false,
			"Missing query string 'orderID'"})
		return
	}

	orderID, err := strconv.Atoi(oid)
	if err != nil {
		fmt.Println("orderID has invalid format")
		c.IndentedJSON(http.StatusBadRequest, deleteOrderResponse{
			-1,
			false,
			"orderID has invalid format"})
		return
	}

	sqlStr := "DELETE FROM Orders WHERE OrderID = ?"
	fmt.Printf("DELETE FROM Orders WHERE OrderID = %d", orderID)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	stmt, err := DBPool.PrepareContext(ctx, sqlStr)
	if err != nil {
		fmt.Printf("PrepareContext failed, err: %v\n", err)
		c.IndentedJSON(http.StatusBadRequest, deleteOrderResponse{
			-1,
			false,
			fmt.Sprintf("PrepareContext failed, err: %v\n", err)})
		return
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, orderID)
	if err != nil {
		fmt.Printf("ExecContext failed, err: %v\n", err)
		c.IndentedJSON(http.StatusBadRequest, deleteOrderResponse{
			-1,
			false,
			fmt.Sprintf("ExecContext failed, err: %v\n", err)})
		return
	}

	c.IndentedJSON(http.StatusOK, deleteOrderResponse{orderID, true, ""})
}

func AdvancedCustomers(c *gin.Context) {
	minDishPrice := c.DefaultQuery("minDishPrice", "0")
	minTime := c.DefaultQuery("minTime", "20000000000000")
	minOrders := c.DefaultQuery("minOrders", "0")

	// TODO: SQL Prepare
	// FIXME: Potential SQL Injection
	sqlStr := "SELECT UserID, UserName, COUNT(OrderID) as numberOfOrders, " +
		"FROM Users NATURAL JOIN Orders " +
		"WHERE DishPrice > ? AND OrderTime > ? " +
		"GROUP BY UserID " +
		"HAVING numberOfOrders > ? " +
		"ORDER BY numberOfOrders DESC;"
	fmt.Println(sqlStr)

	rows, err := DBPool.Query(sqlStr, minDishPrice, minTime, minOrders)
	if err != nil {
		fmt.Printf("query failed, err: %v\n", err)
		c.String(http.StatusBadRequest, "query failed, err: %v\n", err)
		return
	}
	defer rows.Close()

	res := make(advancedCustomersResponse, 0)
	for rows.Next() {
		var row advancedCustomersResponseItem
		err := rows.Scan(&row.UserID, &row.UserName, &row.NumberOfOrders)
		if err != nil {
			fmt.Printf("scan failed, err: %v\n", err)
			c.String(http.StatusBadRequest, "scan failed, err: %v\n", err)
			return
		}
		res = append(res, row)
	}

	c.IndentedJSON(http.StatusOK, res)
}
