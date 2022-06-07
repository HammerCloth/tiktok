package ffmpeg

import (
	"TikTok/config"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"time"
)

/*
将ffmpeg作为一个中间件来调用，
通过SSH的方式，在远程登录FTP服务器
调用部署在服务器上的ffmpeg，来完成视频截图,并存储在对应位置
*/

type Ffmsg struct {
	VideoName string
	ImageName string
}

var ClientSSH *ssh.Client
var Ffchan chan Ffmsg

// InitSSH 建立SSH客户端，但是会不会超时导致无法链接，这个需要做一些措施
func InitSSH() {
	var err error
	//创建sshp登陆配置
	SSHconfig := &ssh.ClientConfig{
		Timeout:         5 * time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            config.UserSSH,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以, 但是不够安全

		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}
	if config.TypeSSH == "password" {
		SSHconfig.Auth = []ssh.AuthMethod{ssh.Password(config.PasswordSSH)}
	}
	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", config.HostSSH, config.PortSSH)
	ClientSSH, err = ssh.Dial("tcp", addr, SSHconfig)
	if err != nil {
		log.Fatal("创建ssh client 失败", err)
	}
	log.Printf("获取到客户端：%v", ClientSSH)
	//建立通道，作为队列使用,并且确立缓冲区大小
	Ffchan = make(chan Ffmsg, config.MaxMsgCount)
	//建立携程用于派遣
	go dispatcher()
	go keepAlive()
}

//通过增加携程，将获取的信息进行派遣，当信息处理失败之后，还会将处理方式放入通道形成的队列中
func dispatcher() {
	for ffmsg := range Ffchan {
		go func(f Ffmsg) {
			err := Ffmpeg(f.VideoName, f.ImageName)
			if err != nil {
				Ffchan <- f
				log.Fatal("派遣失败：重新派遣")
			}
			log.Printf("视频%v截图处理成功", f.VideoName)
		}(ffmsg)
	}
}

// Ffmpeg 通过远程调用ffmpeg命令来创建视频截图
func Ffmpeg(videoName string, imageName string) error {
	session, err := ClientSSH.NewSession()
	if err != nil {
		log.Fatal("创建ssh session 失败", err)
	}
	defer session.Close()
	//执行远程命令 ffmpeg -ss 00:00:01 -i /home/ftpuser/video/1.mp4 -vframes 1 /home/ftpuser/images/4.jpg
	combo, err := session.CombinedOutput("ls;/usr/local/ffmpeg/bin/ffmpeg -ss 00:00:01 -i /home/ftpuser/video/" + videoName + ".mp4 -vframes 1 /home/ftpuser/images/" + imageName + ".jpg")
	if err != nil {
		//log.Fatal("远程执行cmd 失败", err)
		log.Fatal("命令输出:", string(combo))
		return err
	}
	//fmt.Println("命令输出:", string(combo))
	return nil
}

//维持长链接
func keepAlive() {
	time.Sleep(time.Duration(config.SSHHeartbeatTime) * time.Second)
	session, _ := ClientSSH.NewSession()
	session.Close()
}
