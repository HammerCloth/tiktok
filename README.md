# tiktok

<!-- PROJECT SHIELDS -->

![GitHub Repo stars](https://img.shields.io/github/stars/HammerCloth/tiktok?style=plastic)
![GitHub watchers](https://img.shields.io/github/watchers/HammerCloth/tiktok?style=plastic)
![GitHub forks](https://img.shields.io/github/forks/HammerCloth/tiktok?style=plastic)
![GitHub contributors](https://img.shields.io/github/contributors/HammerCloth/tiktok)
[![MIT License][license-shield]][license-url]


<!-- PROJECT LOGO -->
<br />

<p align="center">
  <a href="https://github.com/HammerCloth/tiktok.git/">
    <img src="images/logo2.png" alt="Logo" width="400" height="100">
  </a>

<h3 align="center">抖音简洁版</h3>
  <p align="center">
    <br />
    <a href="https://github.com/HammerCloth/tiktok.git"><strong>探索本项目的文档 »</strong></a>
    <br />
    <br />
  </p>
  </p>

**Attention:** We always welcome contributors to the project. Before adding your contribution, please carefully read our [Git 分支管理规范](https://ypbg9olvt2.feishu.cn/docs/doccnTMRmh7YgMwL2PgZ5moWUsd)和[注释规范](https://juejin.cn/post/7096881555246678046)。

## 目录
- [荣誉展示](#荣誉展示)
- [上手指南](#上手指南)
    - [开发前的配置要求](#开发前的配置要求)
    - [安装步骤](#安装步骤)
    - [演示界面](#演示界面)
    - [演示视频](#演示视频)
- [文件目录说明](#文件目录说明)
- [开发的整体设计](#开发的整体设计)
   - [整体的架构图](#整体的架构图)
   - [数据库的设计](#数据库的设计)
   - [Redis架构的设计](#Redis架构的设计)
   - [RabbitMQ架构的设计](#RabbitMQ架构的设计)
   - [服务模块的设计](#服务模块的设计)
     - [视频模块的设计](#视频模块的设计)
     - [点赞模块的设计](#点赞模块的设计)
     - [关注模块的设计](#关注模块的设计)
     - [用户模块的设计](#用户模块的设计)
     - [评论模块的设计](#评论模块的设计)
- [性能测试](#性能测试)
- [编译项目到linux](#编译项目到linux)
- [使用到的技术](#使用到的技术)
- [未来展望](#未来展望)
  - [分布式服务](#分布式服务)
  - [推荐视频展望](#推荐视频展望)
- [如何参与开源项目](#如何参与开源项目)
- [版本控制](#版本控制)
- [贡献者](#贡献者)
- [鸣谢](#鸣谢)

### 荣誉展示
<a href="https://github.com/HammerCloth/tiktok.git/">
    <img src="images/awards.jpg" alt="Logo" width="800" height="400">
    <img src="images/awards1.jpg" alt="Logo" width="800" height="400">
</a>

### 上手指南

#### 开发前的配置要求

1. go 1.18.1
2. MySQL(数据库sql文件在config包中)
3. 搭建Redis、RabbitMQ环境
4. 配置静态资源服务器：安装Nginx、vsftpd、ffmpeg（相关配置文件在config包中）
5. [最新版抖音客户端软件](https://pan.baidu.com/s/1kXjvYWH12uhvFBARRMBCGg?pwd=6cos)



#### 安装步骤
1. 下载源码
2. 配置SSH、FTP、Redis、静态服务器地址等相关参数
3. 启动服务
4. 在客户端配置相关地址服务端地址即可

```sh
git clone https://github.com/HammerCloth/tiktok.git
```
#### 演示界面
**基础功能演示**

<a href="https://github.com/HammerCloth/tiktok.git/">
    <img src="images/1.png" alt="Logo" width="200" height="400">
    <img src="images/2.png" alt="Logo" width="200" height="400">
    <img src="images/3.png" alt="Logo" width="200" height="400">
    <img src="images/4.png" alt="Logo" width="200" height="400">
</a>

**拓展功能演示**

<a href="https://github.com/HammerCloth/tiktok.git/">
    <img src="images/5.png" alt="Logo" width="200" height="400">
    <img src="images/6.png" alt="Logo" width="200" height="400">
    <img src="images/7.png" alt="Logo" width="200" height="400">
    <img src="images/8.png" alt="Logo" width="200" height="400">
</a>

**设置服务端地址**

<a href="https://github.com/HammerCloth/tiktok.git/">
    <img src="images/9.png" alt="Logo" width="200" height="400">
    <img src="images/10.png" alt="Logo" width="200" height="400">
    <img src="images/11.png" alt="Logo" width="200" height="400">
</a>

#### 演示视频
[![Watch the video](images/video.png)](http://43.138.25.60/tiktok.mp4)

### 文件目录说明

```
tiktok 
├── /.idea/
├── /config/ 配置文件包
├── /controller/ 控制器包
├── /dao/ 数据库访问
├── /document/ 敏感词词库
├── /images/ 图片引用
├── /middleware/ 中间件
│   ├── ffmpeg/ 视频截图
│   ├── ftp/ 文件服务器
│   ├── jwt/ 鉴权
│   ├── rabbitmq/ 消息队列
│   ├── redis/ 缓存
├── /service/ 服务层
├── /util/ 工具
├── .gitignore
├── /go.mod/
├── LICENSE
├── main.go
├── README.md
└── router.go
```

### 开发的整体设计
#### 整体的架构图
<p align="center">
  <a href="https://github.com/HammerCloth/tiktok.git/">
    <img src="images/framework.jpg" alt="Logo" width="1000" height="500">
  </a>

#### 数据库的设计
<p align="center">
  <a href="https://github.com/HammerCloth/tiktok.git/">
    <img src="images/mysql.jpg" alt="Logo" width="800" height="600">
  </a>
</p>

#### Redis架构的设计
<p align="center">
  <a href="https://github.com/HammerCloth/tiktok.git/">
    <img src="images/redis.jpg" alt="Logo" width="1000" height="600">
  </a>
</p>

#### RabbitMQ架构的设计
<p align="center">
  <a href="https://github.com/HammerCloth/tiktok.git/">
    <img src="images/rabbitmq.jpg" alt="Logo" width="1000" height="600">
  </a>
</p>

#### 服务模块的设计

###### 视频模块的设计
视频模块包括视频Feed流获取、视频投稿和获取用户投稿列表。
详情请阅读[视频模块设计说明](https://bytedancecampus1.feishu.cn/docs/doccntmcunjHSMzVUNEhGbxjxJh) 查阅为该模块的详细设计。

###### 点赞模块的设计
点赞模块包括点赞视频、取消赞视频和获取点赞列表。
详情请阅读[点赞模块设计说明](https://bytedancecampus1.feishu.cn/docs/doccn13iJgTIAebIPpMiRqb0Hwb) 查阅为该模块的详细设计。

###### 关注模块的设计
关注模块包括关注、取关、获取关注列表、获取粉丝列表四个基本功能。
详情请阅读[关注模块的设计说明](https://bytedancecampus1.feishu.cn/docs/doccnOsdm29SufPJkDfRs7tLHgx) 查阅为该模块的详细设计。

###### 用户模块的设计
用户与安全模块包括用户注册、用户登录和用户信息三个部分
详情请阅读[用户模块的设计说明](https://bytedancecampus1.feishu.cn/docs/doccn1vusmV9oN1ukTCyLpbJ46f) 查阅为该模块的详细设计。

###### 评论模块的设计
评论模块包括发表评论、删除评论和查看评论。
详情阅读[评论模块的设计说明](https://bytedancecampus1.feishu.cn/docs/doccnDqfcZJW4tTD409NGlYfvCb) 查阅为该模块的详细设计。

### 性能测试
通过命令 go tool pprof -http=:6060 "http://localhost:8080/debug/pprof/profile?seconds=120" 生成了两个版本的火焰图，左图为v1.0，右图为v1.2版本，通过对比两张详细火焰图，优化后的相同方法调用时间更短（添加了相应的中间件）

<p align="center">
<a href="https://github.com/HammerCloth/tiktok.git/">
    <img src="images/1.0.png" alt="Logo" width="500" height="300">
    <img src="images/1.2.png" alt="Logo" width="500" height="300">
</a>

同时包含各个接口的压力测试，详情请阅读[压测报告](https://bytedancecampus1.feishu.cn/docs/doccnoDHHJ84k94G1I2TxHj9Udh) 获得具体的压力数据。
### 编译项目到linux
```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./
```

### 使用到的技术
框架相关：
- [Gin](https://gin-gonic.com/docs/)
- [Gorm](https://gorm.io/docs/)

服务器相关：
- [Nginx](https://www.nginx-cn.net/)
- [vsftpd](https://www.linuxfromscratch.org/blfs/view/svn/server/vsftpd.html)
- [ffmpeg](https://ffmpeg.org/documentation.html)
- [goftp](http://t.zoukankan.com/lvdongjie-p-9554849.html)

中间件相关：
- [Redis](https://redis.io/docs/)
- [RabbitMQ](https://www.rabbitmq.com/documentation.html)

数据库：
- [MySQL](https://dev.mysql.com/doc/)

### 未来展望
#### 分布式服务
利用dubbogo来完成分布式，貔貅作为网关，Nacos作为注册中心，将五个模块分别布置到不同的服务器上，以rpc调用的方式来调用当前模块依赖其他模块的方法，做到分布式处理与解耦。
<p align="center">
  <a href="https://github.com/HammerCloth/tiktok.git/">
    <img src="images/future.png" alt="Logo" width="1000" height="600">
  </a>
</p>

#### 推荐视频展望
队伍创新推荐算法
<p align="center">
  <a href="https://github.com/HammerCloth/tiktok.git/">
    <img src="images/recommend.jpg" alt="Logo" width="1000" height="600">
  </a>
</p>

详情请阅读[视频推荐展望](https://bytedancecampus1.feishu.cn/docx/doxcnVIK62rWgR0iE49UAI2wB0b)

### 如何参与开源项目

贡献使开源社区成为一个学习、激励和创造的绝佳场所。你所作的任何贡献都是**非常感谢**的。

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### 版本控制

该项目使用Git进行版本管理。您可以在repository参看当前可用版本。

### 贡献者
- 司一雄 邮箱:18552541076@163.com
- 刘宗舟 邮箱:1245314855@qq.com
- 蒋宇栋 邮箱:jiangyudong123@qq.com
- 李思源 邮箱:yuanlaisini_002@qq.com
- 李林森 邮箱:1412837463@qq.com

*您也可以查阅仓库为该项目做出贡献的开发者。*

### 版权说明

该项目签署了MIT 授权许可，详情请参阅 [LICENSE.txt](https://github.com/shaojintian/Best_README_template/blob/master/LICENSE.txt)

### 鸣谢

- [字节跳动后端青训营](https://youthcamp.bytedance.com/)

<!-- links -->

[license-shield]: https://img.shields.io/github/license/mrxuexi/tiktok.svg?style=flat-square

[license-url]: https://github.com/mrxuexi/tiktok/blob/master/LICENSE.txt
