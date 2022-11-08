package main

import (
	"database/sql"
	. "example/web-service-gin/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	err := initDB()
	if err != nil {
		fmt.Printf("init db failed, err:%v\n", err)
		return
	}

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/searchRestaurant", searchRestaurant)
	}

	router.Run("localhost:8080")
}

// initDB initializes a TCP connection pool for a Cloud SQL
// instance of MySQL.
func initDB() error {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Warning: %s environment variable not set.", k)
		}
		return v
	}
	// Note: Saving credentials in environment variables is convenient, but not
	// secure - consider a more secure solution such as
	// Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
	// keep secrets safe.
	var (
		dbUser    = mustGetenv("DB_USER")       // e.g. 'my-db-user'
		dbPwd     = mustGetenv("DB_PASS")       // e.g. 'my-db-password'
		dbName    = mustGetenv("DB_NAME")       // e.g. 'my-database'
		dbPort    = mustGetenv("DB_PORT")       // e.g. '3306'
		dbTCPHost = mustGetenv("INSTANCE_HOST") // e.g. '127.0.0.1' ('172.17.0.1' if deployed to GAE Flex)
	)

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPwd, dbTCPHost, dbPort, dbName)

	fmt.Println(dbURI)

	// db is the pool of database connections.
	var err error
	db, err = sql.Open("mysql", dbURI)
	if err != nil {
		return fmt.Errorf("sql.Open: %v", err)
	}

	// Test connection
	//if err = db.Ping(); err != nil {
	//	fmt.Printf("init db failed, err:%v\n", err)
	//	return
	//} else {
	//	fmt.Println("connection success")
	//}

	// config
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return nil
}

func searchRestaurant(c *gin.Context) {
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

	rows, err := db.Query(sqlStr)
	if err != nil {
		fmt.Printf("query failed, err: %v\n", err)
		c.String(http.StatusBadRequest, "query failed, err: %v\n", err)
		return
	}
	defer rows.Close()

	var res SearchRestaurantResponse
	for rows.Next() {
		var row SearchRestaurantResponseItem
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
