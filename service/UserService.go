package service

import (
	"TikTok/config"
	"TikTok/dao"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"strconv"
	"time"
)

type User struct {
	Id            int64
	Name          string
	Password      string
}

func (user User) GetUserList(users *[]User) bool{
	if err := dao.Db.Find(&user).Error; err != nil {
		log.Panicln("err:", err.Error())
		return false
	}
	return true
}

func (user *User) GetUserByUsername(name string) bool {
	dao.Db.Where("name = ?", name).First(&user)
	if user.Name == name {
		return true
	} else {
		return false
	}
}

func (user *User) GetUserById(id int64) bool {
	dao.Db.Where("id = ?", id).First(&user)
	if user.Id == id {
		return true
	} else {
		return false
	}
}


func (user *User) InsertUser() bool {
	if err := dao.Db.Create(&user).Error; err != nil {
		log.Panicln("err:", err.Error())
		return false
	}
	return true
}

func GenerateToken(username string) string {
	u := new(User)
	u.GetUserByUsername(username)
	token := NewToken(u)
	println(token)
	return token
}

func NewToken(u *User) string {
	expiresTime := time.Now().Unix() + int64(config.OneDayOfHours)
	fmt.Printf("%v\n", expiresTime)
	claims := jwt.StandardClaims{
		Audience:  u.Name,
		ExpiresAt: expiresTime,
		Id: 	   strconv.FormatInt(1, 10),
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