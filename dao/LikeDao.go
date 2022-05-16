package dao

//like表的结构。
type Like struct {
	Id       int64 //自增主键
	User_id  int64 //点赞用户id
	Video_id int64 //视频id
	Cancel   int8  //是否点赞，0为点赞，1为取消赞
}
