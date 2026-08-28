import { useState } from 'react';
import { Ticket, TicketComment, updateTicket, addComment } from '../api/tickets';

interface TicketDetailProps {
  ticket: Ticket;
  comments: TicketComment[];
  onUpdate: () => void;
  isLoadingComments?: boolean;
}

const STATUSES = ['open', 'in-progress', 'resolved', 'closed'];
const PRIORITIES = ['low', 'medium', 'high', 'critical'];

export function TicketDetail({ ticket, comments, onUpdate, isLoadingComments }: TicketDetailProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [newComment, setNewComment] = useState('');
  const [isInternalComment, setIsInternalComment] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [editData, setEditData] = useState({
    status: ticket.status,
    priority: ticket.priority,
    resolution: ticket.resolution || '',
  });

  const handleSaveChanges = async () => {
    setLoading(true);
    setError(null);

    try {
      await updateTicket(ticket.id, {
        status: editData.status,
        priority: editData.priority,
        resolution: editData.resolution,
      });
      setIsEditing(false);
      onUpdate();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update ticket');
    } finally {
      setLoading(false);
    }
  };

  const handleAddComment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newComment.trim()) return;

    setLoading(true);
    setError(null);

    try {
      await addComment(ticket.id, {
        content: newComment,
        is_internal: isInternalComment,
      });
      setNewComment('');
      setIsInternalComment(false);
      onUpdate();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to add comment');
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

      {/* Ticket Details */}
      <div className="bg-slate-700/50 border border-slate-600 rounded-lg p-6 space-y-4">
        <div className="flex justify-between items-start mb-4">
          <div>
            <h1 className="text-2xl font-bold text-white mb-1">#{ticket.id} - {ticket.title}</h1>
            <p className="text-sm text-slate-400">
              Created by <span className="font-semibold">{ticket.username}</span> on{' '}
              {new Date(ticket.created_at).toLocaleDateString()}
            </p>
          </div>
          <button
            onClick={() => setIsEditing(!isEditing)}
            disabled={loading}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-slate-600 text-white rounded-lg transition-colors"
          >
            {isEditing ? 'Cancel' : 'Edit'}
          </button>
        </div>

        <p className="text-slate-300">{ticket.description}</p>

        {isEditing ? (
          <div className="space-y-4 pt-4 border-t border-slate-600">
            <div>
              <label className="block text-sm font-medium text-slate-200 mb-2">Status</label>
              <select
                value={editData.status}
                onChange={(e) => setEditData({ ...editData, status: e.target.value })}
                disabled={loading}
                className="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded-lg text-white"
              >
                {STATUSES.map(s => (
                  <option key={s} value={s}>{s}</option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-slate-200 mb-2">Priority</label>
              <select
                value={editData.priority}
                onChange={(e) => setEditData({ ...editData, priority: e.target.value })}
                disabled={loading}
                className="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded-lg text-white"
              >
                {PRIORITIES.map(p => (
                  <option key={p} value={p}>{p}</option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-slate-200 mb-2">Resolution</label>
              <textarea
                value={editData.resolution}
                onChange={(e) => setEditData({ ...editData, resolution: e.target.value })}
                disabled={loading}
                rows={3}
                placeholder="Enter resolution details..."
                className="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded-lg text-white"
              />
            </div>

            <div className="flex gap-2">
              <button
                onClick={handleSaveChanges}
                disabled={loading}
                className="flex-1 px-4 py-2 bg-green-600 hover:bg-green-700 disabled:bg-slate-600 text-white rounded-lg transition-colors"
              >
                {loading ? 'Saving...' : 'Save Changes'}
              </button>
              <button
                onClick={() => setIsEditing(false)}
                disabled={loading}
                className="flex-1 px-4 py-2 bg-slate-600 hover:bg-slate-500 text-white rounded-lg transition-colors"
              >
                Cancel
              </button>
            </div>
          </div>
        ) : (
          <div className="pt-4 border-t border-slate-600 grid grid-cols-2 gap-4">
            <div>
              <p className="text-sm text-slate-400">Status</p>
              <p className="text-white font-semibold capitalize">{ticket.status}</p>
            </div>
            <div>
              <p className="text-sm text-slate-400">Priority</p>
              <p className="text-white font-semibold capitalize">{ticket.priority}</p>
            </div>
            {ticket.assigned_name && (
              <div>
                <p className="text-sm text-slate-400">Assigned To</p>
                <p className="text-white">{ticket.assigned_name}</p>
              </div>
            )}
            {ticket.resolved_at && (
              <div>
                <p className="text-sm text-slate-400">Resolved</p>
                <p className="text-white">{new Date(ticket.resolved_at).toLocaleDateString()}</p>
              </div>
            )}
            {ticket.resolution && (
              <div className="col-span-2">
                <p className="text-sm text-slate-400">Resolution</p>
                <p className="text-white">{ticket.resolution}</p>
              </div>
            )}
          </div>
        )}
      </div>

      {/* Comments Section */}
      <div className="bg-slate-700/50 border border-slate-600 rounded-lg p-6 space-y-4">
        <h2 className="text-lg font-semibold text-white">Comments ({comments.length})</h2>

        {isLoadingComments ? (
          <p className="text-slate-400">Loading comments...</p>
        ) : comments.length === 0 ? (
          <p className="text-slate-400">No comments yet</p>
        ) : (
          <div className="space-y-3 mb-6 max-h-96 overflow-y-auto">
            {comments.map(comment => (
              <div key={comment.id} className="bg-slate-800/50 rounded-lg p-3">
                <div className="flex justify-between items-start mb-2">
                  <div>
                    <p className="font-semibold text-white">{comment.username}</p>
                    <p className="text-xs text-slate-400">
                      {new Date(comment.created_at).toLocaleString()}
                    </p>
                  </div>
                  {comment.is_internal && (
                    <span className="text-xs px-2 py-1 rounded bg-purple-600 text-purple-200">
                      Internal
                    </span>
                  )}
                </div>
                <p className="text-slate-300">{comment.content}</p>
              </div>
            ))}
          </div>
        )}

        {/* Add Comment Form */}
        <form onSubmit={handleAddComment} className="space-y-3 pt-4 border-t border-slate-600">
          <textarea
            value={newComment}
            onChange={(e) => setNewComment(e.target.value)}
            disabled={loading}
            placeholder="Add a comment..."
            rows={3}
            className="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded-lg text-white"
          />
          <div className="flex items-center justify-between">
            <label className="flex items-center gap-2 text-sm text-slate-300 cursor-pointer">
              <input
                type="checkbox"
                checked={isInternalComment}
                onChange={(e) => setIsInternalComment(e.target.checked)}
                disabled={loading}
                className="w-4 h-4"
              />
              Internal comment (admin only)
            </label>
            <button
              type="submit"
              disabled={loading || !newComment.trim()}
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-slate-600 text-white rounded-lg transition-colors"
            >
              {loading ? 'Posting...' : 'Post Comment'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
