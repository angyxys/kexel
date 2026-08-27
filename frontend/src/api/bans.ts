import axiosInstance from './client';

export interface BanInfo {
  vrchat_id: string;
  is_banned: boolean;
  reason: string;
  banned_at: string | null;
  expires_at: string | null;
  is_expired: boolean;
  time_left: string;
  is_permanent: boolean;
}

export interface BanRequest {
  reason: string;
  duration?: number; // Duration in hours (0 = permanent)
  expires_at?: string; // ISO 8601 format
}

export const bansApi = {
  banPlayer: async (playerID: string, request: BanRequest): Promise<{ message: string }> => {
    const response = await axiosInstance.post(`/web/ban/${playerID}`, request);
    return response.data;
  },

  unbanPlayer: async (playerID: string): Promise<{ message: string }> => {
    const response = await axiosInstance.delete(`/web/ban/${playerID}`);
    return response.data;
  },

  getBanInfo: async (playerID: string): Promise<BanInfo> => {
    const response = await axiosInstance.get(`/web/ban/${playerID}`);
    return response.data;
  },

  getBannedPlayers: async (): Promise<BanInfo[]> => {
    const response = await axiosInstance.get('/web/bans');
    return response.data;
  },

  getExpiringSoonBans: async (): Promise<BanInfo[]> => {
    const response = await axiosInstance.get('/web/bans/expiring-soon');
    return response.data;
  },

  cleanupExpiredBans: async (): Promise<{ message: string; unbanned_count: number }> => {
    const response = await axiosInstance.post('/web/bans/cleanup');
    return response.data;
  },
};
