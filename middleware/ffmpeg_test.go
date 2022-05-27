package middleware

import (
	"TikTok/config"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"testing"
	"time"
)

func TestSSH(t *testing.T) {
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
	sshClient, err := ssh.Dial("tcp", addr, SSHconfig)
	if err != nil {
		log.Fatal("创建ssh client 失败", err)
	}
	defer sshClient.Close()
	//创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		log.Fatal("创建ssh session 失败", err)
	}
	defer session.Close()
	//执行远程命令
	combo, err := session.CombinedOutput("whoami; cd /root/huayun; ls -al;echo Hello > hello.txt;echo hello;curl http://baidu.com")
	if err != nil {
		log.Fatal("远程执行cmd 失败", err)
	}
	fmt.Println("命令输出:", string(combo))
}

func TestGo(t *testing.T) {
	for i := 0; i < 10; i++ {
		go TestSSH(t)
	}

}

func TestFfmpeg(t *testing.T) {
	InitSSH()
	Ffmpeg("1", "33")
	//Ffmpeg("1", "4")
}
