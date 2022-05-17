package service

type UserService interface {
	/*
		个人使用
	*/
	// GetTableUserList GetUserList 获得全部TableUser对象
	GetTableUserList() []TableUser

	// GetTableUserByUsername GetUserByUsername 根据username获得TableUser对象
	GetTableUserByUsername(name string) TableUser

	// GetTableUserById GetUserById 根据user_id获得TableUser对象
	GetTableUserById(id int64) TableUser

	// InsertTableUser 将tableUser插入表内
	InsertTableUser(tableUser *TableUser) bool
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

type User struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

type TableUser struct {
	Id       int64
	Name     string
	Password string
}
