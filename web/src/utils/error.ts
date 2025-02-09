import { message } from 'antd';

interface ApiError {
  code?: string;
  message?: string;
  details?: string;
}

export function handleApiError(error: any) {
  const apiError: ApiError = error.response?.data || {};
  
  // 处理常见错误码
  switch (error.response?.status) {
    case 400:
      message.error(apiError.message || '请求参数错误');
      break;
    case 401:
      message.error('登录已过期，请重新登录');
      // 跳转到登录页
      window.location.href = '/login';
      break;
    case 403:
      message.error('没有权限执行此操作');
      break;
    case 404:
      message.error('请求的资源不存在');
      break;
    case 500:
      message.error('服务器错误，请稍后重试');
      break;
    default:
      message.error(apiError.message || '操作失败，请重试');
  }
  
  // 仅在开发环境下输出错误日志
  if (import.meta.env.MODE === 'development' || import.meta.env.VITE_LOG_LEVEL === 'debug') {
    console.error('[API Error]', {
      status: error.response?.status,
      url: error.config?.url,
      error: apiError,
    });
  }

  return Promise.reject(error);
}
