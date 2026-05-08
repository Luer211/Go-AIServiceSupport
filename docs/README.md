# 文档总览

## 项目概述

本项目实现一个基于 **JWT 鉴权 + MySQL 持久化 + Redis 状态缓存与限流 + MQ 异步消费** 的 AI 任务提交与查询系统。

系统核心目标：

- 用户可以注册、登录并获取 JWT；
- 用户可以提交 AI 任务；
- 后端将任务写入 MySQL，并投递到 MQ；
- 消费端异步消费任务并调用 AI 服务；
- 前端通过轮询接口查询任务是否完成；
- Redis 用于任务完成状态缓存与请求限流；
- Zap 用于结构化日志；
- Viper 用于配置管理；
- 接口统一使用响应结构和错误码规范。

## MVP 范围

MVP 主要包含三条业务链路：

1. 用户发起 AI 任务，并通过前端轮询查询任务是否完成；
2. 用户注册；
3. 用户登录。

MVP 主要包含四个接口：

1. 注册接口；
2. 登录接口；
3. 提交 AI 任务接口；
4. 查询任务状态接口。

MVP 主要包含两张数据表：

1. 用户表 `users`；
2. 任务表 `tasks`。

## 推荐阅读顺序

1. 先阅读 [architecture.md](./architecture.md)，理解整体结构和组件职责；
2. 再阅读 [business-flows.md](./business-flows.md)，理解核心业务链路；
3. 接着阅读 [database-schema.md](./database-schema.md) 和 [api-design.md](./api-design.md)，确认数据模型与接口契约；
4. 最后阅读 [roadmap-and-tech-debt.md](./roadmap-and-tech-debt.md)，确认 MVP 边界和后续迭代方向。
