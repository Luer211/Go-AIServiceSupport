# MVP 边界与后续技术债

## MVP 边界

当前 MVP 聚焦后端 API 主链路，不把所有周边能力一次性做完。

MVP 包含：

- 用户注册；
- 用户登录；
- JWT 生成与校验；
- 提交 AI 任务；
- 任务写入 MySQL；
- 任务投递 MQ 的设计与调用入口；
- Redis 存储任务完成状态；
- Redis 做 IP 限流和用户限流；
- 前端轮询查询任务状态；
- 统一响应结构；
- 统一错误码；
- Zap 结构化日志；
- Viper 开发环境配置。

MVP 暂不完整实现：

- 消费端完整服务；
- AI 服务真实调用；
- 复杂 MQ 可靠投递机制；
- 任务结果 Redis 缓存；
- 查询状态时的 MySQL 兜底；
- 多环境完整发布配置；
- 任务取消、运行中等扩展状态。

## 后续可优化点

### Redis 缓存任务结果

当前 Redis 只存储任务完成状态。

后续可以将结果一起存储为 JSON：

```json
{
  "status": "success",
  "result": "AI 返回的结果"
}
```

这样前端查询成功时可以直接拿到结果，减少数据库访问。

### 查询接口兜底查数据库

当前设计中 Redis 查不到就认为任务仍是 `pending`。

但如果 Redis 过期或者丢失，可能会导致已经完成的任务被误判为 `pending`。

更稳妥的方式是：

1. 先查 Redis；
2. Redis 查不到时查 MySQL；
3. 如果 MySQL 中状态是 `success` 或 `failed`，重新写入 Redis；
4. 再返回最终状态。

### 任务状态扩展

后续可以扩展更多任务状态：

```go
const (
    TaskStatusPending   = "pending"
    TaskStatusRunning   = "running"
    TaskStatusSuccess   = "success"
    TaskStatusFailed    = "failed"
    TaskStatusCancelled = "cancelled"
)
```

任务状态说明：
- `pending`：等待消费
- `running`：正在执行
- `success`：执行成功
- `failed`：执行失败
- `cancelled`：已取消

### 任务表增加失败原因字段

当前失败时只更新状态为 `failed`，后续可以增加失败原因字段。

示例：

```sql
ALTER TABLE tasks ADD COLUMN error_message TEXT NULL;
```

这样方便前端展示失败原因，也方便后端排查问题。

### 投递 MQ 失败处理

提交任务接口中，如果 MySQL 写入成功但 MQ 投递失败，需要考虑一致性问题。

可以采用以下方式：

1. 本地消息表；
2. 定时任务补偿；
3. 事务消息；
4. Outbox Pattern。

当前项目作为练习项目，可以先简单处理：

- MQ 投递失败时更新任务状态为 `failed`；
- 返回 `CodeTaskSubmitFailed`。

## 迭代建议

建议按以下顺序继续推进：
1. 先完成注册、登录、JWT 中间件、统一响应和错误码；
2. 再完成 `tasks` 写库和状态查询；
3. 接着接入 Redis 状态缓存和限流；
4. 然后接入 MQ 投递；
5. 最后补消费端、AI 服务调用和失败补偿机制。
