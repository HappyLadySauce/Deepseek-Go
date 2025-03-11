# DeepSeek-Go 智能运维监控平台

DeepSeek-Go 是一个集成 AI 的智能运维监控平台，提供系统监控、网络监控、告警管理和 AI 聊天等功能。

## 主要功能

- **系统监控**：实时监控服务器CPU、内存、磁盘等指标
- **网络监控**：监控网络连接、流量和延迟
- **告警管理**：设置告警规则，及时通知异常情况
- **AI 聊天**：集成 DeepSeek AI 模型，提供智能问答和运维建议

## 技术栈

- 前端：Vue 3 + TypeScript + Element Plus
- 后端：Go + Gin + GORM
- AI：DeepSeek 大语言模型

## 最近修复的问题

### 2024-03-11 修复

1. **主题切换功能修复**
   - 修复了 CSS 变量在不同主题间的切换问题
   - 更新了 `base.css` 中的主题类选择器，确保主题变量能全局生效
   - 修改了 `App.vue` 中的主题初始化逻辑，使用 `initTheme()` 函数
   - 更新了 `Aside.vue` 中的主题相关变量，确保菜单颜色正确显示

2. **API 请求路径修复**
   - 修复了 `axios.ts` 中的 baseURL，避免 API 路径重复前缀问题
   - 添加了请求日志，便于调试 API 调用

3. **流式聊天功能修复**
   - 修复了 `chatApi.ts` 中的 `createChatStream` 函数，确保正确的 API 路径
   - 添加了流式连接状态的日志记录

4. **消息发送功能优化**
   - 简化了 `sendMessage` 函数，移除了复杂的流式处理逻辑
   - 添加了更多调试日志，便于追踪消息发送过程
   - 修复了类型错误，确保消息对象符合接口定义

## 使用说明

### 开发环境

1. 克隆仓库
```bash
git clone https://github.com/your-username/deepseek-go.git
cd deepseek-go
```

2. 安装依赖
```bash
# 前端
cd vue
npm install

# 后端
cd ../backend
go mod tidy
```

3. 启动开发服务器
```bash
# 前端
cd vue
npm run dev

# 后端
cd ../backend
go run main.go
```

### 生产环境

1. 构建前端
```bash
cd vue
npm run build
```

2. 构建后端
```bash
cd ../backend
go build -o deepseek-go
```

3. 部署
```bash
./deepseek-go
```

## 配置说明

### 前端配置

前端配置文件位于 `vue/.env` 和 `vue/.env.production`，可以配置 API 地址等参数。

### 后端配置

后端配置文件位于 `backend/config.yaml`，可以配置数据库连接、AI 模型参数等。

## 贡献指南

1. Fork 仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详情请参阅 [LICENSE](LICENSE) 文件。

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
- AI聊天功能（支持历史记录、参数设置、知识库导入）

## AI聊天功能
项目包含一个功能完备的AI聊天模块，提供以下功能：

### 用户界面
- **四栏布局**：左上栏为AI参数设置，左下栏为历史记录，中间栏为聊天框，右栏为知识库导入模块
- **响应式设计**：适配不同尺寸的屏幕
- **暗黑/亮色模式**：与系统整体主题一致

### AI参数设置
- **模型选择**：支持多种模型
  - DeepSeek模型：deepseek-v1-8k, deepseek-v1-32k, deepseek-v1-128k, deepseek-coder
  - Kimi模型：moonshot-v1-8k, moonshot-v1-32k, moonshot-v1-128k, moonshot-v1-auto
- **温度调节**：控制AI回复的多样性（0-1之间）
- **最大Token数**：限制AI回复的最大长度

### 历史记录管理
- **会话列表**：显示所有历史聊天会话
- **创建新会话**：一键开始新的对话
- **切换会话**：在不同会话间快速切换
- **删除会话**：移除不需要的历史会话
- **重命名会话**：自定义会话标题

### 知识库功能
- **文件上传**：支持PDF、Word、TXT等格式文件作为知识库
- **文件管理**：查看和删除已上传的知识库文件
- **上下文增强**：基于知识库内容增强AI回答

### 聊天功能
- **实时对话**：支持与AI模型进行自然语言对话
- **消息历史**：保存完整的对话历史
- **加载状态**：显示AI正在思考的状态
- **快捷发送**：支持Enter快捷键发送消息
- **流式输出**：支持流式输出AI回复，提供更好的用户体验

## API接口文档

所有需要认证的接口都必须在请求头中包含有效的JWT令牌：
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 用户认证接口

#### 用户注册
- **请求**: `POST /api/v1/auth/register`
- **描述**: 注册新用户账号
- **请求体**:
```json
{
  "username": "testuser",
  "password": "password123",
  "email": "test@example.com"
}
```
- **返回值**:
```json
{
  "token": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 用户登录
- **请求**: `POST /api/v1/auth/login`
- **描述**: 使用用户名和密码登录
- **请求体**:
```json
{
  "username": "testuser",
  "password": "password123"
}
```
- **返回值**:
```json
{
  "token": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 聊天相关接口

#### 普通聊天
- **请求**: `POST /api/v1/chat/completions`
- **描述**: 发送消息获取AI回复（非流式）
- **请求体**:
```json
{
  "session_id": 0,          // 会话ID，0表示创建新会话
  "message": "你好，请介绍一下你自己",
  "ai_config_id": 0,        // AI配置ID，0表示使用默认配置
  "knowledge_ids": [1, 2]   // 知识库ID列表，空数组表示不使用知识库
}
```
- **返回值**:
```json
{
  "message": "聊天成功",
  "data": {
    "id": 123,
    "role": "assistant",
    "content": "你好！我是Deepseek-Chat，一个基于AI的智能助手。我可以回答问题、提供信息、参与对话等。有什么我可以帮助你的吗？",
    "created_at": "2025-03-11T15:30:45Z"
  },
  "session": {
    "id": 45,
    "title": "新对话",
    "last_message": "你好！我是Deepseek-Chat...",
    "updated_at": "2025-03-11T15:30:45Z"
  }
}
```

#### 流式聊天
- **请求**: `POST /api/v1/chat/stream`
- **描述**: 流式获取AI回复，实时显示
- **请求体**:
```json
{
  "session_id": 0,
  "message": "写一篇关于人工智能的短文",
  "ai_config_id": 0,
  "knowledge_ids": []
}
```
- **返回值**: Server-Sent Events (SSE) 格式的流式数据
```
data: {"id":"chatcmpl-123","delta":{"content":"人"},"finish_reason":null}

data: {"id":"chatcmpl-123","delta":{"content":"工"},"finish_reason":null}

data: {"id":"chatcmpl-123","delta":{"content":"智"},"finish_reason":null}

data: {"id":"chatcmpl-123","delta":{"content":"能"},"finish_reason":null}

...

data: {"id":"chatcmpl-123","delta":{"content":"."},"finish_reason":"stop"}

data: [DONE]
```

#### 获取会话列表
- **请求**: `GET /api/v1/chat/sessions?page=1&page_size=10`
- **描述**: 获取用户的聊天会话列表
- **返回值**:
```json
{
  "message": "获取会话列表成功",
  "data": {
    "sessions": [
      {
        "id": 45,
        "title": "关于AI的讨论",
        "last_message": "人工智能正在改变我们的生活...",
        "created_at": "2025-03-11T15:30:45Z",
        "updated_at": "2025-03-11T15:45:12Z"
      },
      {
        "id": 44,
        "title": "编程问题解答",
        "last_message": "你可以使用递归算法解决这个问题...",
        "created_at": "2025-03-10T09:20:33Z",
        "updated_at": "2025-03-10T09:35:18Z"
      }
    ],
    "total": 24,
    "page": 1,
    "page_size": 10
  }
}
```

#### 获取会话消息
- **请求**: `GET /api/v1/chat/sessions/45?page=1&page_size=20`
- **描述**: 获取特定会话的消息历史
- **返回值**:
```json
{
  "message": "获取会话消息成功",
  "data": {
    "messages": [
      {
        "id": 301,
        "role": "user",
        "content": "你好，请介绍一下你自己",
        "created_at": "2025-03-11T15:30:40Z"
      },
      {
        "id": 302,
        "role": "assistant",
        "content": "你好！我是Deepseek-Chat，一个基于AI的智能助手。我可以回答问题、提供信息、参与对话等。有什么我可以帮助你的吗？",
        "created_at": "2025-03-11T15:30:45Z"
      }
    ],
    "total": 2,
    "page": 1,
    "page_size": 20
  }
}
```

#### 更新会话信息
- **请求**: `PUT /api/v1/chat/sessions/45`
- **描述**: 更新会话标题等信息
- **请求体**:
```json
{
  "title": "AI助手使用指南"
}
```
- **返回值**:
```json
{
  "message": "更新会话成功",
  "data": {
    "id": 45,
    "title": "AI助手使用指南",
    "last_message": "你好！我是Deepseek-Chat...",
    "updated_at": "2025-03-11T16:05:22Z"
  }
}
```

#### 删除会话
- **请求**: `DELETE /api/v1/chat/sessions/45`
- **描述**: 删除指定的聊天会话
- **返回值**:
```json
{
  "message": "删除会话成功"
}
```

### 知识库相关接口

#### 上传知识库文件
- **请求**: `POST /api/v1/knowledge/upload`
- **描述**: 上传文件到知识库（使用multipart/form-data）
- **表单字段**:
  - `file`: 文件数据
- **返回值**:
```json
{
  "message": "文件上传成功",
  "data": {
    "id": 12,
    "file_name": "AI发展白皮书.pdf",
    "file_size": 1254789,
    "file_type": ".pdf",
    "status": "processing",
    "created_at": "2025-03-11T16:10:33Z"
  }
}
```

#### 获取文件列表
- **请求**: `GET /api/v1/knowledge/files?page=1&page_size=10`
- **描述**: 获取用户上传的知识库文件列表
- **返回值**:
```json
{
  "message": "获取文件列表成功",
  "data": {
    "files": [
      {
        "id": 12,
        "file_name": "AI发展白皮书.pdf",
        "file_size": 1254789,
        "file_type": ".pdf",
        "status": "completed",
        "created_at": "2025-03-11T16:10:33Z"
      },
      {
        "id": 11,
        "file_name": "机器学习入门.docx",
        "file_size": 852147,
        "file_type": ".docx",
        "status": "completed",
        "created_at": "2025-03-10T14:25:17Z"
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 10
  }
}
```

#### 获取文件详情
- **请求**: `GET /api/v1/knowledge/files/12`
- **描述**: 获取知识库文件的详细信息
- **返回值**:
```json
{
  "message": "获取文件详情成功",
  "data": {
    "id": 12,
    "file_name": "AI发展白皮书.pdf",
    "file_size": 1254789,
    "file_type": ".pdf",
    "status": "completed",
    "vector_count": 45,
    "created_at": "2025-03-11T16:10:33Z",
    "processed_at": "2025-03-11T16:12:05Z"
  }
}
```

#### 删除文件
- **请求**: `DELETE /api/v1/knowledge/files/12`
- **描述**: 删除知识库文件
- **返回值**:
```json
{
  "message": "删除文件成功"
}
```

### AI配置相关接口

#### 创建配置
- **请求**: `POST /api/v1/ai-config/`
- **描述**: 创建AI配置
- **请求体**:
```json
{
  "model_name": "deepseek-chat",
  "temperature": 0.7,
  "max_tokens": 2048,
  "provider": "deepseek",
  "is_default": true
}
```
- **返回值**:
```json
{
  "message": "创建AI配置成功",
  "data": {
    "id": 5,
    "model_name": "deepseek-v1-32k",
    "temperature": 0.7,
    "max_tokens": 2048,
    "provider": "deepseek",
    "is_default": true,
    "created_at": "2025-03-11T16:20:11Z"
  }
}
```

#### 获取所有配置
- **请求**: `GET /api/v1/ai-config/`
- **描述**: 获取用户所有AI配置
- **返回值**:
```json
{
  "message": "获取AI配置成功",
  "data": [
    {
      "id": 5,
      "model_name": "deepseek-v1-32k",
      "temperature": 0.7,
      "max_tokens": 2048,
      "provider": "deepseek",
      "is_default": true,
      "created_at": "2025-03-11T16:20:11Z"
    },
    {
      "id": 4,
      "model_name": "moonshot-v1-8k",
      "temperature": 0.8,
      "max_tokens": 1024,
      "provider": "kimi",
      "is_default": false,
      "created_at": "2025-03-10T11:15:33Z"
    }
  ]
}
```

#### 获取默认配置
- **请求**: `GET /api/v1/ai-config/default`
- **描述**: 获取用户的默认AI配置
- **返回值**:
```json
{
  "message": "获取默认配置成功",
  "data": {
    "id": 5,
    "model_name": "deepseek-v1-32k",
    "temperature": 0.7,
    "max_tokens": 2048,
    "provider": "deepseek",
    "is_default": true,
    "created_at": "2025-03-11T16:20:11Z"
  }
}
```

#### 获取可用模型列表
- **请求**: `GET /api/v1/ai-config/models`
- **描述**: 获取系统支持的AI模型列表
- **返回值**:
```json
{
  "message": "获取可用模型列表成功",
  "data": {
    "deepseek": [
      {
        "name": "deepseek-v1-8k",
        "provider": "deepseek", 
        "description": "基础模型，支持8K上下文"
      },
      {
        "name": "deepseek-v1-32k",
        "provider": "deepseek",
        "description": "基础模型，支持32K上下文"
      },
      {
        "name": "deepseek-v1-128k",
        "provider": "deepseek",
        "description": "基础模型，支持128K上下文"
      },
      {
        "name": "deepseek-coder",
        "provider": "deepseek",
        "description": "代码专用模型"
      }
    ],
    "kimi": [
      {
        "name": "moonshot-v1-8k",
        "provider": "kimi",
        "description": "基础模型，支持8K上下文"
      },
      {
        "name": "moonshot-v1-32k",
        "provider": "kimi",
        "description": "基础模型，支持32K上下文"
      },
      {
        "name": "moonshot-v1-128k",
        "provider": "kimi",
        "description": "基础模型，支持128K上下文"
      },
      {
        "name": "moonshot-v1-auto",
        "provider": "kimi",
        "description": "自动选择模型，根据上下文长度"
      }
    ]
  }
}
```

#### 获取单个配置
- **请求**: `GET /api/v1/ai-config/5`
- **描述**: 获取特定的AI配置
- **返回值**:
```json
{
  "message": "获取配置成功",
  "data": {
    "id": 5,
    "model_name": "deepseek-v1-32k",
    "temperature": 0.7,
    "max_tokens": 2048,
    "provider": "deepseek",
    "is_default": true,
    "created_at": "2025-03-11T16:20:11Z"
  }
}
```

#### 更新配置
- **请求**: `PUT /api/v1/ai-config/5`
- **描述**: 更新特定的AI配置
- **请求体**:
```json
{
  "model_name": "deepseek-v1-32k",
  "temperature": 0.5,
  "max_tokens": 3000,
  "provider": "deepseek",
  "is_default": true
}
```
- **返回值**:
```json
{
  "message": "更新配置成功",
  "data": {
    "id": 5,
    "model_name": "deepseek-v1-32k",
    "temperature": 0.5,
    "max_tokens": 3000,
    "provider": "deepseek",
    "is_default": true,
    "updated_at": "2025-03-11T16:30:22Z"
  }
}
```

#### 删除配置
- **请求**: `DELETE /api/v1/ai-config/5`
- **描述**: 删除特定的AI配置（非默认配置）
- **返回值**:
```json
{
  "message": "删除配置成功"
}
```

### 邮箱验证接口

#### 发送验证码
- **请求**: `POST /api/v1/auth/send-verification`
- **描述**: 向指定邮箱发送验证码
- **请求体**:
```json
{
  "email": "user@example.com"
}
```
- **返回值**:
```json
{
  "message": "验证码已发送到您的邮箱，有效期3分钟"
}
```

#### 验证邮箱
- **请求**: `POST /api/v1/auth/verify-email`
- **描述**: 验证用户输入的验证码
- **请求体**:
```json
{
  "email": "user@example.com",
  "code": "123456"
}
```
- **返回值**:
```json
{
  "message": "邮箱验证成功",
  "verified": true
}
```

## 支持的AI提供商
本项目同时支持两种AI提供商：

1. **DeepSeek**
   - 基础URL: https://api.deepseek.com
   - 模型: deepseek-v1-8k, deepseek-v1-32k, deepseek-v1-128k, deepseek-coder

2. **Kimi (Moonshot AI)**
   - 基础URL: https://api.moonshot.cn/v1
   - 模型: moonshot-v1-8k, moonshot-v1-32k, moonshot-v1-128k, moonshot-v1-auto

### AI配置示例
在`config.yml`文件中配置API密钥和基础URL:

```yaml
ai:
  # 深度求索配置
  deepseek:
    api_key: "your_api_key_here"
    base_url: "https://api.deepseek.com"
  
  # 月之暗面 Kimi 配置
  kimi:
    api_key: "your_api_key_here"
    base_url: "https://api.moonshot.cn/v1"
```

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
│     │  ├─ auth.go       # 认证控制器
│     │  ├─ users.go      # 用户控制器
│     │  ├─ chat.go       # 聊天控制器
│     │  ├─ knowledge.go  # 知识库控制器
│     │  ├─ ai_config.go  # AI配置控制器
│     │  └─ email_verification.go # 邮箱验证控制器
│     ├─ global/         # 全局变量
│     │  └─ global.go
│     ├─ middleware/     # 中间件
│     │  ├─ auth.go
│     │  ├─ cors.go
│     │  └─ logger.go
│     ├─ models/         # 数据模型
│     │  ├─ base.go
│     │  ├─ user.go
│     │  ├─ chat.go      # 聊天相关模型
│     │  └─ email_verification.go # 邮箱验证模型
│     ├─ router/         # 路由配置
│     │  └─ router.go
│     ├─ utils/          # 工具函数
│     │  ├─ auth/        # 认证相关
│     │  ├─ email/       # 邮件相关
│     │  ├─ ai/          # AI模型相关
│     │  │  ├─ model.go  # AI模型接口定义
│     │  │  ├─ deepseek.go # DeepSeek模型实现
│     │  │  └─ kimi.go   # Kimi模型实现
│     │  └─ response/    # 响应相关
│     ├─ test/           # 测试文件
│     ├─ go.mod
│     ├─ go.sum
│     └─ main.go
│
└─ vue/src                   # 前端Vue项目目录
     ├─ assets/             # 静态资源
     ├─ components/         # 组件
     ├─ hooks/              # 自定义钩子
     ├─ router/             # 路由配置
     ├─ stores/             # 状态管理
     │  └─ chat.ts         # 聊天相关状态
     ├─ utils/              # 工具函数
     │  └─ chatApi.ts      # 聊天API封装
     ├─ views/              # 页面视图
     │  └─ chat/           # 聊天相关页面
     │     └─ index.vue    # 聊天主页面
     └─ App.vue            # 根组件
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

## 联系方式
如有问题，请提交Issue或联系项目维护者。

