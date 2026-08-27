import axiosInstance from './client';

export interface InvitationCode {
  id: number;
  code: string;
  role: string;
  max_uses: number;
  uses: number;
  expires_at: string | null;
  is_active: boolean;
  created_at: string;
}

export interface CreateInvitationRequest {
  role: string;
  max_uses?: number;
  expires_at?: string;
}

export interface InvitationStats {
  total_active: number;
  total_uses: number;
  total_capacity: number;
}

export const invitationsApi = {
  createInvitation: async (request: CreateInvitationRequest): Promise<InvitationCode> => {
    const response = await axiosInstance.post('/web/invitations', request);
    return response.data;
  },

  getMyInvitations: async (): Promise<InvitationCode[]> => {
    const response = await axiosInstance.get('/web/invitations');
    return response.data;
  },

  validateInvitation: async (code: string): Promise<{ valid: boolean; role: string; uses: number; max_uses: number }> => {
    const response = await axiosInstance.get('/auth/invitations/validate', {
      params: { code },
    });
    return response.data;
  },

  revokeInvitation: async (id: number): Promise<{ message: string }> => {
    const response = await axiosInstance.delete(`/web/invitations/${id}`);
    return response.data;
  },

  getInvitationStats: async (): Promise<InvitationStats> => {
    const response = await axiosInstance.get('/web/invitations/stats');
    return response.data;
  },
};
