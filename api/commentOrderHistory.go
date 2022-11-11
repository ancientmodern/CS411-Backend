package api

import (
	"context"
	. "example/web-service-gin/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func UpdateComment(c *gin.Context) {

	var req updateCommentRequest
	if err := c.BindJSON(&req); err != nil {
		fmt.Printf("BindJSON failed, err: %v\n", err)
		c.String(http.StatusBadRequest, "BindJSON failed, err: %v\n", err)
		return
	}
	orderID := req.OrderID
	rating := req.Rating
	content := req.Content

	var insertBool bool
	insertBool = false
	var params []interface{}
	params = append(params, orderID, rating, content)

	sqlStr := fmt.Sprintf("SELECT Rating, Content FROM Comments WHERE OrderID = %d", orderID)
	fmt.Println(sqlStr)

	rows, err := DBPool.Query(sqlStr)
	if err != nil {
		fmt.Printf("query failed, err: %v\n", err)
		c.String(http.StatusBadRequest, "query failed, err: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var row getCommentResponse
		err := rows.Scan(&row.Rating, &row.Content)
		if err != nil {
			insertBool = true
		}
	}

	if insertBool {
		addsqlStr := "INSERT INTO Comments (OrderID, Rating, Content) VALUES (?, ?, ?)"
		fmt.Println(addsqlStr)

		ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()
		stmt, err := DBPool.PrepareContext(ctx, addsqlStr)
		if err != nil {
			fmt.Printf("PrepareContext failed, err: %v\n", err)
			c.IndentedJSON(http.StatusBadRequest, updateCommentResponse{
				false,
				fmt.Sprintf("PrepareContext failed, err: %v\n", err),
			})
			return
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, params...)
		if err != nil {
			fmt.Printf("ExecContext failed, err: %v\n", err)
			c.IndentedJSON(http.StatusBadRequest, updateCommentResponse{
				false,
				fmt.Sprintf("ExecContext failed, err: %v\n", err),
			})
			return
		}
	} else {
		updatesqlStr := fmt.Sprintf(
			"UPDATE Comments "+
				"SET Rating = %d, Content = %s "+
				"WHERE orderID = ?",
			rating, content,
		)
		fmt.Println(updatesqlStr)
		ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFunc()
		stmt, err := DBPool.PrepareContext(ctx, updatesqlStr)
		if err != nil {
			fmt.Printf("PrepareContext failed, err: %v\n", err)
			c.IndentedJSON(http.StatusBadRequest, updateCommentResponse{
				false,
				fmt.Sprintf("PrepareContext failed, err: %v\n", err),
			})
			return
		}
		defer stmt.Close()
		_, err = stmt.ExecContext(ctx, orderID)
		if err != nil {
			fmt.Printf("ExecContext failed, err: %v\n", err)
			c.IndentedJSON(http.StatusBadRequest, updateCommentResponse{
				false,
				fmt.Sprintf("ExecContext failed, err: %v\n", err),
			})
			return
		}
	}
	c.IndentedJSON(http.StatusOK, updateCommentResponse{true, ""})
}

func GetComment(c *gin.Context) {
	orderID := c.DefaultQuery("orderID", "")
	sqlStr := fmt.Sprintf("SELECT Rating, Content FROM Comments WHERE OrderID = %s", orderID)
	fmt.Println(sqlStr)

	rows, err := DBPool.Query(sqlStr)
	if err != nil {
		fmt.Printf("query failed, err: %v\n", err)
		c.String(http.StatusBadRequest, "query failed, err: %v\n", err)
		return
	}
	defer rows.Close()

	var res getCommentResponse
	for rows.Next() {
		err := rows.Scan(&res.Rating, &res.Content)
		if err != nil {

		}
	}
	c.IndentedJSON(http.StatusOK, res)
}

func DeleteComment(c *gin.Context) {
	oid := c.Query("orderID")
	if oid == "" {
		fmt.Println("Missing query string 'orderID'")
		c.String(http.StatusBadRequest, "Missing query string 'orderID'")
	}

	orderID, err := strconv.Atoi(oid)
	if err != nil {
		fmt.Println("orderID has invalid format")
		c.IndentedJSON(http.StatusBadRequest, deleteCommentResponse{
			false,
			fmt.Sprintf("orderID has invalid format"),
		})
		return
	}
	sqlStr := "DELETE FROM Comments WHERE OrderID = ?"
	fmt.Println("DELETE FROM Comments WHERE OrderID = #{orderID}")

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	stmt, err := DBPool.PrepareContext(ctx, sqlStr)
	if err != nil {
		fmt.Printf("PrepareContext failed, err: %v\n", err)
		c.IndentedJSON(http.StatusBadRequest, deleteCommentResponse{
			false,
			fmt.Sprintf("PrepareContext failed, err: %v\n", err),
		})
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, orderID)
	if err != nil {
		fmt.Printf("ExecContext failed, err: %v\n", err)
		c.IndentedJSON(http.StatusBadRequest, deleteCommentResponse{
			false,
			fmt.Sprintf("ExecContext failed, err: %v\n", err),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, deleteCommentResponse{true, ""})
}

func SearchOrderHistory(c *gin.Context) {
	var req historyOrderRequest
	if c.Bind(&req) != nil {
		c.String(http.StatusBadRequest, "Query parameters are not correct")
		return
	}

	userID := c.DefaultQuery("userID", "")
	minPrice := c.DefaultQuery("minPrice", "0")
	maxPrice := c.DefaultQuery("maxPrice", "100000")
	minTime := c.DefaultQuery("minTime", "20000000000000")
	maxTime := c.DefaultQuery("maxTime", "20221111235959")
	order := c.DefaultQuery("orderBy", "orderID")
	ascend := c.DefaultQuery("ascend", "ASC")

	if ascend == "false" {
		ascend = "DESC"
	}

	sqlStr := fmt.Sprintf(
		"SELECT OrderID, OrderTime, DishPrice, DishID, RiderID "+
			"FROM Orders "+
			"WHERE UserID LIKE '%s' AND OrderTime >= %s AND OrderTime <= %s "+
			"AND DishPrice >= %s AND DishPrice <= %s "+
			"ORDER BY %s %s",
		userID, minTime, maxTime, minPrice, maxPrice, order, ascend,
	)
	fmt.Println(sqlStr)

	rows, err := DBPool.Query(sqlStr)
	if err != nil {
		fmt.Printf("query failed, err: %v\n", err)
		c.String(http.StatusBadRequest, "query failed, err: %v\n", err)
		return
	}
	defer rows.Close()

	res := make(historyOrderResponse, 0)
	for rows.Next() {
		var row historyOrderResponseItem
		err := rows.Scan(&row.OrderID, &row.OrderTime, &row.DishPrice, &row.DishID, &row.RiderID)
		if err != nil {
			fmt.Printf("scan failed, err: %v\n", err)
			c.String(http.StatusBadRequest, "scan failed, err: %v\n", err)
			return
		}
		res = append(res, row)
	}

	c.IndentedJSON(http.StatusOK, res)
}
