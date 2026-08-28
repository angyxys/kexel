import { useState, useEffect } from 'react';
import { Navigation } from '../components/Navigation';
import { TicketList } from '../components/TicketList';
import { TicketDetail } from '../components/TicketDetail';
import {
  Ticket,
  TicketComment,
  CreateTicketRequest,
  getUserTickets,
  getTicketStats,
  getTicketComments,
  createTicket,
} from '../api/tickets';

type TabType = 'my-tickets' | 'create' | 'all';

export function Tickets() {
  const [activeTab, setActiveTab] = useState<TabType>('my-tickets');
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [selectedTicket, setSelectedTicket] = useState<Ticket | null>(null);
  const [comments, setComments] = useState<TicketComment[]>([]);
  const [stats, setStats] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Create form state
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    category: 'bug',
    priority: 'medium',
  });
  const [creating, setCreating] = useState(false);

  useEffect(() => {
    loadData();
  }, [activeTab]);

  useEffect(() => {
    if (selectedTicket) {
      loadComments();
    }
  }, [selectedTicket]);

  const loadData = async () => {
    setLoading(true);
    setError(null);

    try {
      const [ticketsData, statsData] = await Promise.all([
        getUserTickets(),
        getTicketStats(),
      ]);
      setTickets(ticketsData.data || []);
      setStats(statsData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load data');
    } finally {
      setLoading(false);
    }
  };

  const loadComments = async () => {
    if (!selectedTicket) return;

    try {
      const response = await getTicketComments(selectedTicket.id);
      setComments(response.data || []);
    } catch (err) {
      console.error('Failed to load comments:', err);
    }
  };

  const handleCreateTicket = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.title.trim() || !formData.description.trim()) {
      setError('Title and description are required');
      return;
    }

    setCreating(true);
    setError(null);

    try {
      const request: CreateTicketRequest = {
        title: formData.title,
        description: formData.description,
        category: formData.category,
        priority: formData.priority,
      };

      const newTicket = await createTicket(request);
      setTickets([newTicket, ...tickets]);
      setFormData({
        title: '',
        description: '',
        category: 'bug',
        priority: 'medium',
      });
      setActiveTab('my-tickets');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create ticket');
    } finally {
      setCreating(false);
    }
  };

  return (
    <>
      <Navigation />
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-white">Support Tickets</h1>
          <p className="text-slate-400 mt-2">Report bugs, request features, or get support</p>
        </div>

        {error && (
          <div className="mb-6 bg-red-500/20 border border-red-500 rounded-lg p-4">
            <p className="text-red-200">{error}</p>
          </div>
        )}

        {/* Stats */}
        {stats && (
          <div className="grid grid-cols-2 md:grid-cols-5 gap-4 mb-8">
            <div className="bg-slate-700/50 border border-slate-600 rounded-lg p-4">
              <p className="text-sm text-slate-400">Total</p>
              <p className="text-2xl font-bold text-white">{stats.total}</p>
            </div>
            <div className="bg-blue-500/20 border border-blue-500 rounded-lg p-4">
              <p className="text-sm text-blue-200">Open</p>
              <p className="text-2xl font-bold text-blue-100">{stats.open}</p>
            </div>
            <div className="bg-yellow-500/20 border border-yellow-500 rounded-lg p-4">
              <p className="text-sm text-yellow-200">In Progress</p>
              <p className="text-2xl font-bold text-yellow-100">{stats.in_progress}</p>
            </div>
            <div className="bg-green-500/20 border border-green-500 rounded-lg p-4">
              <p className="text-sm text-green-200">Resolved</p>
              <p className="text-2xl font-bold text-green-100">{stats.resolved}</p>
            </div>
            <div className="bg-slate-600/20 border border-slate-500 rounded-lg p-4">
              <p className="text-sm text-slate-300">Closed</p>
              <p className="text-2xl font-bold text-slate-100">{stats.closed}</p>
            </div>
          </div>
        )}

        {/* Tabs */}
        <div className="flex gap-2 mb-6 border-b border-slate-600">
          <button
            onClick={() => setActiveTab('my-tickets')}
            className={`px-4 py-2 font-medium transition-colors ${
              activeTab === 'my-tickets'
                ? 'text-blue-400 border-b-2 border-blue-400'
                : 'text-slate-400 hover:text-white'
            }`}
          >
            My Tickets
          </button>
          <button
            onClick={() => setActiveTab('create')}
            className={`px-4 py-2 font-medium transition-colors ${
              activeTab === 'create'
                ? 'text-blue-400 border-b-2 border-blue-400'
                : 'text-slate-400 hover:text-white'
            }`}
          >
            Create Ticket
          </button>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Left: Ticket List or Form */}
          <div>
            {activeTab === 'my-tickets' && (
              <TicketList
                tickets={tickets}
                onSelectTicket={setSelectedTicket}
                isLoading={loading}
              />
            )}

            {activeTab === 'create' && (
              <div className="bg-slate-700/50 border border-slate-600 rounded-lg p-6">
                <h2 className="text-lg font-semibold text-white mb-4">Create New Ticket</h2>
                <form onSubmit={handleCreateTicket} className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-slate-200 mb-2">
                      Category
                    </label>
                    <select
                      value={formData.category}
                      onChange={(e) =>
                        setFormData({ ...formData, category: e.target.value })
                      }
                      disabled={creating}
                      className="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded-lg text-white"
                    >
                      <option value="bug">🐛 Bug Report</option>
                      <option value="feature-request">✨ Feature Request</option>
                      <option value="other">📝 Other</option>
                    </select>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-slate-200 mb-2">
                      Priority
                    </label>
                    <select
                      value={formData.priority}
                      onChange={(e) =>
                        setFormData({ ...formData, priority: e.target.value })
                      }
                      disabled={creating}
                      className="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded-lg text-white"
                    >
                      <option value="low">Low</option>
                      <option value="medium">Medium</option>
                      <option value="high">High</option>
                      <option value="critical">Critical</option>
                    </select>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-slate-200 mb-2">
                      Title *
                    </label>
                    <input
                      type="text"
                      value={formData.title}
                      onChange={(e) =>
                        setFormData({ ...formData, title: e.target.value })
                      }
                      disabled={creating}
                      placeholder="Brief description of the issue..."
                      className="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded-lg text-white"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-slate-200 mb-2">
                      Description *
                    </label>
                    <textarea
                      value={formData.description}
                      onChange={(e) =>
                        setFormData({ ...formData, description: e.target.value })
                      }
                      disabled={creating}
                      placeholder="Detailed description, steps to reproduce, etc..."
                      rows={6}
                      className="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded-lg text-white"
                    />
                  </div>

                  <button
                    type="submit"
                    disabled={creating || !formData.title.trim() || !formData.description.trim()}
                    className="w-full px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-slate-600 text-white rounded-lg font-medium transition-colors"
                  >
                    {creating ? 'Creating...' : 'Create Ticket'}
                  </button>
                </form>
              </div>
            )}
          </div>

          {/* Right: Ticket Detail */}
          <div>
            {selectedTicket ? (
              <TicketDetail
                ticket={selectedTicket}
                comments={comments}
                onUpdate={loadData}
                isLoadingComments={false}
              />
            ) : (
              <div className="bg-slate-700/50 border border-slate-600 rounded-lg p-6 text-center">
                <p className="text-slate-400">
                  {activeTab === 'my-tickets'
                    ? 'Select a ticket to view details'
                    : 'Create a new ticket to get started'}
                </p>
              </div>
            )}
          </div>
        </div>
      </div>
    </>
  );
}
