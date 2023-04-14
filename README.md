## 基于 go-zero 实现网盘系统

>学习视频：<a href="https://www.bilibili.com/video/BV1cr4y1s7H4">【项目实战】基于 go-zero 实现网盘系统</a> \
>学习项目：<a href="https://github.com/GetcharZp/cloud-disk">  Github:  CloudDisk</a>


文档：

-  [XORM - Simple and Powerful ORM for Go](https://xorm.io/zh/)
-  [go-zero](https://go-zero.dev/cn/)


Apifox 生成的在线接口文档，仅供参考：

[https://www.apifox.cn/apidoc/shared-74bcbbbe-b20f-4ed6-afa8-02752683a8c9/api-20729960](https://www.apifox.cn/apidoc/shared-74bcbbbe-b20f-4ed6-afa8-02752683a8c9/api-20729960)

### go-zero 简单使用

Goctl 安装，参考： [Goctl 安装 | go-zero](https://go-zero.dev/cn/docs/goctl/installation)

使用 go-zero 官方的工具生成代码：

```shell
# 创建 API 服务
goctl api new core
```

运行生成的代码：

```go
# 启动服务
go run core.go -f etc/core-api.yaml
```


在 core/internal/logic/corelogic.go 中修改业务逻辑

访问：http://localhost:8888/from/you

**go-zero 业务开发流程总结**：

- 在 core.api 中声明接口格式（包括路径、请求参数结构、响应数据结构）
- 通过以下指令生成接口相关代码：

```go
goctl api go -api core.api -dir . -style goZero
```

- 在生成的代码中 xxx_logic.go 中的 TODO 处补全接口逻辑
- 通过以下指令启动服务：

```go
go run core.go -f etc/core-api.yaml
```

生成的单体 api 服务目录：

```go
|____go.mod
|____go.sum
|____greet.api // api接口与类型定义
|____etc // 网关层配置文件
| |____greet-api.yaml
|____internal
| |____config // 配置-对应etc下配置文件
| | |____config.go
| |____handler // 视图函数层, 路由与处理器
| | |____routes.go
| | |____greethandler.go
| |____logic // 逻辑处理
| | |____greetlogic.go
| |____svc // 依赖资源, 封装 rpc 对象的地方
| | |____servicecontext.go
| |____types // 中间类型
| | |____types.go
|____greet.go // main.go 入口
```

###系统模块
- 用户模块

    - 密码登录 
    - 刷新Authorization 
    - 邮箱注册 
    - 用户详情 
    - 用户容量
- 存储池模块
    - 中心存储池资源管理 
      - 文件上传 
      - 文件秒传 
      - 文件分片上传 
      - 对接腾讯对象存储
    - 个人存储池资源管理 
      - 文件关联存储 
      - 文件列表 
      - 文件名称修改 
      - 文件夹创建 
      - 文件删除 
      - 文件移动
- 文件分享模块 
  - 创建分享记录 
  - 获取资源详情 
  - 资源保存

### 数据库分析

**用户信息**：存储用户基本信息，用于登录

```sql
CREATE TABLE `user_basic` (
	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`identity` varchar(36) DEFAULT NULL,

	`name` varchar(60) DEFAULT NULL,
	`password` varchar(32) DEFAULT NULL,
	`email` varchar(100) DEFAULT NULL,

	`created_at` datetime DEFAULT NULL,
	`updated_at` datetime DEFAULT NULL,
	`deleted_at` datetime DEFAULT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
```

**公共文件存储池**：存储文件信息

```sql
CREATE TABLE `repository_pool` (
	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`identity` varchar(36) DEFAULT NULL,

	`hash` varchar(32) DEFAULT NULL COMMENT '文件的唯一标识',
	`name` varchar(255) DEFAULT NULL COMMENT '文件名称',
	`ext` varchar(30) DEFAULT NULL COMMENT '文件扩展名',
	`size` int(11) DEFAULT NULL COMMENT '文件大小',
	`path` varchar(255) DEFAULT NULL COMMENT '文件路径',

	`created_at` datetime DEFAULT NULL,
	`updated_at` datetime DEFAULT NULL,
	`deleted_at` datetime DEFAULT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
```

**用户存储池**：对公共文件存储池中文件信息的引用

```sql
CREATE TABLE `user_repository` (
	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`identity` varchar(36) DEFAULT NULL,

	`parent_id` int(11) DEFAULT NULL COMMENT '父级文件层级, 0-【文件夹】',
	`user_identity` varchar(36) DEFAULT NULL COMMENT '对应用户的唯一标识',
	`repository_identity` varchar(36) DEFAULT NULL COMMENT '公共池中文件的唯一标识',
	`ext` varchar(255) DEFAULT NULL COMMENT '文件或文件夹类型',
	`name` varchar(255) DEFAULT NULL COMMENT '用户定义的文件名',

	`created_at` datetime DEFAULT NULL,
	`updated_at` datetime DEFAULT NULL,
	`deleted_at` datetime DEFAULT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;
```

**文件分享**：

```sql
CREATE TABLE `share_basic` (
	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`identity` varchar(36) DEFAULT NULL,

	`user_identity` varchar(36) DEFAULT NULL COMMENT '对应用户的唯一标识',
	`repository_identity` varchar(36) DEFAULT NULL COMMENT '公共池中文件的唯一标识',
	`user_repository_identity` varchar(36) DEFAULT NULL COMMENT '用户池子中的唯一标识',
	`expired_time` int(11) DEFAULT NULL COMMENT '失效时间，单位秒,【0-永不失效】',
	`click_num` int(11) DEFAULT '0' COMMENT '点击次数',

	`created_at` datetime DEFAULT NULL,
	`updated_at` datetime DEFAULT NULL,
	`deleted_at` datetime DEFAULT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
```

#### 项目中使用到第三方库

Go 邮箱库： [jordan-wright/email: Robust and flexible email library for Go)](https://github.com/jordan-wright/email)

Go UUID 库： [satori/go.uuid: UUID package for Go](https://github.com/satori/go.uuid)

腾讯云 COS 后台地址： [https://console.cloud.tencent.com/cos/bucket](https://console.cloud.tencent.com/cos/bucket)

腾讯云 COS 帮助文档： [https://cloud.tencent.com/document/product/436/31215](https://cloud.tencent.com/document/product/436/31215)
