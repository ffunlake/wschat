#!/bin/bash

# 定义颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 定义SQLite数据库文件路径
DB_FILE="../data/wschat.db"
DB_DIR="../data"

# 创建数据目录（如果不存在）
echo -e "${GREEN}Creating data directory...${NC}"
mkdir -p $DB_DIR

# 检查SQLite3是否安装
if ! command -v sqlite3 &> /dev/null; then
    echo -e "${RED}Error: sqlite3 is not installed. Please install it first.${NC}"
    exit 1
fi

# 检查数据库文件是否已存在
if [ -f "$DB_FILE" ]; then
    echo -e "${GREEN}Database file already exists. Backing up...${NC}"
    cp "$DB_FILE" "${DB_FILE}.backup.$(date +%Y%m%d%H%M%S)"
fi

# 创建customers表
echo -e "${GREEN}Creating database schema...${NC}"
sqlite3 $DB_FILE <<EOF
-- 创建customers表
CREATE TABLE IF NOT EXISTS customers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    customer_id VARCHAR(100) NOT NULL UNIQUE,
    customer_name VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建messages表
CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    customer_id VARCHAR(100) NOT NULL,
    message TEXT NOT NULL,
    sender VARCHAR(20) NOT NULL, -- 'customer' or 'bot'
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建feedback表
CREATE TABLE IF NOT EXISTS feedback (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    customer_id VARCHAR(100) NOT NULL,
    rating INTEGER NOT NULL, -- 1-5 rating
    comment TEXT,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建触发器以更新updated_at
CREATE TRIGGER IF NOT EXISTS update_customers_timestamp 
AFTER UPDATE ON customers
BEGIN
    UPDATE customers SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
EOF

# 检查表是否创建成功
TABLES=$(sqlite3 $DB_FILE "SELECT name FROM sqlite_master WHERE type='table';")
if [[ $TABLES == *"customers"* ]] && [[ $TABLES == *"messages"* ]] && [[ $TABLES == *"feedback"* ]]; then
    echo -e "${GREEN}Database schema created successfully.${NC}"
    echo -e "${GREEN}Created tables: ${NC}"
    echo "$TABLES" | grep -v "sqlite_"
else
    echo -e "${RED}Failed to create all tables.${NC}"
    echo -e "${GREEN}Tables created: ${NC}"
    echo "$TABLES" | grep -v "sqlite_"
    exit 1
fi

# 添加执行权限
chmod +x "$0"

echo -e "${GREEN}Installation completed successfully!${NC}"
