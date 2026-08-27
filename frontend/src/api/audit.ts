import axiosInstance from './client';

export interface AuditLog {
  id: number;
  user_id: number;
  username: string;
  action: string;
  resource_type: string;
  resource_id: string;
  description: string;
  ip_address: string;
  created_at: string;
}

export interface AuditListResponse {
  data: AuditLog[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface AuditStats {
  total_logs: number;
  today_logs: number;
  unique_users: number;
}

export const auditApi = {
  listLogs: async (
    page = 1,
    pageSize = 20,
    filters?: {
      user_id?: string;
      action?: string;
      resource_type?: string;
      resource_id?: string;
      start_date?: string;
      end_date?: string;
    }
  ): Promise<AuditListResponse> => {
    const params = new URLSearchParams();
    params.append('page', page.toString());
    params.append('page_size', pageSize.toString());

    if (filters) {
      if (filters.user_id) params.append('user_id', filters.user_id);
      if (filters.action) params.append('action', filters.action);
      if (filters.resource_type) params.append('resource_type', filters.resource_type);
      if (filters.resource_id) params.append('resource_id', filters.resource_id);
      if (filters.start_date) params.append('start_date', filters.start_date);
      if (filters.end_date) params.append('end_date', filters.end_date);
    }

    const response = await axiosInstance.get(`/web/audit-logs?${params.toString()}`);
    return response.data;
  },

  getStats: async (): Promise<AuditStats> => {
    const response = await axiosInstance.get('/web/audit-logs/stats');
    return response.data;
  },

  exportLogs: async (filters?: {
    action?: string;
    resource_type?: string;
    start_date?: string;
    end_date?: string;
  }): Promise<Blob> => {
    const params = new URLSearchParams();
    if (filters) {
      if (filters.action) params.append('action', filters.action);
      if (filters.resource_type) params.append('resource_type', filters.resource_type);
      if (filters.start_date) params.append('start_date', filters.start_date);
      if (filters.end_date) params.append('end_date', filters.end_date);
    }

    const response = await axiosInstance.get(`/web/audit-logs/export?${params.toString()}`, {
      responseType: 'blob',
    });
    return response.data;
  },

  getResourceHistory: async (resourceId: string): Promise<AuditLog[]> => {
    const response = await axiosInstance.get(`/web/audit-logs/resource/${resourceId}`);
    return response.data;
  },
};
