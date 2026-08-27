import { useEffect, useState } from 'react';
import { statsApi, DashboardOverview } from '../api/stats';
import { StatsCard } from '../components/StatsCard';
import { RecentActivity } from '../components/RecentActivity';
import { Navigation } from '../components/Navigation';

export function Stats() {
  const [overview, setOverview] = useState<DashboardOverview | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    loadDashboard();
    // Refresh every 30 seconds
    const interval = setInterval(loadDashboard, 30000);
    return () => clearInterval(interval);
  }, []);

  const loadDashboard = async () => {
    try {
      setLoading(true);
      setError('');
      const data = await statsApi.getDashboardOverview();
      setOverview(data);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load dashboard');
    } finally {
      setLoading(false);
    }
  };

  if (!overview) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800 flex items-center justify-center">
        <p className="text-slate-300">{loading ? 'Loading...' : error}</p>
      </div>
    );
  }

  const kpi = overview.kpi;
  const stats = overview.ban_stats;

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800">
      <Navigation />

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {error && (
          <div className="mb-6 p-4 bg-red-900/30 border border-red-700 rounded-lg text-red-200 text-sm">
            {error}
          </div>
        )}

        {/* KPI Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-8">
          <StatsCard
            title="Total Players"
            value={kpi.total_players}
            icon="👥"
            color="blue"
          />
          <StatsCard
            title="Active Players"
            value={kpi.total_active}
            icon="✅"
            color="green"
          />
          <StatsCard
            title="Banned"
            value={kpi.total_banned}
            icon="🚫"
            color="red"
          />
          <StatsCard
            title="VIPs"
            value={kpi.total_vips}
            icon="⭐"
            color="yellow"
          />
          <StatsCard
            title="Moderators"
            value={kpi.total_mods}
            icon="🛡️"
            color="purple"
          />
          <StatsCard
            title="Owners"
            value={kpi.total_owners}
            icon="👑"
            color="blue"
          />
        </div>

        {/* Ban Statistics */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-8">
          <div className="bg-slate-800/50 rounded-lg p-6 border border-slate-700">
            <h3 className="text-lg font-semibold text-white mb-4">Ban Statistics</h3>
            <div className="space-y-3">
              <div className="flex justify-between items-center p-3 bg-slate-700/30 rounded-lg">
                <span className="text-slate-300">Total Banned</span>
                <span className="text-xl font-bold text-white">{stats.total_banned}</span>
              </div>
              <div className="flex justify-between items-center p-3 bg-slate-700/30 rounded-lg">
                <span className="text-slate-300">Permanent</span>
                <span className="text-xl font-bold text-red-400">{stats.permanent_bans}</span>
              </div>
              <div className="flex justify-between items-center p-3 bg-slate-700/30 rounded-lg">
                <span className="text-slate-300">Temporary</span>
                <span className="text-xl font-bold text-yellow-400">{stats.temporary_bans}</span>
              </div>
              <div className="flex justify-between items-center p-3 bg-slate-700/30 rounded-lg">
                <span className="text-slate-300">Expiring Today</span>
                <span className="text-xl font-bold text-orange-400">{stats.expiring_today}</span>
              </div>
              <div className="flex justify-between items-center p-3 bg-slate-700/30 rounded-lg">
                <span className="text-slate-300">Expiring This Week</span>
                <span className="text-xl font-bold text-blue-400">{stats.expiring_week}</span>
              </div>
            </div>
          </div>

          {/* Role Distribution */}
          <div className="bg-slate-800/50 rounded-lg p-6 border border-slate-700">
            <h3 className="text-lg font-semibold text-white mb-4">Role Distribution</h3>
            <div className="space-y-3">
              {overview.role_distribution.map((role) => (
                <div key={role.role} className="flex items-center gap-3">
                  <div className="flex-1">
                    <div className="flex justify-between items-center mb-1">
                      <span className="text-sm capitalize text-slate-300">{role.role}</span>
                      <span className="text-sm font-bold text-white">{role.count}</span>
                    </div>
                    <div className="w-full bg-slate-700 rounded-full h-2 overflow-hidden">
                      <div
                        className={`h-full rounded-full ${
                          role.role === 'owner'
                            ? 'bg-purple-500'
                            : role.role === 'mod'
                              ? 'bg-blue-500'
                              : role.role === 'vip'
                                ? 'bg-yellow-500'
                                : 'bg-slate-500'
                        }`}
                        style={{
                          width: `${(role.count / kpi.total_players) * 100}%`,
                        }}
                      />
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Recent Activity */}
        <RecentActivity activities={overview.recent_activities} loading={loading} />
      </main>
    </div>
  );
}
