import { getAuthHeaders } from './auth';

export interface PatreonInfo {
  id: number;
  campaign_id: string;
  is_enabled: boolean;
  last_sync_at?: string;
  tier_mapping: Record<string, string>;
  member_count?: number;
  created_at: string;
}

export interface TierMappingRequest {
  tier_mapping: Record<string, string>;
}

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000/web';

export function getPatreonOAuthURL(state: string): string {
  return `${API_BASE_URL}/integrations/patreon/oauth-url?state=${encodeURIComponent(state)}`;
}

export async function handlePatreonOAuthCallback(code: string): Promise<PatreonInfo> {
  const response = await fetch(`${API_BASE_URL}/integrations/patreon/oauth-callback`, {
    method: 'POST',
    headers: {
      ...getAuthHeaders(),
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ code }),
  });

  if (!response.ok) {
    throw new Error('Failed to authenticate with Patreon');
  }

  return response.json();
}

export async function getPatreonIntegration(): Promise<PatreonInfo> {
  const response = await fetch(`${API_BASE_URL}/integrations/patreon`, {
    method: 'GET',
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    throw new Error('Patreon integration not found');
  }

  return response.json();
}

export async function configurePatreonTierMapping(tierMapping: Record<string, string>): Promise<{ message: string }> {
  const response = await fetch(`${API_BASE_URL}/integrations/patreon/tier-mapping`, {
    method: 'PATCH',
    headers: {
      ...getAuthHeaders(),
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ tier_mapping: tierMapping }),
  });

  if (!response.ok) {
    throw new Error('Failed to configure tier mapping');
  }

  return response.json();
}

export async function syncPatreonMembers(): Promise<{ message: string }> {
  const response = await fetch(`${API_BASE_URL}/integrations/patreon/sync`, {
    method: 'POST',
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    throw new Error('Failed to sync Patreon members');
  }

  return response.json();
}

export async function disconnectPatreon(): Promise<{ message: string }> {
  const response = await fetch(`${API_BASE_URL}/integrations/patreon`, {
    method: 'DELETE',
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    throw new Error('Failed to disconnect Patreon');
  }

  return response.json();
}
