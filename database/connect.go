package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

var DBPool *sql.DB

// InitDB initializes a TCP connection pool for a Cloud SQL
// instance of MySQL.
func InitDB() error {
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
		dbUser    = mustGetenv("DB_USER")       // e.g. 'my-DBPool-user'
		dbPwd     = mustGetenv("DB_PASS")       // e.g. 'my-DBPool-password'
		dbName    = mustGetenv("DB_NAME")       // e.g. 'my-database'
		dbPort    = mustGetenv("DB_PORT")       // e.g. '3306'
		dbTCPHost = mustGetenv("INSTANCE_HOST") // e.g. '127.0.0.1' ('172.17.0.1' if deployed to GAE Flex)
	)

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPwd, dbTCPHost, dbPort, dbName)

	fmt.Println(dbURI)

	// DBPool is the pool of database connections.
	var err error
	DBPool, err = sql.Open("mysql", dbURI)
	if err != nil {
		return fmt.Errorf("sql.Open: %v", err)
	}

	// Test connection
	//if err = DBPool.Ping(); err != nil {
	//	fmt.Printf("init DBPool failed, err:%v\n", err)
	//	return
	//} else {
	//	fmt.Println("connection success")
	//}

	// config
	DBPool.SetConnMaxLifetime(time.Minute * 3)
	DBPool.SetMaxOpenConns(10)
	DBPool.SetMaxIdleConns(10)

	return nil
}
