import axiosInstance from './client';
import { Player } from '../types';

export const playersApi = {
  listPlayers: async (): Promise<Player[]> => {
    const response = await axiosInstance.get('/web/players');
    return response.data;
  },

  getPlayer: async (id: string): Promise<Player> => {
    const response = await axiosInstance.get(`/web/player/${id}`);
    return response.data;
  },

  createPlayer: async (data: {
    vrchat_id: string;
    roles?: ('user' | 'vip' | 'mod' | 'owner')[];
    is_banned?: boolean;
  }): Promise<{ player: Player }> => {
    const response = await axiosInstance.post('/web/player', data);
    return response.data;
  },

  updatePlayer: async (
    id: string,
    data: {
      roles?: ('user' | 'vip' | 'mod' | 'owner')[];
      is_banned?: boolean;
    }
  ): Promise<{ player: Player }> => {
    const response = await axiosInstance.put(`/web/player/${id}`, data);
    return response.data;
  },

  deletePlayer: async (id: string): Promise<void> => {
    await axiosInstance.delete(`/web/player/${id}`);
  },
};
