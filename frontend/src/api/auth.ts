import axiosInstance from './client';
import { AuthResponse, User } from '../types';

export const authApi = {
  register: async (
    data:
      | { username: string; email: string; password: string; invitation_code?: string }
      | (string | undefined)
  ): Promise<{ user: User }> => {
    // Support both old style (3 params) and new style (object)
    let payload;
    if (typeof data === 'string') {
      // Old style: register(username, email, password)
      // This branch should not happen with the new implementation
      payload = { username: data, email: '', password: '' };
    } else {
      payload = data;
    }

    const response = await axiosInstance.post('/auth/register', payload);
    return response.data;
  },

  login: async (username: string, password: string): Promise<AuthResponse> => {
    const response = await axiosInstance.post('/auth/login', {
      username,
      password,
    });
    return response.data;
  },

  refresh: async (refreshToken: string): Promise<AuthResponse> => {
    const response = await axiosInstance.post('/auth/refresh', {
      refresh_token: refreshToken,
    });
    return response.data;
  },

  logout: async (refreshToken: string): Promise<void> => {
    await axiosInstance.post('/auth/logout', {
      refresh_token: refreshToken,
    });
  },

  getCurrentUser: async (): Promise<{ user_id: number }> => {
    const response = await axiosInstance.get('/web/me');
    return response.data;
  },
};
