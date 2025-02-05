import axios, { AxiosRequestConfig, AxiosResponse } from 'axios';
import { handleApiError } from './error';
import { logger } from './logger';

const baseURL = import.meta.env.VITE_API_BASE_URL || '';

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
    // 从localStorage获取token
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
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
    // 记录响应日志
    logger.logResponse(response);
    // 直接返回响应数据，不需要再次访问 .data
    return response;
  },
  (error) => {
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