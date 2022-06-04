package config

import "time"

// Secret 密钥
var Secret = "tiktok"

// OneDayOfHours 时间
var OneDayOfHours = 60 * 60 * 24
var OneMinute = 60 * 1
var OneMonth = 60 * 60 * 24 * 30
var OneYear = 365 * 60 * 60 * 24
var ExpireTime = time.Hour * 48 // 设置Redis数据热度消散时间。

// VideoCount 每次获取视频流的数量
const VideoCount = 5

// ConConfig ftp服务器地址
const ConConfig = "43.138.25.60:21"
const FtpUser = "ftpuser"
const FtpPsw = "424193726"
const HeartbeatTime = 2 * 60

// PlayUrlPrefix 存储的图片和视频的链接
const PlayUrlPrefix = "http://43.138.25.60/"
const CoverUrlPrefix = "http://43.138.25.60/images/"

// HostSSH SSH配置
const HostSSH = "43.138.25.60"
const UserSSH = "ftpuser"
const PasswordSSH = "424193726"
const TypeSSH = "password"
const PortSSH = 22
const MaxMsgCount = 100
const SSHHeartbeatTime = 10 * 60

const ValidComment = 0   //评论状态：有效
const InvalidComment = 1 //评论状态：取消
const DateTime = "2006-01-02 15:04:05"

//const ChanCapacity = 10 //chan管道容量，暂时没定

const IsLike = 0     //点赞的状态
const Unlike = 1     //取消赞的状态
const LikeAction = 1 //点赞的行为
const Attempts = 3   //操作数据库的最大尝试次数

const DefaultRedisValue = -1 //redis中key对应的预设值，防脏读
