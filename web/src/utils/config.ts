// 环境变量配置管理

// API配置
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api';

// 认证配置
export const AUTH_TOKEN_KEY = import.meta.env.VITE_AUTH_TOKEN_KEY || 'token';
export const AUTH_REMEMBER_ME_KEY = import.meta.env.VITE_AUTH_REMEMBER_ME_KEY || 'rememberMe';

// 日志配置
export const LOG_LEVEL = import.meta.env.VITE_LOG_LEVEL || 'info';

// 系统配置
export const APP_TITLE = import.meta.env.VITE_APP_TITLE || 'LVerity';
export const APP_DESCRIPTION = import.meta.env.VITE_APP_DESCRIPTION || '设备验证系统';

// 开发环境标识
export const IS_DEV = import.meta.env.MODE === 'development';

// 获取配置值的工具函数
export function getConfig(key: string, defaultValue: string = ''): string {
  return import.meta.env[`VITE_${key}`] || defaultValue;
}