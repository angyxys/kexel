import { useEffect, useState } from 'react';
import { integrationsApi, DiscordIntegration } from '../api/integrations';

export function DiscordManager() {
  const [integration, setIntegration] = useState<DiscordIntegration | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [showSetup, setShowSetup] = useState(false);
  const [testing, setTesting] = useState(false);

  // Setup form
  const [formBotToken, setFormBotToken] = useState('');
  const [formGuildID, setFormGuildID] = useState('');
  const [settingUpForm, setSettingUpForm] = useState(false);

  // Channels form
  const [formModLogChannel, setFormModLogChannel] = useState('');
  const [formAnnouncementChannel, setFormAnnouncementChannel] = useState('');
  const [updatingChannels, setUpdatingChannels] = useState(false);

  useEffect(() => {
    loadIntegration();
  }, []);

  const loadIntegration = async () => {
    try {
      setLoading(true);
      setError('');
      const data = await integrationsApi.getDiscordIntegration();
      setIntegration(data);
      setFormModLogChannel(data.mod_log_channel_id || '');
      setFormAnnouncementChannel(data.announcement_channel_id || '');
    } catch (err: any) {
      if (err.response?.status === 404) {
        setIntegration(null);
      } else {
        setError(err.response?.data?.message || 'Failed to load Discord integration');
      }
    } finally {
      setLoading(false);
    }
  };

  const handleSetupDiscord = async () => {
    if (!formBotToken || !formGuildID) {
      setError('Bot token and Guild ID are required');
      return;
    }

    try {
      setSettingUpForm(true);
      setError('');
      const data = await integrationsApi.setupDiscord(formBotToken, formGuildID);
      setIntegration(data);
      setFormBotToken('');
      setFormGuildID('');
      setShowSetup(false);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to setup Discord');
    } finally {
      setSettingUpForm(false);
    }
  };

  const handleUpdateChannels = async () => {
    if (!integration) return;

    try {
      setUpdatingChannels(true);
      setError('');
      await integrationsApi.configureDiscordChannels(formModLogChannel, formAnnouncementChannel);
      setIntegration({
        ...integration,
        mod_log_channel_id: formModLogChannel,
        announcement_channel_id: formAnnouncementChannel,
      });
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to update channels');
    } finally {
      setUpdatingChannels(false);
    }
  };

  const handleTestConnection = async () => {
    try {
      setTesting(true);
      setError('');
      await integrationsApi.testDiscordConnection();
      alert('Discord connection successful!');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to test connection');
    } finally {
      setTesting(false);
    }
  };

  const handleDisconnect = async () => {
    if (!confirm('Disconnect Discord bot? You can reconnect later.')) return;

    try {
      setError('');
      await integrationsApi.disconnectDiscord();
      setIntegration(null);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to disconnect');
    }
  };

  if (loading) {
    return <div className="text-slate-400">Loading Discord integration...</div>;
  }

  return (
    <div className="space-y-6">
      {error && (
        <div className="p-3 bg-red-900/30 border border-red-700 rounded text-red-200 text-sm">
          {error}
        </div>
      )}

      {!integration ? (
        // Setup Form
        showSetup ? (
          <div className="p-4 rounded-lg bg-slate-700/30 border border-slate-600 space-y-4">
            <h4 className="font-semibold text-white">Connect Discord Bot</h4>

            <div>
              <label className="block text-sm font-medium text-slate-300 mb-1">Bot Token</label>
              <input
                type="password"
                value={formBotToken}
                onChange={(e) => setFormBotToken(e.target.value)}
                placeholder="Enter Discord bot token"
                className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-slate-300 mb-1">Guild ID</label>
              <input
                type="text"
                value={formGuildID}
                onChange={(e) => setFormGuildID(e.target.value)}
                placeholder="Enter Discord server ID"
                className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div className="flex gap-2">
              <button
                onClick={handleSetupDiscord}
                disabled={settingUpForm}
                className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-slate-600 text-white rounded-lg transition-colors"
              >
                {settingUpForm ? 'Connecting...' : 'Connect'}
              </button>
              <button
                onClick={() => setShowSetup(false)}
                className="flex-1 px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg transition-colors"
              >
                Cancel
              </button>
            </div>
          </div>
        ) : (
          <button
            onClick={() => setShowSetup(true)}
            className="w-full px-4 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors font-medium"
          >
            🎮 Connect Discord Bot
          </button>
        )
      ) : (
        // Connected State
        <div className="space-y-4">
          <div className="p-4 rounded-lg bg-green-900/30 border border-green-700">
            <div className="flex items-center justify-between mb-2">
              <h4 className="font-semibold text-green-200">✓ Discord Connected</h4>
              <button
                onClick={handleDisconnect}
                className="px-3 py-1 bg-red-700 hover:bg-red-600 text-red-200 rounded text-xs transition-colors"
              >
                Disconnect
              </button>
            </div>
            <p className="text-sm text-green-100">Guild ID: {integration.guild_id}</p>
          </div>

          <div className="p-4 rounded-lg bg-slate-700/30 border border-slate-600 space-y-4">
            <h4 className="font-semibold text-white">Configure Channels</h4>

            <div>
              <label className="block text-sm font-medium text-slate-300 mb-1">Mod Log Channel</label>
              <input
                type="text"
                value={formModLogChannel}
                onChange={(e) => setFormModLogChannel(e.target.value)}
                placeholder="Channel ID for moderation logs"
                className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-slate-300 mb-1">Announcement Channel</label>
              <input
                type="text"
                value={formAnnouncementChannel}
                onChange={(e) => setFormAnnouncementChannel(e.target.value)}
                placeholder="Channel ID for announcements"
                className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div className="flex gap-2">
              <button
                onClick={handleUpdateChannels}
                disabled={updatingChannels}
                className="flex-1 px-4 py-2 bg-green-600 hover:bg-green-700 disabled:bg-slate-600 text-white rounded-lg transition-colors"
              >
                {updatingChannels ? 'Saving...' : 'Save Channels'}
              </button>
              <button
                onClick={handleTestConnection}
                disabled={testing}
                className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-slate-600 text-white rounded-lg transition-colors"
              >
                {testing ? 'Testing...' : 'Test Connection'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
