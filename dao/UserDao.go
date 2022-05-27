package dao

import (
	"log"
)

// TableUser 对应数据库User表结构的结构体
type TableUser struct {
	Id       int64
	Name     string
	Password string
}

// TableName 修改表名映射
func (tableUser TableUser) TableName() string {
	return "users"
}

// GetTableUserList 获取全部TableUser对象
func GetTableUserList() ([]TableUser, error) {
	tableUsers := []TableUser{}
	if err := Db.Find(&tableUsers).Error; err != nil {
		log.Println(err.Error())
		return tableUsers, err
	}
	return tableUsers, nil
}

// GetTableUserByUsername 根据username获得TableUser对象
func GetTableUserByUsername(name string) (TableUser, error) {
	tableUser := TableUser{}
	if err := Db.Where("name = ?", name).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return tableUser, err
	}
	return tableUser, nil
}

// GetTableUserById 根据user_id获得TableUser对象
func GetTableUserById(id int64) (TableUser, error) {
	tableUser := TableUser{}
	if err := Db.Where("id = ?", id).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return tableUser, err
	}
	return tableUser, nil
}

// InsertTableUser 将tableUser插入表内
func InsertTableUser(tableUser *TableUser) bool {
	if err := Db.Create(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
