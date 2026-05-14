  -- 初始化 MySQL 表结构脚本

  USE ai_task;

  CREATE TABLE IF NOT EXISTS users (
      id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
      username VARCHAR(64) NOT NULL,
      password_hash VARCHAR(255) NOT NULL,
      created_at DATETIME(3) NULL,
      updated_at DATETIME(3) NULL,
      PRIMARY KEY (id),
      UNIQUE KEY idx_users_username (username)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

  CREATE TABLE IF NOT EXISTS tasks (
      id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
      task_id VARCHAR(64) NOT NULL,
      user_id BIGINT UNSIGNED NOT NULL,
      prompt TEXT NOT NULL,
      result LONGTEXT NULL,
      status VARCHAR(16) NOT NULL DEFAULT 'pending',
      created_at DATETIME(3) NULL,
      updated_at DATETIME(3) NULL,
      PRIMARY KEY (id),
      UNIQUE KEY idx_tasks_task_id (task_id),
      KEY idx_tasks_user_id (user_id),
      KEY idx_tasks_status (status)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;