package dao

import (
	"log"
	"sync"
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

type UserDao struct {
}

var (
	userDao  *UserDao
	userOnce sync.Once
)

func NewUserDaoInstance() *UserDao {
	userOnce.Do(
		func() {
			userDao = &UserDao{}
		})
	return userDao
}

func (*UserDao) GetTableUserList() ([]TableUser, error) {
	tableUsers := []TableUser{}
	Init()
	if err := Db.Find(&tableUsers).Error; err != nil {
		log.Println(err.Error())
		return tableUsers, err
	}
	return tableUsers, nil
}

// GetTableUserByUsername 根据username获得TableUser对象
func (*UserDao) GetTableUserByUsername(name string) (TableUser, error) {
	tableUser := TableUser{}
	Init()
	if err := Db.Where("name = ?", name).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return tableUser, err
	}
	return tableUser, nil
}

// GetTableUserById 根据user_id获得TableUser对象
func (*UserDao) GetTableUserById(id int64) (TableUser, error) {
	tableUser := TableUser{}
	Init()
	if err := Db.Where("id = ?", id).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return tableUser, err
	}
	return tableUser, nil
}

// InsertTableUser 将tableUser插入表内
func (*UserDao) InsertTableUser(tableUser *TableUser) bool {
	Init()
	if err := Db.Create(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
