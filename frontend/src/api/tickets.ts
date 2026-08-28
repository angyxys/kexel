import { getAuthHeaders } from './auth';

export interface Ticket {
  id: number;
  user_id: number;
  username: string;
  title: string;
  description: string;
  category: string;
  priority: string;
  status: string;
  assigned_to?: number;
  assigned_name?: string;
  resolution?: string;
  resolved_at?: string;
  comment_count: number;
  created_at: string;
  updated_at: string;
}

export interface TicketComment {
  id: number;
  ticket_id: number;
  user_id: number;
  username: string;
  content: string;
  is_internal: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateTicketRequest {
  title: string;
  description: string;
  category?: string;
  priority?: string;
}

export interface UpdateTicketRequest {
  status?: string;
  priority?: string;
  assigned_to?: string;
  resolution?: string;
}

export interface AddCommentRequest {
  content: string;
  is_internal?: boolean;
}

export interface TicketStats {
  total: number;
  open: number;
  in_progress: number;
  resolved: number;
  closed: number;
}

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000/web';

export async function createTicket(request: CreateTicketRequest): Promise<Ticket> {
  const response = await fetch(`${API_BASE_URL}/tickets`, {
    method: 'POST',
    headers: {
      ...getAuthHeaders(),
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(request),
  });

  if (!response.ok) {
    throw new Error('Failed to create ticket');
  }

  return response.json();
}

export async function getTicket(id: number): Promise<Ticket> {
  const response = await fetch(`${API_BASE_URL}/tickets/${id}`, {
    method: 'GET',
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    throw new Error('Failed to fetch ticket');
  }

  return response.json();
}

export async function getUserTickets(): Promise<{ data: Ticket[] }> {
  const response = await fetch(`${API_BASE_URL}/tickets`, {
    method: 'GET',
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    throw new Error('Failed to fetch tickets');
  }

  return response.json();
}

export async function getAllTickets(page: number = 1, pageSize: number = 20): Promise<{ data: Ticket[] }> {
  const response = await fetch(
    `${API_BASE_URL}/tickets/all?page=${page}&page_size=${pageSize}`,
    {
      method: 'GET',
      headers: getAuthHeaders(),
    }
  );

  if (!response.ok) {
    throw new Error('Failed to fetch tickets');
  }

  return response.json();
}

export async function filterTickets(
  status?: string,
  priority?: string,
  category?: string,
  page: number = 1,
  pageSize: number = 20
): Promise<{ data: Ticket[] }> {
  const params = new URLSearchParams();
  if (status) params.append('status', status);
  if (priority) params.append('priority', priority);
  if (category) params.append('category', category);
  params.append('page', page.toString());
  params.append('page_size', pageSize.toString());

  const response = await fetch(`${API_BASE_URL}/tickets/filter?${params}`, {
    method: 'GET',
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    throw new Error('Failed to filter tickets');
  }

  return response.json();
}

export async function updateTicket(id: number, request: UpdateTicketRequest): Promise<{ message: string }> {
  const response = await fetch(`${API_BASE_URL}/tickets/${id}`, {
    method: 'PATCH',
    headers: {
      ...getAuthHeaders(),
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(request),
  });

  if (!response.ok) {
    throw new Error('Failed to update ticket');
  }

  return response.json();
}

export async function deleteTicket(id: number): Promise<{ message: string }> {
  const response = await fetch(`${API_BASE_URL}/tickets/${id}`, {
    method: 'DELETE',
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    throw new Error('Failed to delete ticket');
  }

  return response.json();
}

export async function addComment(id: number, request: AddCommentRequest): Promise<TicketComment> {
  const response = await fetch(`${API_BASE_URL}/tickets/${id}/comments`, {
    method: 'POST',
    headers: {
      ...getAuthHeaders(),
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(request),
  });

  if (!response.ok) {
    throw new Error('Failed to add comment');
  }

  return response.json();
}

export async function getTicketComments(id: number): Promise<{ data: TicketComment[] }> {
  const response = await fetch(`${API_BASE_URL}/tickets/${id}/comments`, {
    method: 'GET',
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    throw new Error('Failed to fetch comments');
  }

  return response.json();
}

export async function getTicketStats(): Promise<TicketStats> {
  const response = await fetch(`${API_BASE_URL}/tickets/stats`, {
    method: 'GET',
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    throw new Error('Failed to fetch ticket stats');
  }

  return response.json();
}
