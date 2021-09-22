# CyDrive

## 需求说明
CyDrive 是一个支持自动同步与共享功能的网盘软件，提供基础的文件上传下载与追加功能。此外还提供自动同步与文件共享功能。其存储分离的设计可以让存储容量可以方便地水平扩展，普通人可以使用 CyDrive 来搭建自己的网盘服务，并将自己的若干废弃设备的存储资源同时利用起来。

在日常中我们除了需要经常将文件共享给其它人，也可能会需要共享给自己的其他设备，有时候只是共享一些临时的数据，不需要存储到网盘，如一段文本，一张图片等。因此 CyDrive 也提供了这样的功能，每个用户除了具有自己的网盘存储容量，还具有一个不大的临时数据存储容量，其存储用户共享的临时内容，这个存储区域使用 LRU 与过期机制来淘汰旧数据。

[需求规格说明书](docs/需求规格说明书.md)

## 设计
其分为服务端和客户端两个部分，其中服务端又分为存储层和服务层两层
- 服务端
  - 存储层
    - 提供基础的存储服务
    - 可能提供多副本功能来增强持久性
  - 服务层
    - 支持自定义存储目标（本地，远端机器，或许考虑增加对 S3，OSS 等的支持）
    - 提供文件的分享功能
    - 当用户需要读写数据时，为用户提供对存储层的索引功能
    - 提供用户的账号及其权限管理功能
- 客户端
  - 自动同步（类似于 OneDrive）
  - 与其它人，或自己的其它设备共享数据

### 服务端

将存储层分离是为了存储容量与网络带宽的可扩展性，CyDrive 的目标之一是可以让大家将自己的废弃设备的存储资源利用起来。为了对抗老旧设备的不稳定性（如存储硬件损坏），CyDrive 可能会加入多副本功能。
#### 存储层接口

- 文件的上传下载与追加（append）

#### 服务层接口
- 获取指定文件的索引
- 获取一个目录下的文件列表
- 账号的注册、登录、注销、授权、查询
- 文件的分享（包括临时文件，即不存储在用户网盘中的文件，如临时的图片、消息等）
- 获取系统元信息（如存储层用量，容量）

#### 副本机制
// todo

### 客户端
客户端侧提供跨平台支持，其中桌面平台提供自动同步功能，移动平台受限于续航和网络等因素，可以用户自行设置同步频率或者仅打开应用时进行同步等其它同步策略。另外客户端还提供文件分享功能，对于图片一类常见的文件类型提供预览，而不需下载到本地。

#### 客户端 SDK
- 从/向 服务端 下载/上传/追加 文件
- 自身账号的注册、登录、注销、查询
- 共享文件给其它用户，或其它设备（可以是网盘中的文件，或自己的本地文件）
- 自动同步（指定一个文件夹，将这个文件夹作为网盘根目录，进行自动同步）

除 SDK 外，客户端还需要以下功能：
- 检测自身账号的其它设备在线情况，并可以选择向一些设备发送文件
- 提供常见文件的预览（文本，图片等）


