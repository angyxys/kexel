import { InvitationManager } from '../components/InvitationManager';
import { Navigation } from '../components/Navigation';

export function Invitations() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800">
      <Navigation />

      <main className="p-8">
        <div className="max-w-4xl mx-auto">
          <div className="mb-8">
            <h1 className="text-4xl font-bold text-white mb-2">Invitation Management</h1>
            <p className="text-slate-400">Create and manage invitation codes for new users</p>
          </div>

        <div className="grid gap-6">
          <InvitationManager />

          {/* Information Cards */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="bg-slate-800/50 rounded-lg p-4 border border-slate-700">
              <div className="flex items-center mb-2">
                <span className="text-2xl mr-2">📝</span>
                <h3 className="font-semibold text-white">How It Works</h3>
              </div>
              <p className="text-sm text-slate-400">
                Generate invitation codes and share them with users. They can use these codes during registration to join your server with assigned roles.
              </p>
            </div>

            <div className="bg-slate-800/50 rounded-lg p-4 border border-slate-700">
              <div className="flex items-center mb-2">
                <span className="text-2xl mr-2">🔧</span>
                <h3 className="font-semibold text-white">Customization</h3>
              </div>
              <p className="text-sm text-slate-400">
                Configure role assignments, usage limits, and expiration dates for each invitation code.
              </p>
            </div>

            <div className="bg-slate-800/50 rounded-lg p-4 border border-slate-700">
              <div className="flex items-center mb-2">
                <span className="text-2xl mr-2">🔐</span>
                <h3 className="font-semibold text-white">Security</h3>
              </div>
              <p className="text-sm text-slate-400">
                Revoke codes at any time and track usage to prevent unauthorized access.
              </p>
            </div>
          </div>

          {/* Usage Examples */}
          <div className="bg-slate-800/50 rounded-lg p-6 border border-slate-700">
            <h3 className="text-lg font-semibold text-white mb-4">Registration Link</h3>
            <p className="text-sm text-slate-400 mb-3">
              Share the registration link with your invitation code pre-filled:
            </p>
            <div className="bg-slate-900 rounded p-3 text-slate-300 font-mono text-sm overflow-x-auto">
              https://yourdomain.com/register?code=XXXXXXXXXXXXXXXX
            </div>
          </div>
        </div>
        </div>
      </main>
    </div>
  );
}
