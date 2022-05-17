package dao

//评论信息-数据库中的结构体-dao层使用
type Comment struct {
	Id          int64  //评论id
	User_id     int64  //评论用户id
	Video_id    int64  //视频id
	Content     string //评论内容
	Create_date string //评论发布的日期mm-dd
	Cancel      int32  //取消评论为1，发布评论为0
}
