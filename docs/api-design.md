# HTTP API 设计

系统主要包含四个接口。

1. 用户注册
- 请求方法：`POST`
- 请求路径：`/api/v1/auth/register`
- 是否需要JWT：否
- 说明：用于用户注册，创建新的用户账号。

2. 用户登录
- 请求方法：`POST`
- 请求路径：`/api/v1/auth/login`
- 是否需要JWT：否
- 说明：用于用户登录，登录成功后返回 JWT。

3. 提交 AI 任务
- 请求方法：`POST`
- 请求路径：`/api/v1/tasks`
- 是否需要JWT：是
- 说明：用户提交一个 AI 任务，服务端异步处理。

4. 查询任务状态
- 请求方法：`GET`
- 请求路径：`/api/v1/tasks/{task_id}`
- 是否需要JWT：是
- 说明：前端根据任务 ID 轮询查询任务处理状态。

## 统一响应结构

所有接口返回给客户端的数据，都使用统一响应结构。

```go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}
```

成功响应示例：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "abc123",
    "status": "pending"
  }
}
```

失败响应示例：

```json
{
  "code": 40000,
  "message": "invalid params"
}
```

## 错误码

```go
const (
    CodeSuccess = 0

    CodeInvalidParams = 40000
    CodeUnauthorized  = 40100
    CodeForbidden     = 40300
    CodeNotFound      = 40400
    CodeTooManyReq    = 42900

    CodeUserExists       = 40001
    CodeInvalidLogin     = 40002
    CodeTaskNotFound     = 40003
    CodeTaskSubmitFailed = 50001

    CodeInternalError = 50000
)
```

错误码及其含义：
- `0`：成功
- `40000`：请求参数错误
- `40100`：未认证或 Token 无效
- `40300`：无权限
- `40400`：资源不存在
- `42900`：请求过于频繁
- `40001`：用户名已存在
- `40002`：用户名或密码错误
- `40003`：任务不存在
- `50001`：任务提交失败
- `50000`：系统内部错误

## 注册接口

### 接口说明

用户提交用户名和密码后，后端完成以下操作：

1. 校验参数；
2. 查询用户名是否已经存在；
3. 对密码进行哈希；
4. 写入 `users` 表；
5. 返回用户基本信息。

### 请求信息

```text
POST /api/v1/auth/register
```

### 请求结构体

```go
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=64"`
    Password string `json:"password" binding:"required,min=6,max=64"`
}
```

### 响应结构体

```go
type RegisterResponse struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 1,
    "username": "test_user"
  }
}
```

## 登录接口

### 接口说明

用户提交用户名和密码后，后端完成以下操作：

1. 根据用户名查询用户；
2. 校验密码哈希；
3. 生成 JWT；
4. 返回 Token 和过期时间。

### 请求信息

```text
POST /api/v1/auth/login
```

### 请求结构体

```go
type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}
```

### 响应结构体

```go
type LoginResponse struct {
    Token     string `json:"token"`
    ExpiresIn int64  `json:"expires_in"`
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.xxx",
    "expires_in": 7200
  }
}
```

## 提交 AI 任务接口

### 接口说明

用户携带 JWT 提交 prompt 后，后端完成以下操作：

1. 校验 JWT；
2. 获取 `user_id`；
3. 执行 IP 限流；
4. 执行用户限流；
5. 创建全局唯一的 `task_id`；
6. 写入 MySQL，状态为 `pending`；
7. 投递任务消息到 MQ；
8. 返回 `task_id` 给前端。

### 请求信息

```text
POST /api/v1/tasks
```

### 请求头

```http
Authorization: Bearer <token>
```

### 请求结构体

```go
type CreateTaskRequest struct {
    Prompt string `json:"prompt" binding:"required,min=1,max=5000"`
}
```

### 响应结构体

```go
type CreateTaskResponse struct {
    TaskID string `json:"task_id"`
    Status string `json:"status"`
}
```

### 请求示例

```json
{
  "prompt": "请帮我总结这段文本"
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "task_abc123",
    "status": "pending"
  }
}
```

## 查询任务状态接口

### 接口说明

这是前端定时轮询调用的接口。前端拿到 `task_id` 后，可以每隔 3 秒请求一次。

后端完成以下操作：

1. 校验 JWT；
2. 获取 `user_id`；
3. 校验任务是否属于当前用户；
4. 查询 Redis；
5. Redis 查不到，返回 `pending`；
6. Redis 查到 `success`，返回 `success`；
7. Redis 查到 `failed`，返回 `failed`。

### 请求信息

```text
GET /api/v1/tasks/{task_id}
```

### 请求头

```http
Authorization: Bearer <token>
```

### 请求结构体

```go
type GetTaskStatusRequest struct {
    TaskID string `uri:"task_id" binding:"required"`
}
```

### 响应结构体

```go
type GetTaskStatusResponse struct {
    TaskID string `json:"task_id"`
    Status string `json:"status"`
}
```

### 响应示例：任务处理中

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "task_abc123",
    "status": "pending"
  }
}
```

### 响应示例：任务成功

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "task_abc123",
    "status": "success"
  }
}
```

### 响应示例：任务失败

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "task_abc123",
    "status": "failed"
  }
}
```
