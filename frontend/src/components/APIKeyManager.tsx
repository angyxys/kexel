import { useEffect, useState } from 'react';
import { apiKeysApi, APIKey, CreateAPIKeyResponse } from '../api/api-keys';

export function APIKeyManager() {
  const [keys, setKeys] = useState<APIKey[]>([]);
  const [availableScopes, setAvailableScopes] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [showCreate, setShowCreate] = useState(false);
  const [newKey, setNewKey] = useState<CreateAPIKeyResponse | null>(null);
  const [copied, setCopied] = useState(false);

  // Form state
  const [formName, setFormName] = useState('');
  const [formScopes, setFormScopes] = useState<string[]>([]);
  const [formExpiresIn, setFormExpiresIn] = useState<number | ''>('');
  const [formRateLimit, setFormRateLimit] = useState(1000);
  const [creating, setCreating] = useState(false);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      setError('');
      const [keysData, scopesData] = await Promise.all([apiKeysApi.getAPIKeys(), apiKeysApi.getAvailableScopes()]);
      setKeys(keysData);
      setAvailableScopes(scopesData);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load API keys');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateKey = async () => {
    if (!formName || formScopes.length === 0) {
      setError('Please fill in all required fields');
      return;
    }

    try {
      setCreating(true);
      setError('');
      const response = await apiKeysApi.createAPIKey({
        name: formName,
        scopes: formScopes,
        expires_in: formExpiresIn ? (formExpiresIn as number) : undefined,
        rate_limit: formRateLimit,
      });
      setNewKey(response);
      setFormName('');
      setFormScopes([]);
      setFormExpiresIn('');
      setFormRateLimit(1000);
      await loadData();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to create API key');
    } finally {
      setCreating(false);
    }
  };

  const handleRevokeKey = async (keyId: number) => {
    if (!confirm('Are you sure you want to revoke this API key?')) return;

    try {
      await apiKeysApi.revokeAPIKey(keyId);
      setKeys(keys.map((k) => (k.id === keyId ? { ...k, is_active: false } : k)));
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to revoke API key');
    }
  };

  const handleDeleteKey = async (keyId: number) => {
    if (!confirm('This will permanently delete this API key. Continue?')) return;

    try {
      await apiKeysApi.deleteAPIKey(keyId);
      setKeys(keys.filter((k) => k.id !== keyId));
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to delete API key');
    }
  };

  const toggleScope = (scope: string) => {
    if (formScopes.includes(scope)) {
      setFormScopes(formScopes.filter((s) => s !== scope));
    } else {
      setFormScopes([...formScopes, scope]);
    }
  };

  const copyToClipboard = () => {
    if (newKey?.key) {
      navigator.clipboard.writeText(newKey.key);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  return (
    <div className="space-y-6">
      {error && (
        <div className="p-3 bg-red-900/30 border border-red-700 rounded text-red-200 text-sm">
          {error}
        </div>
      )}

      {/* New Key Display */}
      {newKey && (
        <div className="p-4 rounded-lg bg-green-900/30 border border-green-700">
          <h4 className="font-semibold text-green-200 mb-3">✓ API Key Created Successfully</h4>
          <p className="text-sm text-green-100 mb-3">
            ⚠️ Copy this key now. You won't be able to see it again!
          </p>
          <div className="bg-slate-900 rounded p-3 font-mono text-sm text-slate-300 break-all mb-3 border border-slate-600">
            {newKey.key}
          </div>
          <div className="flex gap-2">
            <button
              onClick={copyToClipboard}
              className="flex-1 px-3 py-2 bg-green-600 hover:bg-green-700 text-green-100 rounded text-sm transition-colors"
            >
              {copied ? '✓ Copied!' : '📋 Copy Key'}
            </button>
            <button
              onClick={() => setNewKey(null)}
              className="flex-1 px-3 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded text-sm transition-colors"
            >
              Done
            </button>
          </div>
        </div>
      )}

      {/* Create Form */}
      {showCreate && !newKey && (
        <div className="p-4 rounded-lg bg-slate-700/30 border border-slate-600 space-y-4">
          <h4 className="font-semibold text-white">Create New API Key</h4>

          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1">Name</label>
            <input
              type="text"
              value={formName}
              onChange={(e) => setFormName(e.target.value)}
              placeholder="My Integration"
              className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-300 mb-2">Scopes</label>
            <div className="grid grid-cols-2 gap-2">
              {availableScopes.map((scope) => (
                <label key={scope} className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={formScopes.includes(scope)}
                    onChange={() => toggleScope(scope)}
                    className="w-4 h-4"
                  />
                  <span className="text-sm text-slate-300">{scope}</span>
                </label>
              ))}
            </div>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="block text-sm font-medium text-slate-300 mb-1">Expires In (days)</label>
              <input
                type="number"
                value={formExpiresIn}
                onChange={(e) => setFormExpiresIn(e.target.value ? parseInt(e.target.value) : '')}
                placeholder="Leave empty for no expiry"
                className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-slate-300 mb-1">Rate Limit (req/hr)</label>
              <input
                type="number"
                value={formRateLimit}
                onChange={(e) => setFormRateLimit(parseInt(e.target.value) || 1000)}
                className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>

          <div className="flex gap-2">
            <button
              onClick={handleCreateKey}
              disabled={creating || !formName || formScopes.length === 0}
              className="flex-1 px-4 py-2 bg-green-600 hover:bg-green-700 disabled:bg-slate-600 text-white rounded-lg transition-colors font-medium"
            >
              {creating ? 'Creating...' : 'Create API Key'}
            </button>
            <button
              onClick={() => {
                setShowCreate(false);
                setFormName('');
                setFormScopes([]);
                setFormExpiresIn('');
                setFormRateLimit(1000);
              }}
              className="flex-1 px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg transition-colors"
            >
              Cancel
            </button>
          </div>
        </div>
      )}

      {/* Header */}
      <div className="flex justify-between items-center">
        <h3 className="text-lg font-semibold text-white">API Keys</h3>
        {!showCreate && !newKey && (
          <button
            onClick={() => setShowCreate(true)}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors text-sm"
          >
            ➕ New Key
          </button>
        )}
      </div>

      {/* Keys List */}
      {loading ? (
        <p className="text-slate-400">Loading API keys...</p>
      ) : keys.length === 0 ? (
        <p className="text-slate-400 text-sm">No API keys yet</p>
      ) : (
        <div className="space-y-3">
          {keys.map((key) => (
            <div key={key.id} className="p-4 bg-slate-800/50 rounded-lg border border-slate-700 hover:border-slate-500 transition-colors">
              <div className="flex items-start justify-between mb-2">
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-2">
                    <span className="font-semibold text-white">{key.name}</span>
                    {key.is_active ? (
                      <span className="px-2 py-1 rounded text-xs font-medium bg-green-900 text-green-200">Active</span>
                    ) : (
                      <span className="px-2 py-1 rounded text-xs font-medium bg-red-900 text-red-200">Revoked</span>
                    )}
                  </div>

                  <div className="space-y-1 text-xs text-slate-400 mb-2">
                    <p>
                      🔑 Key: <span className="font-mono text-slate-300">{key.key_prefix}...</span>
                    </p>
                    <p>📅 Created: {new Date(key.created_at).toLocaleDateString()}</p>
                    {key.last_used_at && <p>⏱️ Last used: {new Date(key.last_used_at).toLocaleString()}</p>}
                    {key.expires_at && <p>⏰ Expires: {new Date(key.expires_at).toLocaleDateString()}</p>}
                  </div>

                  <div className="flex flex-wrap gap-1">
                    {key.scopes.map((scope) => (
                      <span key={scope} className="px-2 py-1 rounded text-xs bg-slate-700 text-slate-300">
                        {scope}
                      </span>
                    ))}
                  </div>
                </div>

                {key.is_active && (
                  <div className="flex gap-2 ml-4">
                    <button
                      onClick={() => handleRevokeKey(key.id)}
                      className="px-3 py-1 bg-yellow-700 hover:bg-yellow-600 text-yellow-200 rounded text-xs transition-colors whitespace-nowrap"
                    >
                      Revoke
                    </button>
                    <button
                      onClick={() => handleDeleteKey(key.id)}
                      className="px-3 py-1 bg-red-700 hover:bg-red-600 text-red-200 rounded text-xs transition-colors"
                    >
                      Delete
                    </button>
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
