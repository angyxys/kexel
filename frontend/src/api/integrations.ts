import axiosInstance from './client';

export interface DiscordIntegration {
  id: number;
  guild_id: string;
  mod_log_channel_id: string;
  announcement_channel_id: string;
  is_connected: boolean;
}

export interface RateLimitBlock {
  ip_address: string;
  endpoint: string;
  request_count: number;
  blocked_at: string;
}

export const integrationsApi = {
  // Discord
  setupDiscord: async (botToken: string, guildId: string): Promise<DiscordIntegration> => {
    const response = await axiosInstance.post('/web/integrations/discord', { bot_token: botToken, guild_id: guildId });
    return response.data;
  },

  getDiscordIntegration: async (): Promise<DiscordIntegration> => {
    const response = await axiosInstance.get('/web/integrations/discord');
    return response.data;
  },

  configureDiscordChannels: async (modLogChannelId: string, announcementChannelId: string): Promise<{ message: string }> => {
    const response = await axiosInstance.patch('/web/integrations/discord/channels', {
      mod_log_channel_id: modLogChannelId,
      announcement_channel_id: announcementChannelId,
    });
    return response.data;
  },

  testDiscordConnection: async (): Promise<{ message: string }> => {
    const response = await axiosInstance.post('/web/integrations/discord/test');
    return response.data;
  },

  disconnectDiscord: async (): Promise<{ message: string }> => {
    const response = await axiosInstance.delete('/web/integrations/discord');
    return response.data;
  },

  // Rate Limiting
  getRateLimitBlocks: async (): Promise<RateLimitBlock[]> => {
    const response = await axiosInstance.get('/admin/rate-limit/blocks');
    return response.data.data;
  },
};
