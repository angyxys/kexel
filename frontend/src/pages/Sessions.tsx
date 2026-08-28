import { SessionManager } from '../components/SessionManager';
import { Navigation } from '../components/Navigation';

export function Sessions() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800">
      <Navigation />

      <main className="p-8">
        <div className="max-w-4xl mx-auto">
          <div className="mb-8">
            <h1 className="text-4xl font-bold text-white mb-2">Session Management</h1>
            <p className="text-slate-400">Monitor and manage your active sessions across devices</p>
          </div>

          <div className="grid gap-6">
            <SessionManager />

            {/* Information Cards */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="bg-slate-800/50 rounded-lg p-4 border border-slate-700">
                <div className="flex items-center mb-2">
                  <span className="text-2xl mr-2">🔒</span>
                  <h3 className="font-semibold text-white">Security</h3>
                </div>
                <p className="text-sm text-slate-400">
                  Monitor all active sessions on your account. Each session is tracked with device information and IP address.
                </p>
              </div>

              <div className="bg-slate-800/50 rounded-lg p-4 border border-slate-700">
                <div className="flex items-center mb-2">
                  <span className="text-2xl mr-2">⚠️</span>
                  <h3 className="font-semibold text-white">Suspicious Activity</h3>
                </div>
                <p className="text-sm text-slate-400">
                  If you see sessions you don't recognize, logout them immediately. Sessions automatically expire after 7 days.
                </p>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
