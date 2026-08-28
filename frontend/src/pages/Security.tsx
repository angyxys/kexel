import { TwoFactorManager } from '../components/TwoFactorManager';
import { APIKeyManager } from '../components/APIKeyManager';
import { Navigation } from '../components/Navigation';

export function Security() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800">
      <Navigation />

      <main className="p-8">
        <div className="max-w-4xl mx-auto">
          <div className="mb-8">
            <h1 className="text-4xl font-bold text-white mb-2">Security Settings</h1>
            <p className="text-slate-400">Manage your account security and privacy</p>
          </div>

          <div className="grid gap-6">
            {/* 2FA Manager */}
            <TwoFactorManager />

            {/* API Keys Manager */}
            <div className="bg-slate-800/50 rounded-lg p-6 border border-slate-700">
              <APIKeyManager />
            </div>

            {/* Security Tips */}
            <div className="bg-slate-800/50 rounded-lg p-6 border border-slate-700 space-y-4">
              <h3 className="text-lg font-semibold text-white">🛡️ Security Best Practices</h3>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="p-4 rounded-lg bg-slate-700/30 border border-slate-600">
                  <h4 className="font-semibold text-white mb-2">Strong Passwords</h4>
                  <p className="text-sm text-slate-400">
                    Use unique passwords with a mix of uppercase, lowercase, numbers, and symbols.
                  </p>
                </div>

                <div className="p-4 rounded-lg bg-slate-700/30 border border-slate-600">
                  <h4 className="font-semibold text-white mb-2">Enable 2FA</h4>
                  <p className="text-sm text-slate-400">
                    Use two-factor authentication to add an extra layer of security to your account.
                  </p>
                </div>

                <div className="p-4 rounded-lg bg-slate-700/30 border border-slate-600">
                  <h4 className="font-semibold text-white mb-2">Review Sessions</h4>
                  <p className="text-sm text-slate-400">
                    Check your active sessions regularly and logout any devices you don't recognize.
                  </p>
                </div>

                <div className="p-4 rounded-lg bg-slate-700/30 border border-slate-600">
                  <h4 className="font-semibold text-white mb-2">Monitor Activity</h4>
                  <p className="text-sm text-slate-400">
                    Review your audit logs to detect any suspicious activities on your account.
                  </p>
                </div>
              </div>
            </div>

            {/* Account Data */}
            <div className="bg-slate-800/50 rounded-lg p-6 border border-slate-700">
              <h3 className="text-lg font-semibold text-white mb-4">📊 Account Data</h3>

              <div className="space-y-3">
                <div className="flex items-center justify-between p-3 bg-slate-700/30 rounded-lg">
                  <span className="text-slate-300">Download Account Data</span>
                  <button className="px-3 py-1 bg-blue-600 hover:bg-blue-700 text-white rounded text-sm transition-colors">
                    Download
                  </button>
                </div>

                <div className="flex items-center justify-between p-3 bg-slate-700/30 rounded-lg">
                  <span className="text-slate-300">Delete Account</span>
                  <button className="px-3 py-1 bg-red-700 hover:bg-red-600 text-red-200 rounded text-sm transition-colors">
                    Delete
                  </button>
                </div>
              </div>

              <p className="text-xs text-slate-400 mt-4">
                ⚠️ These actions are permanent and cannot be undone.
              </p>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
