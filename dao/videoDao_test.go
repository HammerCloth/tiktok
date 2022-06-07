package dao

import (
	"TikTok/middleware/ftp"
	"fmt"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	Init()
}

func TestFind(t *testing.T) {
	Init()
	var tv TableVideo
	result := Db.First(&tv)
	fmt.Println(result.RowsAffected)
	fmt.Println(tv.Id)
	fmt.Println(tv.AuthorId)
	fmt.Println(tv.CoverUrl)
	fmt.Println(tv.PlayUrl)
	fmt.Println(tv.PublishTime)
	fmt.Println(tv.Title)
}

func TestGetVideosByAuthorId(t *testing.T) {
	Init()
	data, err := GetVideosByAuthorId(2)
	if err != nil {
		print(err)
	}
	for _, video := range data {
		fmt.Println(video)
	}
}

func TestGetVideoByVideoId(t *testing.T) {
	Init()
	data, err := GetVideoByVideoId(3)
	if err != nil {
		print(err)
	}
	fmt.Println(data)

}

func TestGetVideosByLastTime(t *testing.T) {
	Init()
	data, err := GetVideosByLastTime(time.Now())
	if err != nil {
		return
	}
	for _, video := range data {
		fmt.Println(video)
	}
}
func TestVideoFtp(t *testing.T) {
	ftp.InitFTP()
	//file, err := os.Open("/Users/siyixiong/Movies/bilibil/bilibil20211219/樱花少女.mp4")
	//if err != nil {
	//	panic(err)
	//}
	//err = VideoFTP(file, "k2")
	//if err != nil {
	//	return
	//}
	//defer file.Close()
	//ffmpeg.exe -ss 00:00:01 -i spring.mp4 -vframes 1 bb.jpg
	//imageName := uuid.NewV4().String() + ".jpg"
	//cmdArguments := []string{"-ss", "00:00:01", "-i", "/home/ftpuser/video/" + "1" + ".mp4", "-vframes", "1", "/home/ftpuser/images/" + imageName}
	//cmd := exec.Command("ffmpeg", cmdArguments...)
	//var out bytes.Buffer
	//cmd.Stdout = &out
	//err := cmd.Run()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("command output: %q", out.String())
}

func TestSave(t *testing.T) {
	Save("test", "test", 10024, "aaa")
}

func TestGetVideoIdsByAuthorId(t *testing.T) {
	Init()
	id, err := GetVideoIdsByAuthorId(20003)
	if err != nil {
		return
	}
	fmt.Println(id)
}
