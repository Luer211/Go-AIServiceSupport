# 数据库表、索引、约束

系统包含两张核心数据表：

1. 用户表 `users`；
2. 任务表 `tasks`。

## users 表

### 建表 SQL

```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(64) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
```

## tasks 表

### 建表 SQL

```sql
CREATE TABLE tasks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    task_id VARCHAR(64) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    prompt TEXT NOT NULL,
    result LONGTEXT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'pending',
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
```

### 索引设计

```sql
CREATE INDEX idx_user_id ON tasks(user_id);

-- UNIQUE 约束会自动创建索引，因此 task_id 不需要额外创建普通索引
-- CREATE INDEX idx_task_id ON tasks(task_id);

CREATE INDEX idx_status ON tasks(status);
```

说明：`task_id` 已经有唯一约束，MySQL 会自动为唯一约束创建索引，因此无需重复创建 `idx_task_id`。

## 任务状态枚举

```go
const (
    TaskStatusPending = "pending"
    TaskStatusSuccess = "success"
    TaskStatusFailed  = "failed"
)
```

任务状态说明：
- `pending`：任务已提交，但尚未完成
- `success`：任务执行成功
- `failed`：任务执行失败
