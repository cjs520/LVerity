#!/bin/bash

# 设置变量
BACKUP_DIR="/backup/mysql"
DATE=$(date +%Y%m%d_%H%M%S)
DB_USER=${MYSQL_USER:-lverity}
DB_PASS=${MYSQL_PASSWORD:-lverity123}
DB_NAME=${MYSQL_DATABASE:-lverity}

# 创建备份目录
mkdir -p $BACKUP_DIR

# 执行备份
mysqldump -h localhost -u $DB_USER -p$DB_PASS $DB_NAME > "$BACKUP_DIR/${DB_NAME}_${DATE}.sql"

# 压缩备份文件
gzip "$BACKUP_DIR/${DB_NAME}_${DATE}.sql"

# 删除30天前的备份
find $BACKUP_DIR -name "*.sql.gz" -mtime +30 -delete

echo "数据库备份完成：$BACKUP_DIR/${DB_NAME}_${DATE}.sql.gz"