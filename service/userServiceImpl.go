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

type UserServiceImpl struct {
	FollowService
	LikeService
}

// GetTableUserList 获得全部TableUser对象
func (usi *UserServiceImpl) GetTableUserList() []dao.TableUser {
	tableUsers, err := dao.GetTableUserList()
	if err != nil {
		log.Println("Err:", err.Error())
		return tableUsers
	}
	return tableUsers
}

// GetTableUserByUsername 根据username获得TableUser对象
func (usi *UserServiceImpl) GetTableUserByUsername(name string) dao.TableUser {
	tableUser, err := dao.GetTableUserByUsername(name)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return tableUser
	}
	log.Println("Query User Success")
	return tableUser
}

// GetTableUserById 根据user_id获得TableUser对象
func (usi *UserServiceImpl) GetTableUserById(id int64) dao.TableUser {
	tableUser, err := dao.GetTableUserById(id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return tableUser
	}
	log.Println("Query User Success")
	return tableUser
}

// InsertTableUser 将tableUser插入表内
func (usi *UserServiceImpl) InsertTableUser(tableUser *dao.TableUser) bool {
	flag := dao.InsertTableUser(tableUser)
	if flag == false {
		log.Println("插入失败")
		return false
	}
	return true
}

// GetUserById 未登录情况下,根据user_id获得User对象
func (usi *UserServiceImpl) GetUserById(id int64) (User, error) {
	user := User{
		Id:             0,
		Name:           "",
		FollowCount:    0,
		FollowerCount:  0,
		IsFollow:       false,
		TotalFavorited: 0,
		FavoriteCount:  0,
	}
	tableUser, err := dao.GetTableUserById(id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return user, err
	}
	log.Println("Query User Success")
	followCount, _ := usi.GetFollowingCnt(id)
	if err != nil {
		log.Println("Err:", err.Error())
	}
	followerCount, _ := usi.GetFollowerCnt(id)
	if err != nil {
		log.Println("Err:", err.Error())
	}
	u := GetLikeService() //解决循环依赖
	totalFavorited, _ := u.TotalFavourite(id)
	favoritedCount, _ := u.FavouriteVideoCount(id)
	user = User{
		Id:             id,
		Name:           tableUser.Name,
		FollowCount:    followCount,
		FollowerCount:  followerCount,
		IsFollow:       false,
		TotalFavorited: totalFavorited,
		FavoriteCount:  favoritedCount,
	}
	return user, nil
}

// GetUserByIdWithCurId 已登录(curID)情况下,根据user_id获得User对象
func (usi *UserServiceImpl) GetUserByIdWithCurId(id int64, curId int64) (User, error) {
	user := User{
		Id:             0,
		Name:           "",
		FollowCount:    0,
		FollowerCount:  0,
		IsFollow:       false,
		TotalFavorited: 0,
		FavoriteCount:  0,
	}
	tableUser, err := dao.GetTableUserById(id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return user, err
	}
	log.Println("Query User Success")
	followCount, err := usi.GetFollowingCnt(id)
	if err != nil {
		log.Println("Err:", err.Error())
	}
	followerCount, err := usi.GetFollowerCnt(id)
	if err != nil {
		log.Println("Err:", err.Error())
	}
	isfollow, err := usi.IsFollowing(curId, id)
	if err != nil {
		log.Println("Err:", err.Error())
	}
	u := GetLikeService() //解决循环依赖
	totalFavorited, _ := u.TotalFavourite(id)
	favoritedCount, _ := u.FavouriteVideoCount(id)
	user = User{
		Id:             id,
		Name:           tableUser.Name,
		FollowCount:    followCount,
		FollowerCount:  followerCount,
		IsFollow:       isfollow,
		TotalFavorited: totalFavorited,
		FavoriteCount:  favoritedCount,
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
func NewToken(u dao.TableUser) string {
	expiresTime := time.Now().Unix() + int64(config.OneDayOfHours)
	fmt.Printf("expiresTime: %v\n", expiresTime)
	id64 := u.Id
	fmt.Printf("id: %v\n", strconv.FormatInt(id64, 10))
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
