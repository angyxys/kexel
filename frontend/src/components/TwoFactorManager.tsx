import { useEffect, useState } from 'react';
import { twoFactorApi, TOTPSetup, TOTPStatus } from '../api/2fa';
import { useAuthStore } from '../store/authStore';

export function TwoFactorManager() {
  const { user } = useAuthStore();
  const [status, setStatus] = useState<TOTPStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [showSetup, setShowSetup] = useState(false);
  const [setup, setSetup] = useState<TOTPSetup | null>(null);
  const [verifyCode, setVerifyCode] = useState('');
  const [verifying, setVerifying] = useState(false);
  const [showBackupCodes, setShowBackupCodes] = useState(false);

  useEffect(() => {
    loadStatus();
  }, []);

  const loadStatus = async () => {
    try {
      setLoading(true);
      setError('');
      const data = await twoFactorApi.getTOTPStatus();
      setStatus(data);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load 2FA status');
    } finally {
      setLoading(false);
    }
  };

  const handleStartSetup = async () => {
    try {
      setError('');
      if (!user?.email) {
        setError('Email not found');
        return;
      }
      const setupData = await twoFactorApi.setupTOTP(user.email);
      setSetup(setupData);
      setShowSetup(true);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to setup 2FA');
    }
  };

  const handleVerifySetup = async () => {
    if (!verifyCode) {
      setError('Please enter a code');
      return;
    }

    try {
      setVerifying(true);
      setError('');
      await twoFactorApi.verifyTOTP(verifyCode);
      setShowSetup(false);
      setVerifyCode('');
      setSetup(null);
      await loadStatus();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Invalid code');
    } finally {
      setVerifying(false);
    }
  };

  const handleDisable2FA = async () => {
    if (!confirm('This will disable 2FA on your account. Continue?')) return;

    try {
      setError('');
      await twoFactorApi.disableTOTP();
      await loadStatus();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to disable 2FA');
    }
  };

  if (loading) {
    return <div className="text-slate-400">Loading 2FA status...</div>;
  }

  return (
    <div className="bg-slate-800/50 rounded-lg p-6 border border-slate-700 space-y-6">
      <div>
        <h3 className="text-lg font-semibold text-white mb-2">Two-Factor Authentication</h3>
        <p className="text-sm text-slate-400">
          Add an extra layer of security to your account with 2FA using an authenticator app.
        </p>
      </div>

      {error && (
        <div className="p-3 bg-red-900/30 border border-red-700 rounded text-red-200 text-sm">
          {error}
        </div>
      )}

      {/* Status */}
      {status && !showSetup && (
        <div className="p-4 rounded-lg bg-slate-700/30 border border-slate-600">
          <div className="flex items-center justify-between mb-3">
            <span className="text-slate-300 font-medium">Status</span>
            {status.is_enabled ? (
              <span className="px-3 py-1 rounded-lg text-xs font-medium bg-green-900 text-green-200">
                🔒 Enabled
              </span>
            ) : (
              <span className="px-3 py-1 rounded-lg text-xs font-medium bg-slate-900 text-slate-200">
                ⭕ Disabled
              </span>
            )}
          </div>

          {status.is_enabled && (
            <div className="space-y-2 text-sm text-slate-400 mt-4">
              <p>📝 Backup codes remaining: <strong className="text-white">{status.backup_codes_left}</strong></p>
              {status.enabled_at && (
                <p>✓ Enabled on: <strong className="text-white">{new Date(status.enabled_at).toLocaleDateString()}</strong></p>
              )}
            </div>
          )}
        </div>
      )}

      {/* Setup Form */}
      {showSetup && setup ? (
        <div className="space-y-4 p-4 rounded-lg bg-slate-700/30 border border-slate-600">
          <div>
            <h4 className="font-semibold text-white mb-2">Step 1: Scan QR Code</h4>
            <p className="text-sm text-slate-400 mb-3">
              Scan this QR code with your authenticator app (Google Authenticator, Authy, Microsoft Authenticator, etc):
            </p>
            <div className="flex justify-center bg-white p-4 rounded-lg">
              <img src={setup.qr_code} alt="TOTP QR Code" style={{ maxWidth: '250px' }} />
            </div>
          </div>

          <div>
            <h4 className="font-semibold text-white mb-2">Step 2: Manual Entry (Optional)</h4>
            <p className="text-xs text-slate-400 mb-2">If you can't scan the QR code, enter this secret manually:</p>
            <div className="bg-slate-900 rounded p-3 font-mono text-sm text-slate-300 break-all border border-slate-600">
              {setup.secret}
            </div>
          </div>

          <div>
            <h4 className="font-semibold text-white mb-2">Step 3: Verify Code</h4>
            <p className="text-xs text-slate-400 mb-2">Enter the 6-digit code from your authenticator app:</p>
            <input
              type="text"
              maxLength={6}
              placeholder="000000"
              value={verifyCode}
              onChange={(e) => setVerifyCode(e.target.value.replace(/\D/g, ''))}
              className="w-full px-4 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white text-center text-2xl tracking-widest focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div className="flex gap-3">
            <button
              onClick={handleVerifySetup}
              disabled={verifying || verifyCode.length !== 6}
              className="flex-1 px-4 py-2 bg-green-600 hover:bg-green-700 disabled:bg-slate-600 text-white rounded-lg transition-colors font-medium"
            >
              {verifying ? 'Verifying...' : 'Verify & Enable'}
            </button>
            <button
              onClick={() => {
                setShowSetup(false);
                setSetup(null);
                setVerifyCode('');
              }}
              className="flex-1 px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg transition-colors"
            >
              Cancel
            </button>
          </div>
        </div>
      ) : (
        /* Actions */
        <div className="flex gap-3">
          {!status?.is_enabled ? (
            <button
              onClick={handleStartSetup}
              className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors font-medium"
            >
              🔐 Enable 2FA
            </button>
          ) : (
            <>
              <button
                onClick={() => setShowBackupCodes(!showBackupCodes)}
                className="flex-1 px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg transition-colors"
              >
                {showBackupCodes ? 'Hide Backup Codes' : 'View Backup Codes'}
              </button>
              <button
                onClick={handleDisable2FA}
                className="flex-1 px-4 py-2 bg-red-700 hover:bg-red-600 text-red-200 rounded-lg transition-colors"
              >
                ⚠️ Disable 2FA
              </button>
            </>
          )}
        </div>
      )}

      {/* Backup Codes */}
      {showBackupCodes && setup?.backup_codes && (
        <div className="p-4 rounded-lg bg-yellow-900/20 border border-yellow-700">
          <h4 className="font-semibold text-yellow-200 mb-3">⚠️ Save Your Backup Codes</h4>
          <p className="text-sm text-yellow-100 mb-3">
            Store these codes in a safe place. You can use them to access your account if you lose access to your authenticator app.
          </p>
          <div className="grid grid-cols-2 gap-2 bg-slate-900 p-3 rounded font-mono text-sm text-slate-300">
            {setup.backup_codes.map((code, i) => (
              <div key={i} className="p-2 bg-slate-800 rounded border border-slate-600">
                {code}
              </div>
            ))}
          </div>
          <button
            onClick={() => {
              const text = setup.backup_codes.join('\n');
              navigator.clipboard.writeText(text);
              alert('Backup codes copied!');
            }}
            className="mt-3 w-full px-3 py-2 bg-yellow-700 hover:bg-yellow-600 text-yellow-100 rounded text-sm transition-colors"
          >
            📋 Copy All Codes
          </button>
        </div>
      )}

      {/* Info */}
      {!status?.is_enabled && (
        <div className="p-4 rounded-lg bg-blue-900/20 border border-blue-700">
          <h4 className="font-semibold text-blue-200 mb-2">Supported Authenticators</h4>
          <ul className="text-sm text-blue-100 space-y-1">
            <li>• Google Authenticator</li>
            <li>• Authy</li>
            <li>• Microsoft Authenticator</li>
            <li>• FreeOTP</li>
            <li>• 1Password</li>
          </ul>
        </div>
      )}
    </div>
  );
}
