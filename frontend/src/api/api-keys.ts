import axiosInstance from './client';

export interface APIKey {
  id: number;
  name: string;
  key_prefix: string;
  scopes: string[];
  is_active: boolean;
  last_used_at: string | null;
  created_at: string;
  expires_at: string | null;
  rate_limit: number;
}

export interface CreateAPIKeyRequest {
  name: string;
  scopes: string[];
  expires_in?: number; // Days
  rate_limit?: number;
}

export interface CreateAPIKeyResponse extends APIKey {
  key: string; // Only in creation response
}

export const apiKeysApi = {
  createAPIKey: async (request: CreateAPIKeyRequest): Promise<CreateAPIKeyResponse> => {
    const response = await axiosInstance.post('/web/api-keys', request);
    return response.data;
  },

  getAPIKeys: async (): Promise<APIKey[]> => {
    const response = await axiosInstance.get('/web/api-keys');
    return response.data.data;
  },

  deleteAPIKey: async (keyId: number): Promise<{ message: string }> => {
    const response = await axiosInstance.delete(`/web/api-keys/${keyId}`);
    return response.data;
  },

  revokeAPIKey: async (keyId: number): Promise<{ message: string }> => {
    const response = await axiosInstance.post(`/web/api-keys/${keyId}/revoke`);
    return response.data;
  },

  updateRateLimit: async (keyId: number, rateLimit: number): Promise<{ message: string }> => {
    const response = await axiosInstance.patch(`/web/api-keys/${keyId}/rate-limit`, {
      rate_limit: rateLimit,
    });
    return response.data;
  },

  updateScopes: async (keyId: number, scopes: string[]): Promise<{ message: string }> => {
    const response = await axiosInstance.patch(`/web/api-keys/${keyId}/scopes`, {
      scopes,
    });
    return response.data;
  },

  getAvailableScopes: async (): Promise<string[]> => {
    const response = await axiosInstance.get('/web/api-keys/scopes/available');
    return response.data.scopes;
  },
};
