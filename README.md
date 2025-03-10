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
- 邮箱验证功能（限制邮箱长度和类型）

## 邮箱验证功能
项目支持邮箱验证功能，并对用户注册时使用的邮箱进行了限制：

### 支持的邮箱类型
- **QQ邮箱**：用户名部分必须为纯数字，长度为5-11位数字
- **163邮箱**：用户名部分长度为6-18位字符
- **Gmail邮箱**：用户名部分最多30个字符

### 邮件服务器配置
系统支持通过SSL/TLS安全连接到邮件服务器，特别适用于QQ邮箱等要求加密连接的邮件服务。配置示例：

```yaml
email:
  host: "smtp.qq.com"        # SMTP服务器地址
  port: 465                  # SSL端口通常为465
  username: "your_email@qq.com"  # 完整邮箱地址
  password: "authorization_code" # QQ邮箱授权码
  from: "系统名称 <your_email@qq.com>"  # 必须使用 "名称 <邮箱>" 格式
  enable_ssl: true           # 启用SSL连接
  server_name: "smtp.qq.com" # 服务器名称
```

> **重要提示**：
> 1. QQ邮箱需要在QQ邮箱设置中开启SMTP服务并获取授权码，授权码用作password配置项
> 2. from字段必须严格遵循 `"名称 <邮箱地址>"` 格式，例如 `"DeepSeek系统 <12345678@qq.com>"`，否则QQ邮箱服务器会拒绝发送
> 3. username字段应该是纯邮箱地址，与from中的邮箱部分保持一致

### 验证流程
1. 用户注册前，先调用发送验证码接口
2. 用户在3分钟内需要输入正确的验证码完成邮箱验证
3. 邮箱验证成功后，用户才能完成注册流程

### 相关API接口
- **发送验证码**：`POST /auth/send-verification`
- **验证邮箱**：`POST /auth/verify-email`

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

