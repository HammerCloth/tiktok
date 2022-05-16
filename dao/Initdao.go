package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var Db *gorm.DB

func init() {
	var err error
	Db, err = gorm.Open("mysql", "douyin:zjqxy@tcp(43.138.25.60:3306)/douyin/")
	if err != nil {
		log.Panicln("err:", err.Error())
	}
}
