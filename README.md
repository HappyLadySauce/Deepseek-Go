# Deepseek-Go

## 项目简介
本项目是一个基于Go语言和Vue3的全栈开发练手项目，旨在提供一个现代化的Web应用开发框架示例。项目采用前后端分离架构，展示了如何使用当前流行的技术栈构建可扩展的Web应用。

## 技术栈

### 后端技术
- **Go**: 主要开发语言
- **Gin**: 高性能Web框架
- **Gorm**: 优秀的ORM框架，用于数据库操作
- **JWT**: 用户认证
- **Swagger**: API文档生成

### 前端技术
- **Vue3**: 渐进式JavaScript框架
- **TypeScript**: 类型安全的JavaScript超集
- **Vite**: 新一代前端构建工具
- **Element Plus**: UI组件库
- **Pinia**: 状态管理
- **Vue Router**: 路由管理

## 功能特性
- 用户认证与授权
- RESTful API设计
- 数据库CRUD操作
- 前端页面响应式设计
- 统一的错误处理
- 日志记录

## 项目结构
```
├─ README.md                # 项目说明文档
├─ doc/                     # 项目文档目录
│  ├─ api/                 # API文档
│  └─ design/              # 设计文档
│
├─ go/                     # 后端Go项目目录
│  └─ src/
│     ├─ api/             # API接口定义
│     │  └─ v1/           # API版本控制
│     ├─ config/          # 配置文件
│     │  ├─ config.go
│     │  ├─ config.example.yml
│     │  └─ db.go
│     ├─ controller/      # 控制器层
│     │  ├─ auth.go
│     │  └─ user.go
│     ├─ global/         # 全局变量
│     │  └─ global.go
│     ├─ middleware/     # 中间件
│     │  ├─ auth.go
│     │  ├─ cors.go
│     │  └─ logger.go
│     ├─ models/         # 数据模型
│     │  ├─ base.go
│     │  └─ user.go
│     ├─ router/         # 路由配置
│     │  └─ router.go
│     ├─ service/        # 业务逻辑层
│     │  └─ user.go
│     ├─ utils/          # 工具函数
│     │  ├─ auth/
│     │  ├─ email/
│     │  └─ response/
│     ├─ test/           # 测试文件
│     ├─ go.mod
│     ├─ go.sum
│     └─ main.go
│
└─ vue/src                   # 前端Vue项目目录
```

## 快速开始

### 环境要求
- Go 1.23+
- Node.js 16+
- MySQL 8.0+

### 后端启动
```bash
cd go/src
go mod tidy
go run main.go
```

### 前端启动
```bash
cd vue/src
npm install
npm run dev
```

## 开发规范
- 代码风格遵循Go官方规范
- 使用ESLint进行代码检查
- 提交信息遵循Conventional Commits规范

## 贡献指南
欢迎提交Issue和Pull Request

## 许可证
MIT License

## 联系方式
如有问题，请提交Issue或联系项目维护者。

