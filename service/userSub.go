package service

type FollowServiceImpl struct {
}

func (fsi *FollowServiceImpl) IsFollowing(userId int64, targetId int64) (bool, error) {
	return true, nil
}

func (fsi *FollowServiceImpl) GetFollowingCnt(userId int64) (int64, error) {
	return int64(1), nil
}

func (fsi *FollowServiceImpl) GetFollowerCnt(userId int64) (int64, error) {
	return int64(1), nil
}

func (fsi *FollowServiceImpl) AddFollowRelation(userId int64, targetId int64) (bool, error) {
	return true, nil
}

func (fsi *FollowServiceImpl) DeleteFollowRelation(userId int64, targetId int64) (bool, error) {
	return true, nil
}

func (fsi *FollowServiceImpl) GetFollowing(userId int64) ([]User, error) {
	return nil, nil
}

func (fsi *FollowServiceImpl) GetFollowers(userId int64) ([]User, error) {
	return nil, nil
}
