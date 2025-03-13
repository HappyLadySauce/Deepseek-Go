# Golang的CURD

## CURD

`CURD` 是 `create`、`read`、`update` 和 `delete` 的首字母缩写，在数据库操作中频繁出现。通常我们说的CURD就是对数据库的读写操作。

这文章篇主要介绍 `Golang` 对 `Mysql`、`MongoDB` 和 `Redis` 的 `CURD` 相关操作。

### 关系型数据库

在 `Golang` 开发中，关系型数据库通常我们会使用 `gorm.io/gorm` 库。

```bash
go get -u gorm.io/gorm
```

`Grom` （Object-Relationl Mapping）使用结构体映射数据库中的表，将数据库抽象成一个数据库相关的对象。这样，我们在操作数据库的时候，就可以像平时操作对象一样操作它就可以了。

### 非关系型数据库

在 `Golang` 开发中，`MongoDB` 一般使用 MongoDB 官方提供的 `mongo-go-driver` 库；

```bash
go get go.mongodb.org/mongo-driver/mongo
```

`Redis` 则使用 `github.com/go-redis/redis` 库，这是一个广泛使用的 Go 语言 Redis 客户端库。

```bash
# Redis 7
go get github.com/go-redis/redis/v9
```

## Mysql

go\src
│
├── config                 # 配置文件夹，存放配置文件和相关代码
│   ├── config.example.yml # 配置文件的示例，通常用于展示配置项的格式和默认值
│   ├── config.go          # 配置相关的Go代码，可能用于加载和解析配置文件

│   ├── mongodb.go    # MongoDB数据库操作的Go代码，可能包含数据库连接和操作函数

│   └── config.yml         # 项目的实际配置文件，通常在部署时使用
│
├── controller             # 控制器文件夹，存放处理HTTP请求的代码
│   └── global.go          # 全局控制器代码，可能包含一些公共的控制器逻辑或中间件
│
├── middlewares            # 中间件文件夹，存放处理HTTP请求中间逻辑的代码
│   └── cors_middlewares.go # 处理跨域资源共享（CORS）的中间件代码
│
├── models                 # 模型文件夹，存放数据模型的定义
│   └── models.go            # 模型的定义，包含相关的数据结构和方法
│
├── router                 # 路由文件夹，存放定义路由的代码
│   └── router.go          # 路由定义的Go代码，可能包含路由的注册和配置
│
├── go.mod                 # Go模块文件，定义项目的依赖关系
├── go.sum                 # Go依赖文件，记录项目依赖的具体版本
└── main.go                # 项目的入口文件，通常包含main函数，**用于启动应用**































































