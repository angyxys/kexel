import { useEffect, useState } from 'react';
import { invitationsApi, InvitationCode } from '../api/invitations';

export function InvitationManager() {
  const [invitations, setInvitations] = useState<InvitationCode[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [showForm, setShowForm] = useState(false);
  const [role, setRole] = useState('user');
  const [maxUses, setMaxUses] = useState(1);
  const [expiresAt, setExpiresAt] = useState('');
  const [creating, setCreating] = useState(false);

  useEffect(() => {
    loadInvitations();
  }, []);

  const loadInvitations = async () => {
    try {
      setLoading(true);
      setError('');
      const data = await invitationsApi.getMyInvitations();
      setInvitations(data);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load invitations');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateInvitation = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      setCreating(true);
      setError('');

      const request: any = {
        role,
        max_uses: maxUses,
      };

      if (expiresAt) {
        request.expires_at = expiresAt;
      }

      const newInvitation = await invitationsApi.createInvitation(request);
      setInvitations([newInvitation, ...invitations]);
      setShowForm(false);
      setRole('user');
      setMaxUses(1);
      setExpiresAt('');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to create invitation');
    } finally {
      setCreating(false);
    }
  };

  const handleRevokeInvitation = async (id: number) => {
    if (!confirm('Are you sure you want to revoke this invitation?')) return;

    try {
      await invitationsApi.revokeInvitation(id);
      setInvitations(invitations.filter((inv) => inv.id !== id));
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to revoke invitation');
    }
  };

  const getRoleColor = (role: string) => {
    switch (role) {
      case 'owner':
        return 'bg-purple-900 text-purple-200';
      case 'mod':
        return 'bg-blue-900 text-blue-200';
      case 'vip':
        return 'bg-yellow-900 text-yellow-200';
      default:
        return 'bg-slate-700 text-slate-200';
    }
  };

  return (
    <div className="bg-slate-800/50 rounded-lg p-6 border border-slate-700">
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-lg font-semibold text-white">Invitation Codes</h3>
        <button
          onClick={() => setShowForm(!showForm)}
          className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors text-sm"
        >
          {showForm ? 'Cancel' : '➕ Generate Code'}
        </button>
      </div>

      {error && (
        <div className="mb-4 p-3 bg-red-900/30 border border-red-700 rounded text-red-200 text-sm">
          {error}
        </div>
      )}

      {/* Create Form */}
      {showForm && (
        <div className="mb-6 p-4 bg-slate-700/30 rounded-lg border border-slate-600">
          <form onSubmit={handleCreateInvitation} className="space-y-3">
            <div>
              <label className="block text-sm font-medium text-slate-300 mb-1">Role</label>
              <select
                value={role}
                onChange={(e) => setRole(e.target.value)}
                className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="user">User</option>
                <option value="vip">VIP</option>
                <option value="mod">Moderator</option>
                <option value="owner">Owner</option>
              </select>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="block text-sm font-medium text-slate-300 mb-1">Max Uses</label>
                <input
                  type="number"
                  min="-1"
                  value={maxUses}
                  onChange={(e) => setMaxUses(parseInt(e.target.value))}
                  className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
                <p className="text-xs text-slate-400 mt-1">-1 = unlimited</p>
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-300 mb-1">Expires At</label>
                <input
                  type="datetime-local"
                  value={expiresAt}
                  onChange={(e) => setExpiresAt(e.target.value)}
                  className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            <button
              type="submit"
              disabled={creating}
              className="w-full px-4 py-2 bg-green-600 hover:bg-green-700 disabled:bg-slate-600 text-white rounded-lg transition-colors font-medium"
            >
              {creating ? 'Creating...' : 'Generate Invitation Code'}
            </button>
          </form>
        </div>
      )}

      {/* Invitations List */}
      {loading ? (
        <p className="text-slate-400">Loading...</p>
      ) : invitations.length === 0 ? (
        <p className="text-slate-400 text-sm">No invitation codes generated yet</p>
      ) : (
        <div className="space-y-3">
          {invitations.map((inv) => (
            <div
              key={inv.id}
              className="p-4 bg-slate-700/30 rounded-lg border border-slate-600 hover:border-slate-500 transition-colors"
            >
              <div className="flex items-start justify-between mb-2">
                <div>
                  <div className="flex items-center gap-2 mb-2">
                    <span className={`px-2 py-1 rounded text-xs font-medium ${getRoleColor(inv.role)}`}>
                      {inv.role}
                    </span>
                    {inv.is_active ? (
                      <span className="px-2 py-1 rounded text-xs font-medium bg-green-900 text-green-200">
                        Active
                      </span>
                    ) : (
                      <span className="px-2 py-1 rounded text-xs font-medium bg-red-900 text-red-200">
                        Revoked
                      </span>
                    )}
                  </div>
                  <p className="font-mono text-sm text-slate-200 break-all mb-1">{inv.code}</p>
                  <p className="text-xs text-slate-400">
                    Uses: {inv.uses} / {inv.max_uses === -1 ? '∞' : inv.max_uses}
                  </p>
                  {inv.expires_at && (
                    <p className="text-xs text-slate-400">
                      Expires: {new Date(inv.expires_at).toLocaleString()}
                    </p>
                  )}
                </div>

                <div className="flex gap-2">
                  <button
                    onClick={() => {
                      navigator.clipboard.writeText(inv.code);
                      alert('Code copied!');
                    }}
                    className="px-3 py-1 bg-slate-700 hover:bg-slate-600 text-slate-300 rounded text-xs transition-colors"
                  >
                    📋 Copy
                  </button>
                  {inv.is_active && (
                    <button
                      onClick={() => handleRevokeInvitation(inv.id)}
                      className="px-3 py-1 bg-red-700 hover:bg-red-600 text-red-200 rounded text-xs transition-colors"
                    >
                      Revoke
                    </button>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
