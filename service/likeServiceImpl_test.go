package service

import (
	"fmt"
	"testing"
)

func TestIsFavourit(t *testing.T) {
	impl := LikeServiceImpl{}
	bool, _ := impl.IsFavourit(666, 3)
	fmt.Printf("%v", bool)
}

func TestFavouriteCount(t *testing.T) {
	impl := LikeServiceImpl{}
	count, _ := impl.FavouriteCount(666)
	fmt.Printf("%v", count)
}
