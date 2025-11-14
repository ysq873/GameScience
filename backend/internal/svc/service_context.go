// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"backend/internal/config"
	"backend/internal/middleware"
	"backend/internal/store"

	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config         config.Config
	AuthMiddleware rest.Middleware
	DB             *store.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, _ := store.NewDB(c.MysqlDataSource)
	sc := &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewAuthMiddleware(c.KratosPublicURL).Handle,
		DB:             db,
	}
	// 初始化表结构（幂等）
	if sc.DB != nil {
		_ = sc.DB.Exec(`CREATE TABLE IF NOT EXISTS models (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  owner_id VARCHAR(64) NOT NULL,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  price_cents INT NOT NULL,
  cover_url VARCHAR(512),
  file_path VARCHAR(512) NOT NULL,
  status TINYINT NOT NULL DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)`)
		_ = sc.DB.Exec(`CREATE TABLE IF NOT EXISTS orders (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id VARCHAR(64) NOT NULL,
  total_cents INT NOT NULL,
  status TINYINT NOT NULL DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)`)
		_ = sc.DB.Exec(`CREATE TABLE IF NOT EXISTS order_items (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  order_id BIGINT NOT NULL,
  model_id BIGINT NOT NULL,
  price_cents INT NOT NULL,
  title_snapshot VARCHAR(255),
  FOREIGN KEY (order_id) REFERENCES orders(id)
)`)
		_ = sc.DB.Exec(`CREATE TABLE IF NOT EXISTS payments (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  order_id BIGINT NOT NULL,
  amount_cents INT NOT NULL,
  status ENUM('initiated','succeeded','failed') NOT NULL DEFAULT 'initiated',
  gateway VARCHAR(32) NOT NULL,
  idempotency_key VARCHAR(128) UNIQUE,
  callback_id VARCHAR(128),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (order_id) REFERENCES orders(id)
)`)
		_ = sc.DB.Exec(`CREATE TABLE IF NOT EXISTS purchases (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id VARCHAR(64) NOT NULL,
  model_id BIGINT NOT NULL,
  order_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uniq_user_model (user_id, model_id)
)`)
		_ = sc.DB.Exec(`CREATE TABLE IF NOT EXISTS download_tokens (
  token VARCHAR(64) PRIMARY KEY,
  user_id VARCHAR(64) NOT NULL,
  model_id BIGINT NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  one_time TINYINT(1) NOT NULL DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)`)
	}
	return sc
}
