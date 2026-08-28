import { WebhookManager } from '../components/WebhookManager';
import { Navigation } from '../components/Navigation';

export function Webhooks() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800">
      <Navigation />

      <main className="p-8">
        <div className="max-w-4xl mx-auto">
          <div className="mb-8">
            <h1 className="text-4xl font-bold text-white mb-2">Webhooks</h1>
            <p className="text-slate-400">Subscribe to real-time events and integrate with external systems</p>
          </div>

          <div className="grid gap-6">
            <WebhookManager />

            {/* Documentation */}
            <div className="bg-slate-800/50 rounded-lg p-6 border border-slate-700 space-y-4">
              <h3 className="text-lg font-semibold text-white">📚 Webhook Documentation</h3>

              <div className="space-y-4">
                <div>
                  <h4 className="font-semibold text-white mb-2">Event Payload Structure</h4>
                  <div className="bg-slate-900 rounded p-3 text-slate-300 font-mono text-xs overflow-x-auto">
                    <pre>{`{
  "event_type": "player.created",
  "timestamp": "2026-08-27T12:00:00Z",
  "data": {
    "id": "usr_123456",
    "email": "user@example.com"
  }
}`}</pre>
                  </div>
                </div>

                <div>
                  <h4 className="font-semibold text-white mb-2">Signature Verification</h4>
                  <p className="text-sm text-slate-400 mb-2">
                    Each webhook request includes a signature header for verification:
                  </p>
                  <div className="bg-slate-900 rounded p-3 text-slate-300 font-mono text-xs">
                    X-Webhook-Signature: sha256=&lt;hmac_hash&gt;
                  </div>
                </div>

                <div>
                  <h4 className="font-semibold text-white mb-2">Retry Policy</h4>
                  <ul className="text-sm text-slate-400 space-y-1 list-disc list-inside">
                    <li>Failed requests are retried up to 5 times</li>
                    <li>Retry delay increases exponentially: 5, 10, 15, 20, 25 minutes</li>
                    <li>Success = HTTP 2xx status code</li>
                  </ul>
                </div>

                <div>
                  <h4 className="font-semibold text-white mb-2">Available Events</h4>
                  <ul className="text-sm text-slate-400 space-y-1 list-disc list-inside">
                    <li><code className="bg-slate-900 px-2 py-1 rounded">player.created</code> - New player added</li>
                    <li><code className="bg-slate-900 px-2 py-1 rounded">player.updated</code> - Player info changed</li>
                    <li><code className="bg-slate-900 px-2 py-1 rounded">player.deleted</code> - Player removed</li>
                    <li><code className="bg-slate-900 px-2 py-1 rounded">ban.created</code> - New ban issued</li>
                    <li><code className="bg-slate-900 px-2 py-1 rounded">ban.updated</code> - Ban modified</li>
                    <li><code className="bg-slate-900 px-2 py-1 rounded">ban.deleted</code> - Ban removed</li>
                    <li><code className="bg-slate-900 px-2 py-1 rounded">*</code> - Subscribe to all events</li>
                  </ul>
                </div>
              </div>
            </div>

            {/* Best Practices */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="bg-slate-800/50 rounded-lg p-4 border border-slate-700">
                <h4 className="font-semibold text-white mb-2">✓ Best Practices</h4>
                <ul className="text-sm text-slate-400 space-y-1">
                  <li>• Use HTTPS for webhook URLs</li>
                  <li>• Verify HMAC signature</li>
                  <li>• Return 2xx status quickly</li>
                  <li>• Implement idempotency</li>
                </ul>
              </div>

              <div className="bg-slate-800/50 rounded-lg p-4 border border-slate-700">
                <h4 className="font-semibold text-white mb-2">⚠️ Common Issues</h4>
                <ul className="text-sm text-slate-400 space-y-1">
                  <li>• Endpoint returns non-2xx status</li>
                  <li>• Endpoint is unreachable/down</li>
                  <li>• Timeout (10 second limit)</li>
                  <li>• Invalid HMAC signature</li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
