import { useEffect, useState } from 'react';
import { auditApi, AuditLog, AuditStats } from '../api/audit';
import { Navigation } from '../components/Navigation';

export function AuditLogs() {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [stats, setStats] = useState<AuditStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');

  // Filters
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [totalPages, setTotalPages] = useState(1);
  const [action, setAction] = useState('');
  const [resourceType, setResourceType] = useState('');
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');

  useEffect(() => {
    loadData();
  }, [page, action, resourceType, startDate, endDate]);

  const loadData = async () => {
    try {
      setLoading(true);
      setError('');

      const [logsData, statsData] = await Promise.all([
        auditApi.listLogs(page, pageSize, {
          action: action || undefined,
          resource_type: resourceType || undefined,
          start_date: startDate || undefined,
          end_date: endDate || undefined,
        }),
        auditApi.getStats(),
      ]);

      setLogs(logsData.data);
      setTotalPages(logsData.total_pages);
      setStats(statsData);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load audit logs');
    } finally {
      setLoading(false);
    }
  };

  const handleExport = async () => {
    try {
      const blob = await auditApi.exportLogs({
        action: action || undefined,
        resource_type: resourceType || undefined,
        start_date: startDate || undefined,
        end_date: endDate || undefined,
      });

      // Create download link
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `audit-logs-${new Date().toISOString().split('T')[0]}.csv`;
      a.click();
      window.URL.revokeObjectURL(url);
    } catch (err: any) {
      setError('Failed to export logs');
    }
  };

  const resetFilters = () => {
    setPage(1);
    setAction('');
    setResourceType('');
    setStartDate('');
    setEndDate('');
  };

  const getActionBadgeColor = (action: string) => {
    switch (action) {
      case 'POST':
        return 'bg-green-900 text-green-200';
      case 'PUT':
        return 'bg-blue-900 text-blue-200';
      case 'DELETE':
        return 'bg-red-900 text-red-200';
      case 'GET':
        return 'bg-slate-700 text-slate-200';
      default:
        return 'bg-slate-700 text-slate-200';
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800">
      <Navigation />

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Stats Cards */}
        {stats && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
            <div className="bg-slate-800/50 border border-slate-700 rounded-lg p-6">
              <div className="text-slate-400 text-sm font-medium">Total Logs</div>
              <div className="text-3xl font-bold text-white mt-2">{stats.total_logs}</div>
            </div>
            <div className="bg-slate-800/50 border border-slate-700 rounded-lg p-6">
              <div className="text-slate-400 text-sm font-medium">Today</div>
              <div className="text-3xl font-bold text-white mt-2">{stats.today_logs}</div>
            </div>
            <div className="bg-slate-800/50 border border-slate-700 rounded-lg p-6">
              <div className="text-slate-400 text-sm font-medium">Unique Users</div>
              <div className="text-3xl font-bold text-white mt-2">{stats.unique_users}</div>
            </div>
          </div>
        )}

        {/* Error Message */}
        {error && (
          <div className="mb-6 p-4 bg-red-900/30 border border-red-700 rounded-lg text-red-200 text-sm">
            {error}
          </div>
        )}

        {/* Filters */}
        <div className="bg-slate-800/50 border border-slate-700 rounded-lg p-6 mb-8">
          <h3 className="text-lg font-semibold text-white mb-4">Filters</h3>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-4">
            <div>
              <label className="block text-sm font-medium text-slate-300 mb-2">Action</label>
              <select
                value={action}
                onChange={(e) => {
                  setAction(e.target.value);
                  setPage(1);
                }}
                className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="">All Actions</option>
                <option value="POST">Create (POST)</option>
                <option value="PUT">Update (PUT)</option>
                <option value="DELETE">Delete (DELETE)</option>
                <option value="GET">Read (GET)</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-slate-300 mb-2">Resource Type</label>
              <select
                value={resourceType}
                onChange={(e) => {
                  setResourceType(e.target.value);
                  setPage(1);
                }}
                className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="">All Types</option>
                <option value="player">Player</option>
                <option value="user">User</option>
                <option value="invitation">Invitation</option>
                <option value="session">Session</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-slate-300 mb-2">Start Date</label>
              <input
                type="date"
                value={startDate}
                onChange={(e) => {
                  setStartDate(e.target.value);
                  setPage(1);
                }}
                className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-slate-300 mb-2">End Date</label>
              <input
                type="date"
                value={endDate}
                onChange={(e) => {
                  setEndDate(e.target.value);
                  setPage(1);
                }}
                className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>

          <div className="flex gap-2">
            <button
              onClick={resetFilters}
              className="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg transition-colors"
            >
              Reset Filters
            </button>
            <button
              onClick={handleExport}
              className="px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg transition-colors"
            >
              Export CSV
            </button>
          </div>
        </div>

        {/* Logs Table */}
        {loading ? (
          <div className="text-center py-8">
            <p className="text-slate-300">Loading audit logs...</p>
          </div>
        ) : logs.length === 0 ? (
          <div className="text-center py-8 bg-slate-800/30 rounded-lg border border-slate-700">
            <p className="text-slate-300">No logs found matching your criteria</p>
          </div>
        ) : (
          <>
            <div className="overflow-x-auto bg-slate-800/50 rounded-lg border border-slate-700">
              <table className="w-full">
                <thead className="bg-slate-800 border-b border-slate-700">
                  <tr>
                    <th className="px-6 py-3 text-left text-sm font-semibold text-slate-300">Time</th>
                    <th className="px-6 py-3 text-left text-sm font-semibold text-slate-300">User</th>
                    <th className="px-6 py-3 text-left text-sm font-semibold text-slate-300">Action</th>
                    <th className="px-6 py-3 text-left text-sm font-semibold text-slate-300">Resource</th>
                    <th className="px-6 py-3 text-left text-sm font-semibold text-slate-300">Description</th>
                    <th className="px-6 py-3 text-left text-sm font-semibold text-slate-300">IP Address</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-700">
                  {logs.map((log) => (
                    <tr key={log.id} className="hover:bg-slate-700/50 transition-colors">
                      <td className="px-6 py-4 text-sm text-slate-200 whitespace-nowrap">
                        {new Date(log.created_at).toLocaleString()}
                      </td>
                      <td className="px-6 py-4 text-sm text-slate-200">{log.username}</td>
                      <td className="px-6 py-4 text-sm">
                        <span className={`px-2 py-1 rounded text-xs font-medium ${getActionBadgeColor(log.action)}`}>
                          {log.action}
                        </span>
                      </td>
                      <td className="px-6 py-4 text-sm text-slate-200">
                        <div className="text-xs font-mono text-slate-400">{log.resource_type}</div>
                        <div className="text-xs font-mono text-slate-300">{log.resource_id}</div>
                      </td>
                      <td className="px-6 py-4 text-sm text-slate-300 max-w-xs truncate">{log.description}</td>
                      <td className="px-6 py-4 text-sm text-slate-400 font-mono text-xs">{log.ip_address}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Pagination */}
            <div className="mt-6 flex justify-between items-center">
              <div className="text-sm text-slate-400">
                Page {page} of {totalPages} ({logs.length} entries)
              </div>
              <div className="flex gap-2">
                <button
                  onClick={() => setPage(Math.max(1, page - 1))}
                  disabled={page === 1}
                  className="px-4 py-2 bg-slate-700 hover:bg-slate-600 disabled:opacity-50 text-white rounded-lg transition-colors"
                >
                  Previous
                </button>
                <button
                  onClick={() => setPage(Math.min(totalPages, page + 1))}
                  disabled={page === totalPages}
                  className="px-4 py-2 bg-slate-700 hover:bg-slate-600 disabled:opacity-50 text-white rounded-lg transition-colors"
                >
                  Next
                </button>
              </div>
            </div>
          </>
        )}
      </main>
    </div>
  );
}
