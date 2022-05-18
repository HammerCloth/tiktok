package service

// FollowService 定义用户关系接口以及用户关系中的各种方法
type FollowService interface {
	/*
	   一、其他同学需要调用的业务方法。
	*/
	// IsFollowing 根据当前用户id和目标用户id来判断当前用户是否关注了目标用户
	IsFollowing(userId int64, targetId int64) (bool, error)
	// GetFollowerCnt 根据用户id来查询用户被多少其他用户关注
	GetFollowerCnt(userId int64) (int64, error)
	// GetFollowingCnt 根据用户id来查询用户关注了多少其它用户
	GetFollowingCnt(userId int64) (int64, error)
	/*
	   二、直接request需要的业务方法
	*/
	// AddFollowRelation 当前用户关注目标用户
	AddFollowRelation(userId int64, targetId int64) (bool, error)
	// DeleteFollowRelation 当前用户取消对目标用户的关注
	DeleteFollowRelation(userId int64, targetId int64) (bool, error)
	// GetFollowing 获取当前用户的关注列表
	GetFollowing(userId int64) ([]User, error)
	// GetFollowers 获取当前用户的粉丝列表
	GetFollowers(userId int64) ([]User, error)
}
