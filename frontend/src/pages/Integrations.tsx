import { useState, useEffect } from 'react';
import { DiscordManager } from '../components/DiscordManager';
import { PatreonManager } from '../components/PatreonManager';
import { Navigation } from '../components/Navigation';
import { getPatreonIntegration, getPatreonOAuthURL } from '../api/patreon';
import { PatreonInfo } from '../api/patreon';

export function Integrations() {
  const [patreon, setPatreon] = useState<PatreonInfo | null>(null);
  const [loadingPatreon, setLoadingPatreon] = useState(true);

  useEffect(() => {
    loadPatreonData();
  }, []);

  const loadPatreonData = async () => {
    try {
      const data = await getPatreonIntegration();
      setPatreon(data);
    } catch {
      setPatreon(null);
    } finally {
      setLoadingPatreon(false);
    }
  };

  const handlePatreonOAuth = () => {
    const state = Math.random().toString(36).substring(7);
    const url = getPatreonOAuthURL(state);
    window.location.href = url;
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800">
      <Navigation />

      <main className="p-8">
        <div className="max-w-4xl mx-auto">
          <div className="mb-8">
            <h1 className="text-4xl font-bold text-white mb-2">Integrations</h1>
            <p className="text-slate-400">Connect external services and configure system-wide settings</p>
          </div>

          <div className="grid gap-6">
            {/* Discord */}
            <div className="bg-slate-800/50 rounded-lg p-6 border border-slate-700">
              <h3 className="text-lg font-semibold text-white mb-4">🎮 Discord Bot</h3>
              <p className="text-sm text-slate-400 mb-4">
                Connect your Discord server to receive moderation logs and announcements
              </p>
              <DiscordManager />
            </div>

            {/* Patreon */}
            <div className="bg-slate-800/50 rounded-lg p-6 border border-slate-700">
              <h3 className="text-lg font-semibold text-white mb-4">❤️ Patreon Integration</h3>
              <p className="text-sm text-slate-400 mb-4">
                Automatically sync Patreon supporters and grant VIP roles
              </p>
              {loadingPatreon ? (
                <p className="text-slate-400">Loading...</p>
              ) : patreon ? (
                <PatreonManager
                  patreon={patreon}
                  onSync={loadPatreonData}
                  onDisconnect={() => setPatreon(null)}
                />
              ) : (
                <button
                  onClick={handlePatreonOAuth}
                  className="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors"
                >
                  Connect with Patreon
                </button>
              )}
            </div>

            {/* Rate Limiting */}
            <div className="bg-slate-800/50 rounded-lg p-6 border border-slate-700">
              <h3 className="text-lg font-semibold text-white mb-4">⏱️ Rate Limiting</h3>
              <p className="text-sm text-slate-400 mb-4">
                Monitor and manage API rate limiting to prevent abuse
              </p>
              <div className="space-y-2">
                <p className="text-xs text-slate-400">
                  Rate limiting is enabled by default with sensible defaults. Configure custom rules as needed.
                </p>
                <div className="p-3 rounded bg-slate-700/30 border border-slate-600 text-sm text-slate-300">
                  <strong>Default Limits:</strong>
                  <ul className="list-disc list-inside mt-2 space-y-1">
                    <li>1000 requests per hour per IP</li>
                    <li>Custom rules can be configured per endpoint</li>
                    <li>Automatic blocking with exponential backoff</li>
                  </ul>
                </div>
              </div>
            </div>

            {/* Integration Status */}
            <div className="bg-blue-900/30 rounded-lg p-6 border border-blue-700">
              <h3 className="text-lg font-semibold text-blue-200 mb-4">ℹ️ About Integrations</h3>
              <div className="space-y-3 text-sm text-blue-100">
                <p>
                  <strong>Discord Bot:</strong> Requires a bot application from Discord Developer Portal with permissions
                  to send messages and manage roles.
                </p>
                <p>
                  <strong>Patreon Sync:</strong> Connect your Patreon campaign to automatically assign roles based on
                  supporter tier.
                </p>
                <p>
                  <strong>Rate Limiting:</strong> Protects your API from abuse. Rules are applied per IP address and
                  endpoint combination.
                </p>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
