package controller

import (
	"TikTok/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	username := c.Query("username")
	password := c.Query("password")

	usi := service.UserServiceImpl{}

	u := usi.GetTableUserByUsername(username)
	if username == u.Name {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		newUser := service.TableUser{
			Name:     username,
			Password: service.EnCoder(password),
		}
		if usi.InsertTableUser(&newUser) != true {
			println("insert data fail")
		}
		token := service.GenerateToken(username)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   u.Id,
			Token:    token,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	encoderPassword := service.EnCoder(password)
	println(encoderPassword)

	usi := service.UserServiceImpl{}

	u := usi.GetTableUserByUsername(username)

	if encoderPassword == u.Password {
		token := service.GenerateToken(username)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   u.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func UserInfo(c *gin.Context) {
	user_id := c.Query("user_id")
	id, _ := strconv.ParseInt(user_id, 10, 64)

	if u, err := service.UserService.GetUserById(new(service.UserServiceImpl), id); err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     u,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
