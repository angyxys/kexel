import axiosInstance from './client';

export interface SessionInfo {
  id: number;
  device_name: string;
  ip_address: string;
  is_active: boolean;
  last_activity: string;
  login_at: string;
  logout_at: string | null;
  expires_at: string;
}

export interface SessionStats {
  total_active_sessions: number;
  total_sessions: number;
  current_device: string;
  current_ip: string;
}

export const sessionsApi = {
  getSessions: async (): Promise<SessionInfo[]> => {
    const response = await axiosInstance.get('/web/sessions');
    return response.data.data;
  },

  getSessionStats: async (): Promise<SessionStats> => {
    const response = await axiosInstance.get('/web/sessions/stats');
    return response.data;
  },

  logoutSession: async (sessionId: number): Promise<{ message: string }> => {
    const response = await axiosInstance.delete(`/web/sessions/${sessionId}`);
    return response.data;
  },

  logoutAllOtherSessions: async (): Promise<{ message: string }> => {
    const response = await axiosInstance.post('/web/sessions/logout-all');
    return response.data;
  },
};
