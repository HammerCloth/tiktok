package ftp

import (
	"TikTok/config"
	"github.com/dutchcoders/goftp"
	"log"
	"time"
)

var MyFTP *goftp.FTP

func InitFTP() {
	//获取到ftp的链接
	var err error
	MyFTP, err = goftp.Connect(config.ConConfig)
	if err != nil {
		log.Printf("获取到FTP链接失败！！！")
	}
	log.Printf("获取到FTP链接成功%v：", MyFTP)
	//登录
	err = MyFTP.Login(config.FtpUser, config.FtpPsw)
	if err != nil {
		log.Printf("FTP登录失败！！！")
	}
	log.Printf("FTP登录成功！！！")
	//维持长链接
	go keepAlive()
}

func keepAlive() {
	time.Sleep(time.Duration(config.HeartbeatTime) * time.Second)
	MyFTP.Noop()
}
