package service

//Service 定义点赞状态和点赞数量
type LikeService interface {
	/*
	   1.其他模块(video)需要使用的业务方法。
	*/
	//IsFavorit 根据当前视频id判断是否点赞了该视频。
	IsFavourit(videoId int64, userId int64) (bool, error)
	//FavouriteCount  根据当前视频id获取当前视频点赞数量。
	FavouriteCount(videoId int64) (int64, error)

	/*
	   2.request需要实现的功能
	*/

	//当前用户对视频的点赞操作 ,并把这个行为更新到like表中。
	//当前操作行为，1点赞，2取消点赞。
	FavouriteAction(userId int64, videoId int64, action_type int32) error
	//获取当前用户的所有点赞视频，调用videoService的方法
	GetFavouriteList(userId int64) ([]Video, error)
}
