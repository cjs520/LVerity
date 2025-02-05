SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;
USE lverity;

-- 插入默认角色
INSERT INTO roles (id, name, type, description) VALUES
('1', 'admin', 'system', '系统管理员'),
('2', 'operator', 'system', '运营人员'),
('3', 'viewer', 'system', '查看者');

-- 插入默认权限
INSERT INTO permissions (id, resource, action) VALUES
('1', 'system', 'admin'),
('2', 'license', 'manage'),
('3', 'device', 'manage'),
('4', 'user', 'manage'),
('5', 'alert', 'manage'),
('6', 'log', 'view');

-- 插入角色权限关联
INSERT INTO role_permissions (role_id, permission_id) VALUES
-- 管理员拥有所有权限
('1', '1'), ('1', '2'), ('1', '3'), ('1', '4'), ('1', '5'), ('1', '6'),
-- 运营人员拥有部分权限
('2', '2'), ('2', '3'), ('2', '5'),
-- 查看者只有查看权限
('3', '6');

-- 插入默认管理员用户
-- 密码: admin123 (使用 bcrypt 加密)
INSERT INTO users (id, username, password, role_id, status, salt) VALUES 
('1', 'admin', '$2a$10$IVxZm.OZb4GHxVDyHPHzXuQ1Zw0tK0C8C9uTkR6Kb3VE3Eh3Z5K6W', '1', 'active', '123456');
