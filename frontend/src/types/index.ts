export interface User {
  id: number;
  username: string;
  email: string;
  role: 'user' | 'vip' | 'mod' | 'owner';
}

export interface Player {
  vrchat_id: string;
  roles: ('user' | 'vip' | 'mod' | 'owner')[];
  is_banned: boolean;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  token_type: string;
}

export interface ApiError {
  message: string;
  status: number;
}
