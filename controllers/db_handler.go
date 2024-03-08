package controllers

import (
	// "fmt"
	// "os"
	"log"

	// IMPORT DATABASE SQL
	"database/sql"

	// DRIVER GORM
	// "gorm.io/driver/mysql"
	// "gorm.io/gorm"
)

// SQL DB HANDLER
func connect() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/db_uts_pbp")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// GORM DB HANDLER
// func connectGorm() *gorm.DB {
// 	dbHost := os.Getenv("DB_HOST")
// 	fmt.Println(dbHost)
// 	db, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/db_uts_pbp"), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return db
// }
