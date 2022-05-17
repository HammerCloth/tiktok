package service

//发表评论-使用的结构体-service层引用dao层↑的Comment。
//接口定义(CommentService.go)
type CommentService interface {
	/*
		一、其他同学（video）需要使用的方法：
	*/
	//1.根据videoid获取视频评论数量的接口
	CountFromVideoId(id int64) (int64, error)

	/*
		二、评论模块自己request实现的方法：
	*/
	//2、发表评论，传进来评论的基本信息，返回保存是否成功的状态描述
	Send(comment *Comment) error
	//3、删除评论，传入评论id即可，返回错误状态信息
	DelComment(id int64) error
	//4、查看评论列表-返回评论list-在controller层再封装外层的状态信息
	GetList(vedioId int64, userId int64) ([]CommentInfo, error)
}

//查看评论-传出的结构体-service
type CommentInfo struct {
	Id          int64  //评论id
	UserInfo    User   //评论用户的信息-由用户模块赞助
	Content     string //评论内容
	Create_date string //评论发布的日期mm-dd
}

//评论信息-数据库中的结构体-dao层使用-先放这，import不了
type Comment struct {
	Id          int64  //评论id
	User_id     int64  //评论用户id
	Video_id    int64  //视频id
	Content     string //评论内容
	Create_date string //评论发布的日期mm-dd
	Cancel      int32  //取消评论为1，发布评论为0
}
