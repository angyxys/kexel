import axiosInstance from './client';

export interface Webhook {
  id: number;
  name: string;
  url: string;
  events: string[];
  is_active: boolean;
  failure_count: number;
  last_tried_at: string | null;
  last_success_at: string | null;
  created_at: string;
}

export interface WebhookEvent {
  id: number;
  event_type: string;
  status_code: number;
  is_delivered: boolean;
  attempts: number;
  created_at: string;
}

export const webhooksApi = {
  createWebhook: async (name: string, url: string, events: string[]): Promise<Webhook> => {
    const response = await axiosInstance.post('/web/webhooks', { name, url, events });
    return response.data;
  },

  getWebhooks: async (): Promise<Webhook[]> => {
    const response = await axiosInstance.get('/web/webhooks');
    return response.data.data;
  },

  deleteWebhook: async (webhookId: number): Promise<{ message: string }> => {
    const response = await axiosInstance.delete(`/web/webhooks/${webhookId}`);
    return response.data;
  },

  disableWebhook: async (webhookId: number): Promise<{ message: string }> => {
    const response = await axiosInstance.post(`/web/webhooks/${webhookId}/disable`);
    return response.data;
  },

  getWebhookEvents: async (webhookId: number, page: number = 1): Promise<WebhookEvent[]> => {
    const response = await axiosInstance.get(`/web/webhooks/${webhookId}/events`, { params: { page } });
    return response.data.data;
  },

  getAvailableEvents: async (): Promise<string[]> => {
    const response = await axiosInstance.get('/web/webhooks/events/available');
    return response.data.events;
  },
};
