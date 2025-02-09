import axios, { AxiosRequestConfig, AxiosResponse } from 'axios';
import { handleApiError } from '@/utils/error';
import { logger } from '@/utils/logger';
import { API_BASE_URL, AUTH_TOKEN_KEY, AUTH_REMEMBER_ME_KEY } from '@/utils/config';

const baseURL = API_BASE_URL;

const axiosInstance = axios.create({
  baseURL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
axiosInstance.interceptors.request.use(
  (config) => {
    // 添加 token，从 localStorage 或 sessionStorage 中获取
    const token = localStorage.getItem(AUTH_TOKEN_KEY) || sessionStorage.getItem(AUTH_TOKEN_KEY);
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    
    // 添加时间戳防止缓存
    if (config.method === 'get') {
      config.params = {
        ...config.params,
        _t: Date.now(),
      };
    }
    
    // 记录请求日志
    logger.logRequest(config);
    return config;
  },
  (error) => {
    logger.logError(error);
    return Promise.reject(error);
  }
);

// 响应拦截器
axiosInstance.interceptors.response.use(
  (response) => {
    // 如果响应中包含新的 token，更新本地存储
    const newToken = response.headers['x-new-token'];
    if (newToken) {
      // 根据用户的选择存储到 localStorage 或 sessionStorage
      const rememberMe = localStorage.getItem(AUTH_REMEMBER_ME_KEY) === 'true';
      if (rememberMe) {
        localStorage.setItem(AUTH_TOKEN_KEY, newToken);
      } else {
        sessionStorage.setItem(AUTH_TOKEN_KEY, newToken);
      }
    }
    // 记录响应日志
    logger.logResponse(response);
    return response;
  },
  (error) => {
    // 处理 401 未授权错误
    if (error.response && error.response.status === 401) {
      // 清除本地存储的 token
      localStorage.removeItem(AUTH_TOKEN_KEY);
      sessionStorage.removeItem(AUTH_TOKEN_KEY);
      // 重定向到登录页面
      window.location.href = '/login';
    }
    logger.logError(error);
    return handleApiError(error);
  }
);

export const request = {
  get: <T = any>(url: string, config?: AxiosRequestConfig) => 
    axiosInstance.get<any, T>(url, config),

  post: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) => 
    axiosInstance.post<any, T>(url, data, config),

  put: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) => 
    axiosInstance.put<any, T>(url, data, config),

  delete: <T = any>(url: string, config?: AxiosRequestConfig) => 
    axiosInstance.delete<any, T>(url, config),

  patch: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) => 
    axiosInstance.patch<any, T>(url, data, config),
};

export default axiosInstance;