import { Ticket } from '../api/tickets';

interface TicketListProps {
  tickets: Ticket[];
  onSelectTicket: (ticket: Ticket) => void;
  isLoading?: boolean;
}

const priorityColors: Record<string, string> = {
  low: 'bg-blue-500/20 text-blue-200',
  medium: 'bg-yellow-500/20 text-yellow-200',
  high: 'bg-orange-500/20 text-orange-200',
  critical: 'bg-red-500/20 text-red-200',
};

const statusColors: Record<string, string> = {
  open: 'bg-blue-500/20 text-blue-200',
  'in-progress': 'bg-yellow-500/20 text-yellow-200',
  resolved: 'bg-green-500/20 text-green-200',
  closed: 'bg-slate-500/20 text-slate-200',
};

const categoryIcons: Record<string, string> = {
  bug: '🐛',
  'feature-request': '✨',
  other: '📝',
};

export function TicketList({ tickets, onSelectTicket, isLoading }: TicketListProps) {
  if (isLoading) {
    return (
      <div className="text-center py-8">
        <p className="text-slate-400">Loading tickets...</p>
      </div>
    );
  }

  if (tickets.length === 0) {
    return (
      <div className="text-center py-8">
        <p className="text-slate-400">No tickets found</p>
      </div>
    );
  }

  return (
    <div className="space-y-2">
      {tickets.map(ticket => (
        <button
          key={ticket.id}
          onClick={() => onSelectTicket(ticket)}
          className="w-full text-left bg-slate-700/50 hover:bg-slate-600/50 border border-slate-600 rounded-lg p-4 transition-colors"
        >
          <div className="flex items-start justify-between gap-4">
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2 mb-2">
                <span className="text-xl">{categoryIcons[ticket.category] || '📝'}</span>
                <h3 className="font-semibold text-white truncate">#{ticket.id} - {ticket.title}</h3>
              </div>
              <p className="text-sm text-slate-300 line-clamp-2 mb-3">{ticket.description}</p>
              <div className="flex items-center flex-wrap gap-2">
                <span className={`text-xs px-2 py-1 rounded ${priorityColors[ticket.priority] || priorityColors.medium}`}>
                  {ticket.priority}
                </span>
                <span className={`text-xs px-2 py-1 rounded ${statusColors[ticket.status] || statusColors.open}`}>
                  {ticket.status}
                </span>
                {ticket.comment_count > 0 && (
                  <span className="text-xs px-2 py-1 rounded bg-slate-600 text-slate-300">
                    💬 {ticket.comment_count}
                  </span>
                )}
              </div>
            </div>
            <div className="text-right text-xs text-slate-400">
              <p>{new Date(ticket.created_at).toLocaleDateString()}</p>
            </div>
          </div>
        </button>
      ))}
    </div>
  );
}
