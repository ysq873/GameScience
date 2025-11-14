-- 模型（商品）表
CREATE TABLE IF NOT EXISTS models (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  owner_id VARCHAR(64) NOT NULL COMMENT '拥有者/作者的用户ID',
  title VARCHAR(255) NOT NULL COMMENT '模型标题（展示名）',
  description TEXT COMMENT '模型描述',
  price_cents INT NOT NULL COMMENT '价格（单位：分）',
  cover_url VARCHAR(512) COMMENT '封面图URL',
  file_path VARCHAR(512) NOT NULL COMMENT '存储文件路径（用于下载）',
  status TINYINT NOT NULL DEFAULT 0 COMMENT '状态：0=草稿/待审核，1=上架，-1=下架（可自定义）',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) COMMENT='模型（商品）表';

-- 订单表
CREATE TABLE IF NOT EXISTS orders (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  user_id VARCHAR(64) NOT NULL COMMENT '下单用户ID',
  total_cents INT NOT NULL COMMENT '订单总额（分）',
  status TINYINT NOT NULL DEFAULT 0 COMMENT '订单状态',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) COMMENT='订单表';

-- 订单明细表
CREATE TABLE IF NOT EXISTS order_items (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  order_id BIGINT NOT NULL COMMENT '所属订单ID（外键）',
  model_id BIGINT NOT NULL COMMENT '购买的模型ID',
  price_cents INT NOT NULL COMMENT '该条目结算价格（分，快照）',
  title_snapshot VARCHAR(255) COMMENT '购买时的标题快照（避免后续标题变更影响历史）',
  FOREIGN KEY (order_id) REFERENCES orders(id)
) COMMENT='订单明细表';

-- 支付表
CREATE TABLE IF NOT EXISTS payments (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  order_id BIGINT NOT NULL COMMENT '关联订单ID（外键）',
  amount_cents INT NOT NULL COMMENT '支付金额（分）',
  status ENUM('initiated','succeeded','failed') NOT NULL DEFAULT 'initiated' COMMENT '支付状态',
  gateway VARCHAR(32) NOT NULL COMMENT '支付渠道（如 stripe/wechatpay/alipay 等）',
  idempotency_key VARCHAR(128) UNIQUE COMMENT '幂等键（防重复提交/回调）',
  callback_id VARCHAR(128) COMMENT '第三方支付回调/交易ID',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  FOREIGN KEY (order_id) REFERENCES orders(id)
) COMMENT='支付流水表';

-- 购买记录表（授权用）
CREATE TABLE IF NOT EXISTS purchases (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  user_id VARCHAR(64) NOT NULL COMMENT '用户ID',
  model_id BIGINT NOT NULL COMMENT '模型ID',
  order_id BIGINT NOT NULL COMMENT '对应订单ID（便于追溯）',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '购买时间',
  UNIQUE KEY uniq_user_model (user_id, model_id)
) COMMENT='购买记录表：用户与模型的已购关系';

-- 下载令牌表
CREATE TABLE IF NOT EXISTS download_tokens (
  token VARCHAR(64) PRIMARY KEY COMMENT '下载令牌（随机字符串）',
  user_id VARCHAR(64) NOT NULL COMMENT '令牌所属用户ID（授权校验）',
  model_id BIGINT NOT NULL COMMENT '目标模型ID',
  expires_at TIMESTAMP NOT NULL COMMENT '过期时间',
  one_time TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否一次性：1=一次性，用后失效',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '令牌创建时间'
) COMMENT='用于短期/一次性下载的令牌';
