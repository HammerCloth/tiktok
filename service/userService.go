package service

import "TikTok/dao"

type UserService interface {
	/*
		个人使用
	*/
	// GetTableUserList 获得全部TableUser对象
	GetTableUserList() []dao.TableUser

	// GetTableUserByUsername 根据username获得TableUser对象
	GetTableUserByUsername(name string) dao.TableUser

	// GetTableUserById 根据user_id获得TableUser对象
	GetTableUserById(id int64) dao.TableUser

	// InsertTableUser 将tableUser插入表内
	InsertTableUser(tableUser *dao.TableUser) bool
	/*
		他人使用
	*/
	// GetUserById 未登录情况下,根据user_id获得User对象
	GetUserById(id int64) (User, error)

	// GetUserByIdWithCurId 已登录(curID)情况下,根据user_id获得User对象
	GetUserByIdWithCurId(id int64, curId int64) (User, error)

	// 根据token返回id
	// 接口:auth中间件,解析完token,将userid放入context
	//(调用方法:直接在context内拿参数"userId"的值)	fmt.Printf("userInfo: %v\n", c.GetString("userId"))
}

// User 最终封装后,controller返回的User结构体
type User struct {
	Id             int64  `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	FollowCount    int64  `json:"follow_count"`
	FollowerCount  int64  `json:"follower_count"`
	IsFollow       bool   `json:"is_follow"`
	TotalFavorited int64  `json:"total_favorited,omitempty"`
	FavoriteCount  int64  `json:"favorite_count,omitempty"`
}
