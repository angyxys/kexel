import axiosInstance from './client';

export interface TOTPSetup {
  secret: string;
  qr_code: string;
  backup_codes: string[];
}

export interface TOTPStatus {
  is_enabled: boolean;
  backup_codes_left: number;
  enabled_at: string | null;
  last_used_at: string | null;
}

export const twoFactorApi = {
  setupTOTP: async (email: string): Promise<TOTPSetup> => {
    const response = await axiosInstance.post('/web/2fa/setup', { email });
    return response.data;
  },

  verifyTOTP: async (code: string): Promise<{ message: string }> => {
    const response = await axiosInstance.post('/web/2fa/verify', { code });
    return response.data;
  },

  getTOTPStatus: async (): Promise<TOTPStatus> => {
    const response = await axiosInstance.get('/web/2fa/status');
    return response.data;
  },

  disableTOTP: async (): Promise<{ message: string }> => {
    const response = await axiosInstance.delete('/web/2fa');
    return response.data;
  },

  verifyTOTPCode: async (code: string): Promise<{ message: string }> => {
    const response = await axiosInstance.post('/web/2fa/verify-code', { code });
    return response.data;
  },

  verifyBackupCode: async (code: string): Promise<{ message: string }> => {
    const response = await axiosInstance.post('/web/2fa/backup-code', { code });
    return response.data;
  },
};
