-- ERP 系统数据库初始化脚本 (PostgreSQL)

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(50),
    balance_amount DECIMAL(15, 2) DEFAULT 0.00,
    balance_currency VARCHAR(10) DEFAULT 'CNY',
    status_code VARCHAR(50) DEFAULT 'active',
    status_description VARCHAR(255) DEFAULT 'Active',
    status_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 产品表
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    sku VARCHAR(100) NOT NULL UNIQUE,
    price DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    cost DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    stock INTEGER NOT NULL DEFAULT 0,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 订单表
CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    order_no VARCHAR(100) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL REFERENCES users(id),
    total_amount DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    paid_at TIMESTAMP,
    shipped_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 订单项表
CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    subtotal DECIMAL(15, 2) NOT NULL DEFAULT 0.00
);

-- 发票表
CREATE TABLE IF NOT EXISTS invoices (
    id BIGSERIAL PRIMARY KEY,
    invoice_no VARCHAR(100) NOT NULL UNIQUE,
    order_id BIGINT REFERENCES orders(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    amount DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    tax_amount DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    total_amount DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    type VARCHAR(50) NOT NULL DEFAULT 'sales',
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    issued_at TIMESTAMP,
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 充值记录表
CREATE TABLE IF NOT EXISTS recharge_records (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    amount DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    type VARCHAR(50) NOT NULL DEFAULT 'user',
    method VARCHAR(50) NOT NULL DEFAULT 'bank_transfer',
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    remark TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 供应商表
CREATE TABLE IF NOT EXISTS suppliers (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    contact VARCHAR(100),
    phone VARCHAR(50),
    email VARCHAR(255),
    address TEXT,
    balance DECIMAL(15, 2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 采购订单表
CREATE TABLE IF NOT EXISTS purchase_orders (
    id BIGSERIAL PRIMARY KEY,
    order_no VARCHAR(100) NOT NULL UNIQUE,
    supplier_id BIGINT NOT NULL REFERENCES suppliers(id),
    total_amount DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    received_at TIMESTAMP,
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 采购订单项表
CREATE TABLE IF NOT EXISTS purchase_order_items (
    id BIGSERIAL PRIMARY KEY,
    purchase_order_id BIGINT NOT NULL REFERENCES purchase_orders(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_cost DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    subtotal DECIMAL(15, 2) NOT NULL DEFAULT 0.00
);

-- 消费账单表
CREATE TABLE IF NOT EXISTS consumption_bills (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    order_id BIGINT REFERENCES orders(id),
    amount DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'unpaid',
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_invoices_user_id ON invoices(user_id);
CREATE INDEX idx_recharge_records_user_id ON recharge_records(user_id);
CREATE INDEX idx_consumption_bills_user_id ON consumption_bills(user_id);

-- 插入示例数据
INSERT INTO users (username, email, phone, balance) VALUES 
('admin', 'admin@example.com', '13800138000', 0.00),
('user1', 'user1@example.com', '13800138001', 1000.00),
('user2', 'user2@example.com', '13800138002', 500.00);

INSERT INTO products (name, sku, price, cost, stock, description) VALUES 
('产品 A', 'SKU001', 99.00, 50.00, 100, '产品 A 描述'),
('产品 B', 'SKU002', 199.00, 100.00, 50, '产品 B 描述'),
('产品 C', 'SKU003', 299.00, 150.00, 30, '产品 C 描述');

INSERT INTO suppliers (name, contact, phone, email, address) VALUES 
('供应商 A', '张三', '13900139000', 'supplier_a@example.com', '北京市朝阳区'),
('供应商 B', '李四', '13900139001', 'supplier_b@example.com', '上海市浦东新区');
