package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var TMDB *gorm.DB

func TMDBInit(user, pass, host, dbName string) *gorm.DB {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, dbName)
	// dsn := "k:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	TMDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("db connect error")
	}
	fmt.Println("db connect success")
	return TMDB
}

func TMDBClose() {
	db, _ := TMDB.DB()
	db.Close()
}
