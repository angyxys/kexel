import { useState } from 'react';
import { PatreonInfo, configurePatreonTierMapping, syncPatreonMembers, disconnectPatreon } from '../api/patreon';

interface PatreonManagerProps {
  patreon: PatreonInfo | null;
  onSync: () => void;
  onDisconnect: () => void;
}

const KEXEL_ROLES = ['user', 'vip', 'mod', 'owner'];

export function PatreonManager({ patreon, onSync, onDisconnect }: PatreonManagerProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [tierMapping, setTierMapping] = useState<Record<string, string>>(patreon?.tier_mapping || {});
  const [editMode, setEditMode] = useState(false);

  if (!patreon) {
    return (
      <div className="bg-yellow-500/20 border border-yellow-500 rounded-lg p-6">
        <p className="text-yellow-200">Patreon integration not configured</p>
      </div>
    );
  }

  const handleTierMappingChange = (tierId: string, role: string) => {
    setTierMapping(prev => ({
      ...prev,
      [tierId]: role,
    }));
  };

  const handleSaveTierMapping = async () => {
    setLoading(true);
    setError(null);

    try {
      await configurePatreonTierMapping(tierMapping);
      setEditMode(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save tier mapping');
    } finally {
      setLoading(false);
    }
  };

  const handleSync = async () => {
    setLoading(true);
    setError(null);

    try {
      await syncPatreonMembers();
      onSync();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to sync patrons');
    } finally {
      setLoading(false);
    }
  };

  const handleDisconnect = async () => {
    if (!confirm('Are you sure you want to disconnect Patreon?')) return;

    setLoading(true);
    setError(null);

    try {
      await disconnectPatreon();
      onDisconnect();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to disconnect');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      {error && (
        <div className="bg-red-500/20 border border-red-500 rounded-lg p-4">
          <p className="text-red-200">{error}</p>
        </div>
      )}

      <div className="bg-slate-700/50 border border-slate-600 rounded-lg p-6 space-y-4">
        <div>
          <h3 className="text-lg font-semibold text-white mb-4">Patreon Integration Status</h3>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <p className="text-sm text-slate-400">Campaign ID</p>
              <p className="text-white font-mono">{patreon.campaign_id}</p>
            </div>
            <div>
              <p className="text-sm text-slate-400">Status</p>
              <p className={`text-white ${patreon.is_enabled ? 'text-green-400' : 'text-red-400'}`}>
                {patreon.is_enabled ? '✓ Enabled' : '✗ Disabled'}
              </p>
            </div>
            {patreon.last_sync_at && (
              <div className="col-span-2">
                <p className="text-sm text-slate-400">Last Sync</p>
                <p className="text-white">
                  {new Date(patreon.last_sync_at).toLocaleString()}
                </p>
              </div>
            )}
          </div>
        </div>

        <div className="flex gap-2 flex-wrap">
          <button
            onClick={handleSync}
            disabled={loading}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-slate-600 text-white rounded-lg transition-colors"
          >
            {loading ? 'Syncing...' : 'Sync Patrons Now'}
          </button>
          <button
            onClick={() => setEditMode(!editMode)}
            disabled={loading}
            className="px-4 py-2 bg-slate-600 hover:bg-slate-500 text-white rounded-lg transition-colors"
          >
            {editMode ? 'Cancel' : 'Configure Tiers'}
          </button>
          <button
            onClick={handleDisconnect}
            disabled={loading}
            className="px-4 py-2 bg-red-600 hover:bg-red-700 disabled:bg-slate-600 text-white rounded-lg transition-colors"
          >
            Disconnect
          </button>
        </div>
      </div>

      {editMode && (
        <div className="bg-slate-700/50 border border-slate-600 rounded-lg p-6 space-y-4">
          <h3 className="text-lg font-semibold text-white">Configure Tier Mapping</h3>
          <p className="text-sm text-slate-400">Map Patreon tiers to Kexel roles</p>

          <div className="space-y-3">
            {Object.entries(tierMapping).map(([tierId, role]) => (
              <div key={tierId} className="flex items-center gap-2">
                <div className="flex-1">
                  <label className="block text-sm text-slate-300 mb-1">Tier ID: {tierId}</label>
                  <select
                    value={role}
                    onChange={(e) => handleTierMappingChange(tierId, e.target.value)}
                    disabled={loading}
                    className="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded-lg text-white"
                  >
                    <option value="">Select role...</option>
                    {KEXEL_ROLES.map(r => (
                      <option key={r} value={r}>{r}</option>
                    ))}
                  </select>
                </div>
              </div>
            ))}

            {Object.keys(tierMapping).length === 0 && (
              <p className="text-slate-400 text-sm">No tiers configured yet. They will appear here after syncing patrons.</p>
            )}
          </div>

          <div className="flex gap-2">
            <button
              onClick={handleSaveTierMapping}
              disabled={loading}
              className="px-4 py-2 bg-green-600 hover:bg-green-700 disabled:bg-slate-600 text-white rounded-lg transition-colors"
            >
              {loading ? 'Saving...' : 'Save Tier Mapping'}
            </button>
            <button
              onClick={() => setEditMode(false)}
              disabled={loading}
              className="px-4 py-2 bg-slate-600 hover:bg-slate-500 text-white rounded-lg transition-colors"
            >
              Cancel
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
