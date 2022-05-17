package dao

// TableUser 对应数据库User表结构的结构体
type TableUser struct {
	Id       int64
	Name     string
	Password string
}
