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

