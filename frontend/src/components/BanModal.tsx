import { useState } from 'react';
import { bansApi, BanRequest } from '../api/bans';

interface BanModalProps {
  playerID: string;
  onClose: () => void;
  onSuccess: () => void;
}

export function BanModal({ playerID, onClose, onSuccess }: BanModalProps) {
  const [reason, setReason] = useState('');
  const [banType, setBanType] = useState<'permanent' | 'temporary'>('temporary');
  const [duration, setDuration] = useState(24); // hours
  const [expiresAt, setExpiresAt] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleBan = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!reason.trim()) {
      setError('Ban reason is required');
      return;
    }

    try {
      setLoading(true);
      setError('');

      const request: BanRequest = {
        reason,
      };

      if (banType === 'permanent') {
        request.duration = 0; // 0 = permanent
      } else if (expiresAt) {
        // Use custom expiration date
        request.expires_at = expiresAt;
      } else {
        // Use duration in hours
        request.duration = duration;
      }

      await bansApi.banPlayer(playerID, request);
      onSuccess();
      onClose();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to ban player');
    } finally {
      setLoading(false);
    }
  };

  const getPresetDurations = [
    { label: '1 hour', value: 1 },
    { label: '24 hours (1 day)', value: 24 },
    { label: '7 days', value: 24 * 7 },
    { label: '30 days', value: 24 * 30 },
  ];

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
      <div className="bg-slate-800 rounded-lg shadow-xl max-w-md w-full p-6">
        <h2 className="text-xl font-bold text-white mb-4">Ban Player</h2>

        {error && (
          <div className="mb-4 p-3 bg-red-900/30 border border-red-700 rounded text-red-200 text-sm">
            {error}
          </div>
        )}

        <form onSubmit={handleBan} className="space-y-4">
          {/* Player ID (read-only) */}
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-2">Player ID</label>
            <input
              type="text"
              value={playerID}
              disabled
              className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-slate-400 opacity-50"
            />
          </div>

          {/* Ban Reason */}
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-2">Reason *</label>
            <textarea
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              placeholder="Explain why this player is being banned..."
              rows={3}
              className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          {/* Ban Type Selection */}
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-2">Ban Duration</label>
            <div className="space-y-2">
              <label className="flex items-center text-slate-300 cursor-pointer">
                <input
                  type="radio"
                  name="banType"
                  value="temporary"
                  checked={banType === 'temporary'}
                  onChange={(e) => setBanType('temporary')}
                  className="w-4 h-4"
                />
                <span className="ml-2">Temporary Ban</span>
              </label>
              <label className="flex items-center text-slate-300 cursor-pointer">
                <input
                  type="radio"
                  name="banType"
                  value="permanent"
                  checked={banType === 'permanent'}
                  onChange={(e) => setBanType('permanent')}
                  className="w-4 h-4"
                />
                <span className="ml-2">Permanent Ban</span>
              </label>
            </div>
          </div>

          {/* Temporary Ban Options */}
          {banType === 'temporary' && (
            <div className="bg-slate-700/30 p-4 rounded-lg border border-slate-700">
              <div className="space-y-3">
                <div>
                  <label className="block text-sm font-medium text-slate-300 mb-2">
                    Quick Duration Presets
                  </label>
                  <div className="grid grid-cols-2 gap-2">
                    {getPresetDurations.map((preset) => (
                      <button
                        key={preset.value}
                        type="button"
                        onClick={() => setDuration(preset.value)}
                        className={`px-3 py-2 rounded text-sm font-medium transition-colors ${
                          duration === preset.value
                            ? 'bg-blue-600 text-white'
                            : 'bg-slate-700 hover:bg-slate-600 text-slate-300'
                        }`}
                      >
                        {preset.label}
                      </button>
                    ))}
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium text-slate-300 mb-2">
                    Custom Duration (hours)
                  </label>
                  <input
                    type="number"
                    min="1"
                    value={duration}
                    onChange={(e) => setDuration(parseInt(e.target.value) || 1)}
                    className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-slate-300 mb-2">
                    Or select exact expiration date
                  </label>
                  <input
                    type="datetime-local"
                    value={expiresAt}
                    onChange={(e) => setExpiresAt(e.target.value)}
                    className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>

                {!expiresAt && (
                  <div className="text-xs text-slate-400">
                    Ban will expire in: <strong>{duration} hour(s)</strong>
                  </div>
                )}
                {expiresAt && (
                  <div className="text-xs text-slate-400">
                    Ban will expire on: <strong>{new Date(expiresAt).toLocaleString()}</strong>
                  </div>
                )}
              </div>
            </div>
          )}

          {banType === 'permanent' && (
            <div className="bg-red-900/20 p-3 rounded-lg border border-red-700 text-red-200 text-sm">
              ⚠️ This will be a permanent ban. The player can only be unbanned manually.
            </div>
          )}

          {/* Buttons */}
          <div className="flex gap-2 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading}
              className="flex-1 px-4 py-2 bg-red-600 hover:bg-red-700 disabled:bg-slate-600 text-white rounded-lg transition-colors font-medium"
            >
              {loading ? 'Banning...' : 'Ban Player'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
