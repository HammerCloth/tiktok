package service

import (
	"TikTok/config"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

type UserServiceImpl struct {
}

func (usi *UserServiceImpl) GetTableUserList() ([]TableUser, bool) {
	return nil, false
}

func (usi *UserServiceImpl) GetTableUserByUsername(name string) (TableUser, bool) {
	return TableUser{}, false
}

func (usi *UserServiceImpl) GetTableUserById(id int64) (TableUser, bool) {
	return TableUser{}, false
}

func (usi *UserServiceImpl) GetUserById(id int64) (User, error) {
	return User{}, nil
}

func (usi *UserServiceImpl) GetUserByIdWithCurId(id int64, curId int64) (User, error) {
	return User{}, nil
}

func NewToken(u *TableUser) string {
	expiresTime := time.Now().Unix() + int64(config.OneDayOfHours)
	fmt.Printf("%v\n", expiresTime)
	claims := jwt.StandardClaims{
		Audience:  u.Name,
		ExpiresAt: expiresTime,
		Id:        strconv.FormatInt(1, 10),
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
