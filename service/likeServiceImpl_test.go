package service

import (
	"fmt"
	"testing"
)

func TestIsFavourite(t *testing.T) {
	impl := LikeServiceImpl{}
	isFavourite, _ := impl.IsFavourite(666, 3)
	fmt.Printf("%v", isFavourite)
}

func TestFavouriteCount(t *testing.T) {
	impl := LikeServiceImpl{}
	count, _ := impl.FavouriteCount(666)
	fmt.Printf("%v", count)
}
