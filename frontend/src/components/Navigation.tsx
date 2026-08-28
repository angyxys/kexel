import { useNavigate, useLocation } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';

export function Navigation() {
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout } = useAuthStore();

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  const isActive = (path: string) => {
    return location.pathname === path ? 'bg-blue-600 text-white' : 'bg-slate-700 hover:bg-slate-600 text-white';
  };

  return (
    <header className="bg-slate-800/50 border-b border-slate-700 sticky top-0 z-40">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div className="flex justify-between items-center mb-4">
          <div>
            <h1 className="text-2xl font-bold text-white">Kexel</h1>
            <p className="text-sm text-slate-400">Welcome, {user?.username}</p>
          </div>
          <button
            onClick={handleLogout}
            className="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors"
          >
            Logout
          </button>
        </div>

        {/* Navigation */}
        <nav className="flex gap-2 border-t border-slate-700 pt-4 overflow-x-auto">
          <a
            href="/dashboard"
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors whitespace-nowrap ${isActive(
              '/dashboard'
            )}`}
          >
            Players
          </a>
          <a
            href="/audit-logs"
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors whitespace-nowrap ${isActive(
              '/audit-logs'
            )}`}
          >
            Audit Logs
          </a>
          <a
            href="/stats"
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors whitespace-nowrap ${isActive(
              '/stats'
            )}`}
          >
            Statistics
          </a>
          <a
            href="/invitations"
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors whitespace-nowrap ${isActive(
              '/invitations'
            )}`}
          >
            Invitations
          </a>
          <a
            href="/sessions"
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors whitespace-nowrap ${isActive(
              '/sessions'
            )}`}
          >
            Sessions
          </a>
          <a
            href="/security"
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors whitespace-nowrap ${isActive(
              '/security'
            )}`}
          >
            Security
          </a>
          <a
            href="/webhooks"
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors whitespace-nowrap ${isActive(
              '/webhooks'
            )}`}
          >
            Webhooks
          </a>
          <a
            href="/integrations"
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors whitespace-nowrap ${isActive(
              '/integrations'
            )}`}
          >
            Integrations
          </a>
          <a
            href="/banners"
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors whitespace-nowrap ${isActive(
              '/banners'
            )}`}
          >
            Banners
          </a>
          <a
            href="/tickets"
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors whitespace-nowrap ${isActive(
              '/tickets'
            )}`}
          >
            Support
          </a>
        </nav>
      </div>
    </header>
  );
}
