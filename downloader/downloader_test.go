package downloader

import (
	"fmt"
	"testing"

	"github.com/60ke/trais/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Test_getLatestBlockFromDb(t *testing.T) {
	dsn := "root:Liu529966@tcp(127.0.0.1:3306)/bsc_explorer?charset=utf8mb4&parseTime=True&loc=Local"
	BscDB, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	fmt.Println(getLatestBlockFromDb(BscDB, "block"))
	var nums []int64
	BscDB.Table("block").Order("number desc").Limit(1000).Select("Number").Find(&nums)
	maxNum, minNum := nums[len(nums)-1], nums[0]
	for i := minNum; i < maxNum; i++ {
		if !SliceContain(nums, i) {
			fmt.Println(i)
		}
	}
}

func SliceContain[Item comparable](items []Item, item Item) bool {
	for i := 0; i < len(items); i++ {
		if items[i] == item {
			return true
		}
	}
	return false
}

func Test_insertBscBlock(t *testing.T) {
	dsn := "root:Liu529966@tcp(127.0.0.1:3306)/bsc_explorer?charset=utf8mb4&parseTime=True&loc=Local"
	BscDB, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	var block db.BscBlockTable
	var blocks []db.BscBlockTable
	block.Number = 313071
	block.CreditValue = "1"
	block.Hash = "xxx"

	BscDB.Table(block.TableName()).Where("number = ?", block.Number).FirstOrCreate(&block)
	fmt.Println(len(blocks))

	BscDB.Table(block.TableName()).FirstOrCreate(&block)

}
