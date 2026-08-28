import { getAuthHeaders } from './auth';

export interface BannerInfo {
  id: number;
  name: string;
  type: string;
  title: string;
  description: string;
  image_url: string;
  width: number;
  height: number;
  file_size: number;
  is_active: boolean;
  display_order: number;
  created_at: string;
}

export interface UploadBannerRequest {
  image: File;
  type: string;
  title: string;
  description: string;
}

export interface UpdateBannerRequest {
  title: string;
  description: string;
  display_order: number;
  is_active: boolean;
}

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000/web';

export async function uploadBanner(request: UploadBannerRequest): Promise<BannerInfo> {
  const formData = new FormData();
  formData.append('image', request.image);
  formData.append('type', request.type);
  formData.append('title', request.title);
  formData.append('description', request.description);

  const response = await fetch(`${API_BASE_URL}/banners`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: formData,
  });

  if (!response.ok) {
    throw new Error('Failed to upload banner');
  }

  return response.json();
}

export async function getUserBanners(): Promise<{ data: BannerInfo[] }> {
  const response = await fetch(`${API_BASE_URL}/banners`, {
    method: 'GET',
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    throw new Error('Failed to fetch banners');
  }

  return response.json();
}

export async function getBannersByType(type: string): Promise<{ data: BannerInfo[] }> {
  const response = await fetch(`${API_BASE_URL}/banners/${type}`, {
    method: 'GET',
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    throw new Error('Failed to fetch banners by type');
  }

  return response.json();
}

export async function updateBanner(id: number, request: UpdateBannerRequest): Promise<{ message: string }> {
  const response = await fetch(`${API_BASE_URL}/banners/${id}`, {
    method: 'PATCH',
    headers: {
      ...getAuthHeaders(),
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(request),
  });

  if (!response.ok) {
    throw new Error('Failed to update banner');
  }

  return response.json();
}

export async function deleteBanner(id: number): Promise<{ message: string }> {
  const response = await fetch(`${API_BASE_URL}/banners/${id}`, {
    method: 'DELETE',
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    throw new Error('Failed to delete banner');
  }

  return response.json();
}
