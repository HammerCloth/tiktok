package config

// Secret 密钥
var Secret = "tiktok"

// OneDayOfHours 时间
var OneDayOfHours = 60 * 60 * 24
var OneMinute = 60 * 1
var OneMonth = 60 * 60 * 24 * 30

// VideoCount 每次获取视频流的数量
const VideoCount = 5

// ConConfig ftp服务器地址
const ConConfig = "43.138.25.60:21"
const FtpUser = "ftpuser"
const FtpPsw = "424193726"

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
