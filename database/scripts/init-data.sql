SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;
USE lverity;

-- 插入默认角色
INSERT INTO roles (id, name, description, status) VALUES
('1', 'admin', '系统管理员', 'active'),
('2', 'operator', '运营人员', 'active'),
('3', 'viewer', '查看者', 'active');

-- 插入默认权限
INSERT INTO permissions (id, code, name, description) VALUES
('1', 'system:admin', '系统管理', '系统管理权限'),
('2', 'license:manage', '许可证管理', '许可证管理权限'),
('3', 'device:manage', '设备管理', '设备管理权限'),
('4', 'user:manage', '用户管理', '用户管理权限'),
('5', 'alert:manage', '告警管理', '告警管理权限'),
('6', 'log:view', '日志查看', '日志查看权限');

-- 插入角色权限关联
INSERT INTO role_permissions (role_id, permission_id) VALUES
-- 管理员拥有所有权限
('1', '1'), ('1', '2'), ('1', '3'), ('1', '4'), ('1', '5'), ('1', '6'),
-- 运营人员拥有部分权限
('2', '2'), ('2', '3'), ('2', '5'),
-- 查看者只有查看权限
('3', '6');

-- 插入默认管理员用户
-- 密码: admin123 (使用 SHA256 加密)
INSERT INTO users (id, username, password, role_id, status, salt) VALUES 
('1', 'admin', '240be518fabd2724ddb6f04eeb1da5967448d7e831c08c8fa822809f74c720a9', '1', 'active', '123456');
