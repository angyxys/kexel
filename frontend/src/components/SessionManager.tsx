import { useEffect, useState } from 'react';
import { sessionsApi, SessionInfo, SessionStats } from '../api/sessions';

export function SessionManager() {
  const [sessions, setSessions] = useState<SessionInfo[]>([]);
  const [stats, setStats] = useState<SessionStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    loadSessions();
  }, []);

  const loadSessions = async () => {
    try {
      setLoading(true);
      setError('');
      const [sessionsData, statsData] = await Promise.all([
        sessionsApi.getSessions(),
        sessionsApi.getSessionStats(),
      ]);
      setSessions(sessionsData);
      setStats(statsData);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load sessions');
    } finally {
      setLoading(false);
    }
  };

  const handleLogoutSession = async (sessionId: number) => {
    if (!confirm('Are you sure you want to logout this session?')) return;

    try {
      await sessionsApi.logoutSession(sessionId);
      setSessions(sessions.filter((s) => s.id !== sessionId));
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to logout session');
    }
  };

  const handleLogoutAllOther = async () => {
    if (!confirm('This will logout all other sessions. Continue?')) return;

    try {
      await sessionsApi.logoutAllOtherSessions();
      loadSessions();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to logout sessions');
    }
  };

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleString();
  };

  const getTimeAgo = (dateStr: string) => {
    const date = new Date(dateStr);
    const now = new Date();
    const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);

    if (seconds < 60) return 'Just now';
    if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
    if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
    return `${Math.floor(seconds / 86400)}d ago`;
  };

  return (
    <div className="space-y-6">
      {/* Stats */}
      {stats && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="bg-slate-800/50 rounded-lg p-4 border border-slate-700">
            <p className="text-slate-400 text-sm">Active Sessions</p>
            <p className="text-3xl font-bold text-green-400 mt-2">{stats.total_active_sessions}</p>
          </div>

          <div className="bg-slate-800/50 rounded-lg p-4 border border-slate-700">
            <p className="text-slate-400 text-sm">Current Device</p>
            <p className="text-lg font-semibold text-white mt-2">{stats.current_device}</p>
            <p className="text-xs text-slate-400 mt-1">IP: {stats.current_ip}</p>
          </div>
        </div>
      )}

      {/* Error Message */}
      {error && (
        <div className="p-3 bg-red-900/30 border border-red-700 rounded text-red-200 text-sm">
          {error}
        </div>
      )}

      {/* Header */}
      <div className="flex justify-between items-center">
        <h3 className="text-lg font-semibold text-white">Active Sessions</h3>
        {sessions.length > 1 && (
          <button
            onClick={handleLogoutAllOther}
            className="px-3 py-1 bg-red-700 hover:bg-red-600 text-red-200 rounded text-sm transition-colors"
          >
            Logout Other Sessions
          </button>
        )}
      </div>

      {/* Sessions List */}
      {loading ? (
        <p className="text-slate-400">Loading sessions...</p>
      ) : sessions.length === 0 ? (
        <p className="text-slate-400 text-sm">No active sessions</p>
      ) : (
        <div className="space-y-3">
          {sessions.map((session) => (
            <div
              key={session.id}
              className="p-4 bg-slate-800/50 rounded-lg border border-slate-700 hover:border-slate-500 transition-colors"
            >
              <div className="flex items-start justify-between mb-2">
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-2">
                    <span className="font-semibold text-white">📱 {session.device_name}</span>
                    {session.is_active ? (
                      <span className="px-2 py-1 rounded text-xs font-medium bg-green-900 text-green-200">
                        Active
                      </span>
                    ) : (
                      <span className="px-2 py-1 rounded text-xs font-medium bg-slate-900 text-slate-200">
                        Inactive
                      </span>
                    )}
                  </div>

                  <div className="space-y-1 text-sm text-slate-400">
                    <p>🌐 IP: <span className="font-mono text-slate-300">{session.ip_address}</span></p>
                    <p>🔓 Logged in: {formatDate(session.login_at)}</p>
                    <p>⏱️ Last active: {getTimeAgo(session.last_activity)}</p>
                    <p>⏰ Expires: {formatDate(session.expires_at)}</p>
                  </div>
                </div>

                {session.is_active && (
                  <button
                    onClick={() => handleLogoutSession(session.id)}
                    className="ml-4 px-3 py-1 bg-red-700 hover:bg-red-600 text-red-200 rounded text-sm transition-colors whitespace-nowrap"
                  >
                    Logout
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
