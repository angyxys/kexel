import axiosInstance from './client';

export interface KPIStats {
  total_players: number;
  total_vips: number;
  total_mods: number;
  total_owners: number;
  total_banned: number;
  total_active: number;
  total_users: number;
  today_logins: number;
  bans_this_week: number;
  unbans_this_week: number;
}

export interface ActivityLog {
  id: number;
  user: string;
  action: string;
  resource: string;
  timestamp: string;
  details: string;
}

export interface TrendData {
  date: string;
  value: number;
}

export interface RoleDistribution {
  role: string;
  count: number;
}

export interface BanStats {
  total_banned: number;
  permanent_bans: number;
  temporary_bans: number;
  expiring_today: number;
  expiring_week: number;
  most_common_reason: string;
}

export interface DashboardOverview {
  kpi: KPIStats;
  recent_activities: ActivityLog[];
  role_distribution: RoleDistribution[];
  ban_stats: BanStats;
  player_trends: TrendData[];
  timestamp: string;
}

export const statsApi = {
  getKPIStats: async (): Promise<KPIStats> => {
    const response = await axiosInstance.get('/web/stats/kpi');
    return response.data;
  },

  getRecentActivity: async (limit = 20): Promise<ActivityLog[]> => {
    const response = await axiosInstance.get('/web/stats/activity', {
      params: { limit },
    });
    return response.data;
  },

  getPlayerTrends: async (days = 30): Promise<TrendData[]> => {
    const response = await axiosInstance.get('/web/stats/trends', {
      params: { days },
    });
    return response.data;
  },

  getRoleDistribution: async (): Promise<RoleDistribution[]> => {
    const response = await axiosInstance.get('/web/stats/roles');
    return response.data;
  },

  getBanStats: async (): Promise<BanStats> => {
    const response = await axiosInstance.get('/web/stats/bans');
    return response.data;
  },

  getDashboardOverview: async (): Promise<DashboardOverview> => {
    const response = await axiosInstance.get('/web/stats/overview');
    return response.data;
  },
};
