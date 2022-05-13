package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var Db *gorm.DB

func init() {
	var err error
	Db, err = gorm.Open("mysql", "root:qazxc32@tcp(106.14.75.229:3306)/douyin")
	if err != nil {
		log.Panicln("err:", err.Error())
	}
}