package controller

import (
	"TikTok/service"
	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User service.User `json:"user"`
}

func Register(c *gin.Context) {
	/*username := c.Query("username")
	password := c.Query("password")

	u := new(service.User)

	if u.GetUserByUsername(username) {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		newUser := &service.User{
			Name:     username,
			Password: service.EnCoder(password),
		}
		if newUser.InsertUser() != true {
			println("insert data fail")
		}
		token := service.GenerateToken(username)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: pojo.Response{StatusCode: 0},
			UserId:   u.Id,
			Token:    token,
		})
	}*/
}

func Login(c *gin.Context) {
	/*username := c.Query("username")
	password := c.Query("password")
	encoderPassword := service.EnCoder(password)
	println(encoderPassword)
	u := new(service.User)
	u.GetUserByUsername(username)

	if encoderPassword == u.Password {
		token := service.GenerateToken(username)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: pojo.Response{StatusCode: 0},
			UserId:   u.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}*/
}

func UserInfo(c *gin.Context) {
	/*user_id := c.Query("user_id")
	id, _ := strconv.ParseInt(user_id, 10, 64)
	u := new(service.User)
	if  u.GetUserById(id) {
		user := pojo.User{
			Id:            u.Id,
			Name:          u.Name,
			FollowCount:   1,
			FollowerCount: 1,
			IsFollow:      true,
		}
		c.JSON(http.StatusOK, UserResponse{
			Response: pojo.Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}*/
}
