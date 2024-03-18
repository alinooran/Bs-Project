package database

import (
	"fmt"
	"os"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var singleton *gorm.DB
var once sync.Once

func GetConn() *gorm.DB {
	once.Do(func() {
		HOST := os.Getenv("HOST")
		USER := os.Getenv("USER")
		PASSWORD := os.Getenv("PASSWORD")
		DBNAME := os.Getenv("DBNAME")
		PORT := os.Getenv("PORT")
		
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", USER, PASSWORD, HOST, PORT, DBNAME)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("cannot connect to database")
		}
		singleton = db
	})
	return singleton
}
