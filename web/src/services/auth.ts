import { request } from '@/utils/request';

export interface LoginRequest {
  username: string;
  password: string;
  captcha: string;
  captcha_id: string;
}

export interface LoginResponse {
  token: string;
  user: {
    id: number;
    username: string;
    email: string;
    role: string;
    status: string;
    createdAt: string;
  };
}

export const authService = {
  // 获取验证码
  getCaptcha: async () => {
    const response = await request.get('/auth/captcha');
    return response.data;
  },

  // 用户登录
  login: async (data: LoginRequest) => {
    const response = await request.post<LoginResponse>('/auth/login', data);
    return response.data;
  },
};