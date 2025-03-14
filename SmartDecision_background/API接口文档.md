# DeepSeek AI 后端 API 接口文档

## 基础URL
所有API请求的基础URL为: `http://localhost:14020/api/v1`

## 路径结构
用户相关API使用 `/auth/` 前缀，例如: `/api/v1/auth/login`

## 请求格式要求
- 所有POST和PUT请求必须使用**JSON格式**(`application/json`)
- 请求体必须正确序列化为JSON字符串
- 客户端必须设置正确的`Content-Type: application/json`请求头

## 用户认证相关接口

### 1. 用户登录
- **URL**: `/auth/login`
- **方法**: `POST`
- **处理函数**: `Login`
- **Content-Type**: `application/json`
- **请求体**:
  ```json
  {
    "username": "用户名",
    "password": "密码"
  }
  ```
- **响应**:
  ```json
  {
    "username": "用户名",
    "email": "用户邮箱",
    "token": "JWT令牌"
  }
  ```
- **错误响应**:
  ```json
  {
    "error": "错误信息"
  }
  ```

### 2. 用户注册
- **URL**: `/auth/register`
- **方法**: `POST`
- **处理函数**: `Register`
- **Content-Type**: `application/json`
- **请求体**:
  ```json
  {
    "username": "用户名",
    "password": "密码",
    "email": "邮箱",
    "verificationCode": "验证码"
  }
  ```
- **响应**:
  ```json
  {
    "message": "注册成功"
  }
  ```
- **错误响应**:
  ```json
  {
    "error": "错误信息"
  }
  ```

### 3. 发送验证邮件
- **URL**: `/auth/send-verification-email`
- **方法**: `POST`
- **处理函数**: `SendVerificationEmail`
- **Content-Type**: `application/json`
- **请求体**:
  ```json
  {
    "email": "邮箱地址"
  }
  ```
- **响应**:
  ```json
  {
    "email": "邮箱地址",
    "message": "验证邮件发送成功"
  }
  ```
- **错误响应**:
  ```json
  {
    "error": "错误信息"
  }
  ```

### 4. 验证验证码
- **URL**: `/auth/verify-verification-code`
- **方法**: `POST`
- **处理函数**: `VerifyCode`
- **Content-Type**: `application/json`
- **请求体**:
  ```json
  {
    "email": "邮箱地址",
    "code": "验证码"
  }
  ```
- **响应**:
  ```json
  {
    "message": "验证验证码成功",
    "valid": true
  }
  ```
- **错误响应**:
  ```json
  {
    "error": "错误信息"
  }
  ```

### 5. 重置密码
- **URL**: `/auth/reset-password`
- **方法**: `POST`
- **处理函数**: `ResetPassword`
- **Content-Type**: `application/json`
- **请求体**:
  ```json
  {
    "email": "邮箱地址",
    "code": "验证码",
    "newPassword": "新密码"
  }
  ```
- **响应**:
  ```json
  {
    "message": "密码重置成功"
  }
  ```
- **错误响应**:
  ```json
  {
    "error": "错误信息"
  }
  ```

### 6. 更新用户信息
- **URL**: `/auth/update-profile`
- **方法**: `PUT`
- **处理函数**: `UpdateProfile`
- **Content-Type**: `application/json`
- **请求头**: 需要包含授权令牌 `Authorization: Bearer {token}`
- **请求体**:
  ```json
  {
    "username": "新用户名",
    "email": "新邮箱（可选）",
    "currentPassword": "当前密码（如果要更改密码）",
    "newPassword": "新密码（可选）"
  }
  ```
- **响应**:
  ```json
  {
    "username": "更新后的用户名",
    "email": "更新后的邮箱"
  }
  ```
- **错误响应**:
  ```json
  {
    "error": "错误信息"
  }
  ```

## 聊天相关接口

### 1. 发送消息
- **URL**: `/chat`
- **方法**: `POST`
- **处理函数**: `SendMessage`
- **请求头**: 需要包含授权令牌 `Authorization: Bearer {token}`
- **请求体**:
  ```json
  {
    "message": "用户消息内容",
    "sessionId": "会话ID（可选，若不存在则创建新会话）"
  }
  ```
- **响应**:
  ```json
  {
    "reply": "AI回复内容",
    "sessionId": "会话ID"
  }
  ```
- **错误响应**:
  ```json
  {
    "error": "错误信息"
  }
  ```

### 2. 获取聊天历史
- **URL**: `/chat-history`
- **方法**: `GET`
- **处理函数**: `GetChatHistory`
- **请求头**: 需要包含授权令牌 `Authorization: Bearer {token}`
- **请求参数**:
  ```
  ?sessionId=会话ID
  ```
- **响应**:
  ```json
  {
    "messages": [
      {
        "id": "消息ID",
        "content": "消息内容",
        "role": "user/assistant",
        "timestamp": "时间戳"
      }
    ]
  }
  ```
- **错误响应**:
  ```json
  {
    "error": "错误信息"
  }
  ```

### 3. 创建新会话
- **URL**: `/chat-sessions`
- **方法**: `POST`
- **处理函数**: `CreateSession`
- **请求头**: 需要包含授权令牌 `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "sessionId": "新会话ID",
    "title": "新会话标题"
  }
  ```
- **错误响应**:
  ```json
  {
    "error": "错误信息"
  }
  ```

### 4. 获取所有会话
- **URL**: `/chat-sessions`
- **方法**: `GET`
- **处理函数**: `GetSessions`
- **请求头**: 需要包含授权令牌 `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "sessions": [
      {
        "id": "会话ID",
        "title": "会话标题",
        "lastMessageTime": "最后消息时间"
      }
    ]
  }
  ```
- **错误响应**:
  ```json
  {
    "error": "错误信息"
  }
  ```

### 5. 删除会话
- **URL**: `/chat-sessions/{sessionId}`
- **方法**: `DELETE`
- **处理函数**: `DeleteSession`
- **请求头**: 需要包含授权令牌 `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "message": "会话已删除"
  }
  ```
- **错误响应**:
  ```json
  {
    "error": "错误信息"
  }
  ``` 