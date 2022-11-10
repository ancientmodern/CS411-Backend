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

func AddComment(c *gin.Context) {

	orderID := c.DefaultQuery("orderID", "")
	rating := c.DefaultQuery("rating", "0")
	content := c.DefaultQuery("content", "")
	var params []interface{}
	params = append(params, orderID, rating, content)

	sqlStr := "INSERT INTO Comments (OrderID, Rating, Content) VALUES (?, ?, ?)"

	fmt.Println(sqlStr)

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
	commentID, err := queryRes.LastInsertId()
	if err != nil {
		fmt.Printf("LastInsertId failed, err: %v\n", err)
		c.String(http.StatusBadRequest, "LastInsertId failed, err: %v\n", err)
		return
	}

	c.IndentedJSON(http.StatusOK, int(commentID))
}

func UpdateComment(c *gin.Context) {
	cid := c.Query("commentID")
	if cid == "" {
		fmt.Println("Missing query string 'commentID'")
		c.String(http.StatusBadRequest, "Missing query string 'commentID'")
	}
	commentID, err := strconv.Atoi(cid)
	if err != nil {
		fmt.Println("orderID has invalid format")
		c.IndentedJSON(http.StatusBadRequest, updateCommentResponse{
			false,
			"commentID invalid format (str to int)",
		})
		return
	}
	rating := c.DefaultQuery("rating", "0")
	content := c.DefaultQuery("content", "")
	sqlStr := fmt.Sprintf(
		"UPDATE Comments "+
			"SET Rating = %s, Content = %s "+
			"WHERE commentID = ?",
		rating, content,
	)
	fmt.Println(sqlStr)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	stmt, err := DBPool.PrepareContext(ctx, sqlStr)
	if err != nil {
		fmt.Printf("PrepareContext failed, err: %v\n", err)
		c.IndentedJSON(http.StatusBadRequest, updateCommentResponse{
			false,
			fmt.Sprintf("PrepareContext failed, err: %v\n", err),
		})
		return
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, commentID)
	if err != nil {
		fmt.Printf("ExecContext failed, err: %v\n", err)
		c.IndentedJSON(http.StatusBadRequest, updateCommentResponse{
			false,
			fmt.Sprintf("ExecContext failed, err: %v\n", err),
		})
		return
	}
	c.IndentedJSON(http.StatusOK, updateCommentResponse{true, ""})
}
func DeleteComment(c *gin.Context) {
	cid := c.Query("commentID")
	if cid == "" {
		fmt.Println("Missing query string 'commentID'")
		c.String(http.StatusBadRequest, "Missing query string 'commentID'")
	}
	commentID, err := strconv.Atoi(cid)
	if err != nil {
		fmt.Println("commentID has invalid format")
		c.IndentedJSON(http.StatusBadRequest, deleteCommentResponse{
			false,
			fmt.Sprintf("commentID has invalid format"),
		})
		return
	}
	sqlStr := "DELETE FROM Comments WHERE CommentID = ?"

	fmt.Println("DELETE FROM Comments WHERE CommentID = #{commentID}")
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
	_, err = stmt.ExecContext(ctx, commentID)
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
