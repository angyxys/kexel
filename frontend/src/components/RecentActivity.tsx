import { ActivityLog } from '../api/stats';

interface RecentActivityProps {
  activities: ActivityLog[];
  loading?: boolean;
}

export function RecentActivity({ activities, loading = false }: RecentActivityProps) {
  const getActionIcon = (action: string) => {
    switch (action) {
      case 'POST':
        return '➕';
      case 'PUT':
        return '✏️';
      case 'DELETE':
        return '🗑️';
      default:
        return '📝';
    }
  };

  const getActionColor = (action: string) => {
    switch (action) {
      case 'POST':
        return 'text-green-400';
      case 'PUT':
        return 'text-blue-400';
      case 'DELETE':
        return 'text-red-400';
      default:
        return 'text-slate-400';
    }
  };

  if (loading) {
    return (
      <div className="bg-slate-800/50 rounded-lg p-6 border border-slate-700">
        <h3 className="text-lg font-semibold text-white mb-4">Recent Activity</h3>
        <p className="text-slate-400">Loading...</p>
      </div>
    );
  }

  return (
    <div className="bg-slate-800/50 rounded-lg p-6 border border-slate-700">
      <h3 className="text-lg font-semibold text-white mb-4">Recent Activity</h3>

      {activities.length === 0 ? (
        <p className="text-slate-400 text-sm">No recent activity</p>
      ) : (
        <div className="space-y-3">
          {activities.map((activity) => (
            <div
              key={activity.id}
              className="flex items-start gap-3 p-3 bg-slate-700/30 rounded-lg hover:bg-slate-700/50 transition-colors"
            >
              <span className={`text-xl ${getActionColor(activity.action)}`}>
                {getActionIcon(activity.action)}
              </span>
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-1">
                  <p className="font-medium text-white text-sm">{activity.user}</p>
                  <span className="text-xs text-slate-400">{activity.action}</span>
                </div>
                <p className="text-sm text-slate-300 truncate">{activity.resource}</p>
                <p className="text-xs text-slate-400 mt-1">{activity.details}</p>
                <p className="text-xs text-slate-500 mt-1">
                  {new Date(activity.timestamp).toLocaleString()}
                </p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
