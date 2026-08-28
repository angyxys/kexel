import { useEffect, useState } from 'react';
import { webhooksApi, Webhook, WebhookEvent } from '../api/webhooks';

export function WebhookManager() {
  const [webhooks, setWebhooks] = useState<Webhook[]>([]);
  const [availableEvents, setAvailableEvents] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [showCreate, setShowCreate] = useState(false);
  const [expandedWebhook, setExpandedWebhook] = useState<number | null>(null);
  const [webhookEvents, setWebhookEvents] = useState<WebhookEvent[]>([]);
  const [loadingEvents, setLoadingEvents] = useState(false);

  // Form state
  const [formName, setFormName] = useState('');
  const [formURL, setFormURL] = useState('');
  const [formEvents, setFormEvents] = useState<string[]>([]);
  const [creating, setCreating] = useState(false);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      setError('');
      const [webhooksData, eventsData] = await Promise.all([webhooksApi.getWebhooks(), webhooksApi.getAvailableEvents()]);
      setWebhooks(webhooksData);
      setAvailableEvents(eventsData);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load webhooks');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateWebhook = async () => {
    if (!formName || !formURL || formEvents.length === 0) {
      setError('Please fill in all required fields');
      return;
    }

    try {
      setCreating(true);
      setError('');
      await webhooksApi.createWebhook(formName, formURL, formEvents);
      setFormName('');
      setFormURL('');
      setFormEvents([]);
      setShowCreate(false);
      await loadData();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to create webhook');
    } finally {
      setCreating(false);
    }
  };

  const handleDeleteWebhook = async (webhookId: number) => {
    if (!confirm('Are you sure you want to delete this webhook?')) return;

    try {
      await webhooksApi.deleteWebhook(webhookId);
      setWebhooks(webhooks.filter((w) => w.id !== webhookId));
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to delete webhook');
    }
  };

  const handleDisableWebhook = async (webhookId: number) => {
    try {
      await webhooksApi.disableWebhook(webhookId);
      setWebhooks(webhooks.map((w) => (w.id === webhookId ? { ...w, is_active: false } : w)));
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to disable webhook');
    }
  };

  const handleExpandWebhook = async (webhookId: number) => {
    if (expandedWebhook === webhookId) {
      setExpandedWebhook(null);
      return;
    }

    try {
      setExpandedWebhook(webhookId);
      setLoadingEvents(true);
      const events = await webhooksApi.getWebhookEvents(webhookId);
      setWebhookEvents(events);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load webhook events');
    } finally {
      setLoadingEvents(false);
    }
  };

  const toggleEvent = (event: string) => {
    if (formEvents.includes(event)) {
      setFormEvents(formEvents.filter((e) => e !== event));
    } else {
      setFormEvents([...formEvents, event]);
    }
  };

  const getStatusColor = (webhook: Webhook) => {
    if (!webhook.is_active) return 'bg-red-900/30 border-red-700';
    if (webhook.failure_count > 5) return 'bg-yellow-900/30 border-yellow-700';
    if (webhook.last_success_at) return 'bg-green-900/30 border-green-700';
    return 'bg-slate-700/30 border-slate-600';
  };

  return (
    <div className="space-y-6">
      {error && (
        <div className="p-3 bg-red-900/30 border border-red-700 rounded text-red-200 text-sm">
          {error}
        </div>
      )}

      {/* Create Form */}
      {showCreate && (
        <div className="p-4 rounded-lg bg-slate-700/30 border border-slate-600 space-y-4">
          <h4 className="font-semibold text-white">Create New Webhook</h4>

          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1">Webhook Name</label>
            <input
              type="text"
              value={formName}
              onChange={(e) => setFormName(e.target.value)}
              placeholder="My Webhook"
              className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-300 mb-1">Webhook URL</label>
            <input
              type="url"
              value={formURL}
              onChange={(e) => setFormURL(e.target.value)}
              placeholder="https://example.com/webhook"
              className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-300 mb-2">Subscribe to Events</label>
            <div className="grid grid-cols-2 gap-2">
              {availableEvents.map((event) => (
                <label key={event} className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={formEvents.includes(event)}
                    onChange={() => toggleEvent(event)}
                    className="w-4 h-4"
                  />
                  <span className="text-sm text-slate-300">{event}</span>
                </label>
              ))}
            </div>
          </div>

          <div className="flex gap-2">
            <button
              onClick={handleCreateWebhook}
              disabled={creating || !formName || !formURL || formEvents.length === 0}
              className="flex-1 px-4 py-2 bg-green-600 hover:bg-green-700 disabled:bg-slate-600 text-white rounded-lg transition-colors font-medium"
            >
              {creating ? 'Creating...' : 'Create Webhook'}
            </button>
            <button
              onClick={() => {
                setShowCreate(false);
                setFormName('');
                setFormURL('');
                setFormEvents([]);
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
        <h3 className="text-lg font-semibold text-white">Webhooks</h3>
        {!showCreate && (
          <button
            onClick={() => setShowCreate(true)}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors text-sm"
          >
            ➕ New Webhook
          </button>
        )}
      </div>

      {/* Webhooks List */}
      {loading ? (
        <p className="text-slate-400">Loading webhooks...</p>
      ) : webhooks.length === 0 ? (
        <p className="text-slate-400 text-sm">No webhooks yet</p>
      ) : (
        <div className="space-y-3">
          {webhooks.map((webhook) => (
            <div
              key={webhook.id}
              className={`p-4 rounded-lg border transition-colors ${getStatusColor(webhook)}`}
            >
              <div className="flex items-start justify-between mb-2">
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-2">
                    <span className="font-semibold text-white">{webhook.name}</span>
                    {webhook.is_active ? (
                      <span className="px-2 py-1 rounded text-xs font-medium bg-green-900 text-green-200">
                        Active
                      </span>
                    ) : (
                      <span className="px-2 py-1 rounded text-xs font-medium bg-red-900 text-red-200">
                        Disabled
                      </span>
                    )}
                    {webhook.failure_count > 0 && (
                      <span className="px-2 py-1 rounded text-xs font-medium bg-yellow-900 text-yellow-200">
                        ⚠️ {webhook.failure_count} failures
                      </span>
                    )}
                  </div>

                  <p className="text-xs text-slate-400 font-mono break-all mb-2">{webhook.url}</p>

                  <div className="space-y-1 text-xs text-slate-400 mb-2">
                    <p>📅 Created: {new Date(webhook.created_at).toLocaleDateString()}</p>
                    {webhook.last_success_at && (
                      <p>✓ Last success: {new Date(webhook.last_success_at).toLocaleString()}</p>
                    )}
                    {webhook.last_tried_at && (
                      <p>⏱️ Last attempt: {new Date(webhook.last_tried_at).toLocaleString()}</p>
                    )}
                  </div>

                  <div className="flex flex-wrap gap-1">
                    {webhook.events.map((event) => (
                      <span key={event} className="px-2 py-1 rounded text-xs bg-slate-700 text-slate-300">
                        {event}
                      </span>
                    ))}
                  </div>
                </div>

                {webhook.is_active && (
                  <div className="flex gap-2 ml-4">
                    <button
                      onClick={() => handleDisableWebhook(webhook.id)}
                      className="px-3 py-1 bg-yellow-700 hover:bg-yellow-600 text-yellow-200 rounded text-xs transition-colors whitespace-nowrap"
                    >
                      Disable
                    </button>
                    <button
                      onClick={() => handleDeleteWebhook(webhook.id)}
                      className="px-3 py-1 bg-red-700 hover:bg-red-600 text-red-200 rounded text-xs transition-colors"
                    >
                      Delete
                    </button>
                  </div>
                )}
              </div>

              {/* Expandable Events Section */}
              <button
                onClick={() => handleExpandWebhook(webhook.id)}
                className="w-full mt-3 px-3 py-2 text-left text-xs bg-slate-700/50 hover:bg-slate-700 text-slate-300 rounded transition-colors"
              >
                {expandedWebhook === webhook.id ? '▼' : '▶'} Recent Events ({webhookEvents.length})
              </button>

              {expandedWebhook === webhook.id && (
                <div className="mt-3 space-y-2 max-h-48 overflow-y-auto">
                  {loadingEvents ? (
                    <p className="text-xs text-slate-400">Loading events...</p>
                  ) : webhookEvents.length === 0 ? (
                    <p className="text-xs text-slate-400">No events yet</p>
                  ) : (
                    webhookEvents.map((event) => (
                      <div
                        key={event.id}
                        className={`p-2 rounded text-xs ${
                          event.is_delivered ? 'bg-green-900/30 text-green-300' : 'bg-red-900/30 text-red-300'
                        }`}
                      >
                        <div className="flex justify-between items-center">
                          <span className="font-mono">{event.event_type}</span>
                          <span className="text-xs">{event.status_code}</span>
                        </div>
                        <div className="text-xs opacity-75">Attempt {event.attempts}</div>
                      </div>
                    ))
                  )}
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
