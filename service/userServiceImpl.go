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

func (tableUser TableUser) TableName() string {
	return "users"
}

func (usi *UserServiceImpl) GetTableUserList() []TableUser {
	tableUsers := []TableUser{}
	if err := dao.Db.Find(&tableUsers).Error; err != nil {
		log.Panicln("err:", err.Error())
		return tableUsers
	}
	return tableUsers
}

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

func (usi *UserServiceImpl) InsertTableUser(tableUser *TableUser) bool {
	if err := dao.Db.Create(&tableUser).Error; err != nil {
		log.Println("插入失败")
		return false
	}
	return true
}

func (usi *UserServiceImpl) GetUserById(id int64) (User, error) {
	tableUser := TableUser{}
	if err := dao.Db.Where("id = ?", id).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
	}
	if tableUser.Id == id {
		log.Println("查询成功")
	} else {
		log.Println("未找到该用户")
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
	return user, errors.New("query fail")
}

func (usi *UserServiceImpl) GetUserByIdWithCurId(id int64, curId int64) (User, error) {
	tableUser := TableUser{}
	if err := dao.Db.Where("id = ?", id).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
	}
	if tableUser.Id == id {
		log.Println("查询成功")
	} else {
		log.Println("未找到该用户")
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
	return user, errors.New("query fail")
}

func GenerateToken(username string) string {
	u := UserService.GetTableUserByUsername(new(UserServiceImpl), username)
	fmt.Printf("generatetoken: %v\n", u)
	token := NewToken(u)
	println(token)
	return token
}

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

func EnCoder(password string) string {
	h := hmac.New(sha256.New, []byte(password))
	sha := hex.EncodeToString(h.Sum(nil))
	fmt.Println("Result: " + sha)
	return sha
}
