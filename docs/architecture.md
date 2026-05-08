# 架构说明

## 整体架构

系统采用异步任务模型，将用户请求和 AI 服务调用解耦。

```text
用户
  |
前端
  |
后端 API 服务
  |-- MySQL：存储用户和任务
  |-- Redis：缓存任务完成状态 + 请求限流
  |-- MQ：异步投递任务
        |
      消费端
        |
      AI 服务
```

该设计适合 AI 任务耗时较长、前端需要轮询查询结果的场景。通过 MQ 实现异步解耦，通过 Redis 提升状态查询效率，通过 MySQL 保证核心数据持久化。

## 核心组件职责

### 后端 API 服务

后端 API 服务负责接收前端请求，并完成鉴权、限流、数据写入、任务投递和状态查询。

主要职责：

- 注册和登录；
- 生成并校验 JWT；
- 接收用户提交的 AI 任务；
- 执行 IP 限流和用户限流；
- 创建全局唯一 `task_id`；
- 将任务写入 MySQL；
- 将任务消息投递到 MQ；
- 查询 Redis 并返回任务状态；
- 使用统一响应结构和错误码返回接口结果；

### MySQL

MySQL 负责持久化核心业务数据。

主要存储：

- 用户账号信息；
- 任务基本信息；
- 任务状态；
- AI 返回结果；

### Redis

Redis 在系统中承担两个职责：

1. 存储任务完成状态；
2. 请求限流。

任务完成后，消费端将状态写入 Redis，前端轮询状态时后端优先查询 Redis。

### MQ

MQ 用于解耦任务提交和任务执行。

提交任务接口所在的后端服务作为生产端，将任务消息投递到 MQ。消费端从 MQ 中消费任务，调用 AI 服务，并将处理结果写回 MySQL 和 Redis。

### 消费端

消费端负责从 MQ 中消费任务并调用 AI 服务。

MVP 中消费端不做完整实现，只需要在设计上明确其行为：

- 从 MQ 获取任务消息；
- 调用 AI 服务；
- 成功时更新 MySQL 状态和结果；
- 成功时写入 Redis：`task_id -> success`；
- 失败并达到最大重试次数时更新 MySQL 状态为 `failed`；
- 失败并达到最大重试次数时写入 Redis：`task_id -> failed`。

### Zap 结构化日志

使用 Zap 记录结构化日志，方便后续排查问题。

记录：
- 请求路径；
- 请求方法；
- 用户 ID；
- IP；
- `task_id`；
- 错误信息；
- MQ 投递结果；
- MQ 消费结果；
- AI 服务调用耗时；

### Viper 配置管理

使用 Viper 管理项目配置。

配置内容：
- 服务端口；
- MySQL 配置；
- Redis 配置；
- MQ 配置；
- JWT 密钥和过期时间；
- 限流参数；
- Redis 任务状态过期时间；
- MQ 最大重试次数；
- AI 服务地址；

## Redis 设计

### 任务完成状态缓存

消费端在任务执行完成后，将任务状态写入 Redis。

Key：

```text
task:status:{task_id}
```

Value：

```text
success
```

或：

```text
failed
```

示例：

```text
task:status:abc123 -> success
task:status:def456 -> failed
```

任务状态写入 Redis 后，需要设置过期时间，例如：

```text
TTL = 10 分钟 / 30 分钟 / 1 小时
```

具体时间可以根据业务场景配置。

查询任务状态时：

1. 先查 Redis；
2. 如果 Redis 中不存在该 `task_id` 对应的状态，说明任务尚未完成，返回 `pending`；
3. 如果 Redis 中存在 `success`，返回 `success`；
4. 如果 Redis 中存在 `failed`，返回 `failed`。

后续可以将任务结果也一起缓存到 Redis 中，例如将 value 设计为 JSON，这样成功时可以直接从 Redis 返回结果，减少数据库查询压力。

然后考虑到 Redis 的过期问题，后续建议加入 Redis 查不到要去 MySQL 查。

### 限流设计

系统需要做两类限流：

1. IP 限流；
2. 用户限流；

规则示例：

```text
同一个 IP 一分钟最多 30 次请求
同一个用户一分钟最多 30 次请求
```

IP 限流 Key：

```text
rate:ip:{ip}
```

用户限流 Key：

```text
rate:user:{user_id}
```

可以使用 Redis 的 `INCR + EXPIRE` 实现固定窗口限流。

## MQ 设计

### 生产端职责

生产端指提交任务接口所在的后端服务。

职责如下：
1. 接收用户提交的 AI 任务；
2. 校验 JWT；
3. 执行限流；
4. 创建 `task_id`；
5. 写入 MySQL；
6. 投递任务消息到 MQ；
7. 返回 `task_id` 给前端；

### 消费端职责

消费端负责从 MQ 中消费任务并调用 AI 服务。

职责如下：
1. 从 MQ 中获取任务消息；
2. 调用 AI 服务；
3. 如果调用成功：
   - 更新 MySQL 中的任务状态为 `success`；
   - 写入任务结果；
   - 将 `task_id -> success` 写入 Redis；
   - 向 MQ 发送 ACK；
4. 如果调用失败：
   - 进行重试；
   - 如果未达到最大重试次数，可以继续重试；
   - 如果达到最大重试次数：
     - 更新 MySQL 中的任务状态为 `failed`；
     - 将 `task_id -> failed` 写入 Redis；
     - 结束该任务的消费流程。

### MQ 消息结构

```go
type TaskMessage struct {
    TaskID string `json:"task_id"`
    UserID int64  `json:"user_id"`
    Prompt string `json:"prompt"`
}
```

字段名-类型-说明: 
- `TaskID`  string  任务 ID 
- `UserID`  int64   用户 ID 
- `Prompt`  string  用户提交的任务内容 
