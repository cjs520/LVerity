import { User, UserCreateRequest, UserUpdateRequest } from '../models/user';
import api from './request';

export const userService = {
  // 获取用户列表
  list: async (params: { current: number; pageSize: number; [key: string]: any }) => {
    const { data } = await api.get('/users', { params });
    return data;
  },

  // 获取单个用户
  get: async (id: number) => {
    const { data } = await api.get(`/users/${id}`);
    return data as User;
  },

  // 创建用户
  create: async (user: UserCreateRequest) => {
    const { data } = await api.post('/users', user);
    return data as User;
  },

  // 更新用户
  update: async (id: number, user: UserUpdateRequest) => {
    const { data } = await api.put(`/users/${id}`, user);
    return data as User;
  },

  // 删除用户
  delete: async (id: number) => {
    await api.delete(`/users/${id}`);
  },
};
