package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var Db *gorm.DB

func Init() {
	var err error
	//想要正确的处理time.Time,需要带上 parseTime 参数，
	//要支持完整的UTF-8编码，需要将 charset=utf8 更改为 charset=utf8mb4
	Db, err = gorm.Open("mysql", "douyin:zjqxy@tcp(43.138.25.60:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local")
	//打开gorm详细日志
	Db.LogMode(true)
	if err != nil {
		log.Panicln("err:", err.Error())
	}
}
