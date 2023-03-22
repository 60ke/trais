package downloader

import (
	"fmt"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Test_getLatestBlockFromDb(t *testing.T) {
	dsn := "root:Liu529966@tcp(127.0.0.1:3306)/bsc_explorer?charset=utf8mb4&parseTime=True&loc=Local"
	BscDB, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	fmt.Println(getLatestBlockFromDb(BscDB, "block_bak"))
}
