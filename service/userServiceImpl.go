package service

import (
	"TikTok/config"
	"TikTok/dao"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"strconv"
	"time"
)

type UserServiceImpl struct {
	FollowService
}

// TableName 修改表名映射
func (tableUser TableUser) TableName() string {
	return "users"
}

// GetTableUserList 获得全部TableUser对象
func (usi *UserServiceImpl) GetTableUserList() []TableUser {
	tableUsers := []TableUser{}
	if err := dao.Db.Find(&tableUsers).Error; err != nil {
		log.Panicln("err:", err.Error())
		return tableUsers
	}
	return tableUsers
}

// GetTableUserByUsername 根据username获得TableUser对象
func (usi *UserServiceImpl) GetTableUserByUsername(name string) TableUser {
	tableUser := TableUser{}
	if err := dao.Db.Where("name = ?", name).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
	}
	if tableUser.Name == name {
		log.Println("查询成功")
		return tableUser
	} else {
		log.Println("未找到该用户")
		return tableUser
	}
}

// GetTableUserById 根据user_id获得TableUser对象
func (usi *UserServiceImpl) GetTableUserById(id int64) TableUser {
	tableUser := TableUser{}
	if err := dao.Db.Where("id = ?", id).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
	}
	if tableUser.Id == id {
		log.Println("查询成功")
		return tableUser
	} else {
		log.Println("未找到该用户")
		return tableUser
	}
}

// InsertTableUser 将tableUser插入表内
func (usi *UserServiceImpl) InsertTableUser(tableUser *TableUser) bool {
	if err := dao.Db.Create(&tableUser).Error; err != nil {
		log.Println("插入失败")
		return false
	}
	return true
}

// GetUserById 未登录情况下,根据user_id获得User对象
func (usi *UserServiceImpl) GetUserById(id int64) (User, error) {
	tableUser := TableUser{}
	if err := dao.Db.Where("id = ?", id).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
	}
	if tableUser.Id != id {
		log.Println("未找到该用户")
		return User{}, errors.New("query fail")
	} else {
		log.Println("查询成功")
	}
	fsi := new(FollowServiceImpl)
	impl := UserServiceImpl{fsi}
	followCount, _ := impl.FollowService.GetFollowingCnt(id)
	followerCount, _ := impl.FollowService.GetFollowerCnt(id)
	user := User{
		Id:            id,
		Name:          tableUser.Name,
		FollowCount:   followCount,
		FollowerCount: followerCount,
		IsFollow:      false,
	}
	return user, nil
}

// GetUserByIdWithCurId 已登录(curID)情况下,根据user_id获得User对象
func (usi *UserServiceImpl) GetUserByIdWithCurId(id int64, curId int64) (User, error) {
	tableUser := TableUser{}
	if err := dao.Db.Where("id = ?", id).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
	}
	if tableUser.Id != id {
		log.Println("未找到该用户")
		return User{}, errors.New("query fail")
	} else {
		log.Println("查询成功")
	}
	fsi := new(FollowServiceImpl)
	impl := UserServiceImpl{fsi}
	followCount, _ := impl.FollowService.GetFollowingCnt(id)
	followerCount, _ := impl.FollowService.GetFollowerCnt(id)
	isfollow, _ := impl.FollowService.IsFollowing(curId, id)
	user := User{
		Id:            id,
		Name:          tableUser.Name,
		FollowCount:   followCount,
		FollowerCount: followerCount,
		IsFollow:      isfollow,
	}
	return user, nil
}

// GenerateToken 根据username生成一个token
func GenerateToken(username string) string {
	u := UserService.GetTableUserByUsername(new(UserServiceImpl), username)
	fmt.Printf("generatetoken: %v\n", u)
	token := NewToken(u)
	println(token)
	return token
}

// NewToken 根据信息创建token
func NewToken(u TableUser) string {
	expiresTime := time.Now().Unix() + int64(config.OneDayOfHours)
	fmt.Printf("%v\n", expiresTime)
	id64 := u.Id
	fmt.Printf("newtoken: %v\n", strconv.FormatInt(id64, 10))
	claims := jwt.StandardClaims{
		Audience:  u.Name,
		ExpiresAt: expiresTime,
		Id:        strconv.FormatInt(id64, 10),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "tiktok",
		NotBefore: time.Now().Unix(),
		Subject:   "token",
	}
	var jwtSecret = []byte(config.Secret)
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if token, err := tokenClaims.SignedString(jwtSecret); err == nil {
		token = "Bearer " + token
		println("generate token success!\n")
		return token
	} else {
		println("generate token fail\n")
		return "fail"
	}
}

// EnCoder 密码加密
func EnCoder(password string) string {
	h := hmac.New(sha256.New, []byte(password))
	sha := hex.EncodeToString(h.Sum(nil))
	fmt.Println("Result: " + sha)
	return sha
}
