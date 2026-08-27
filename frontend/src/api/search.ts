import axiosInstance from './client';
import { Player } from '../types';

export interface SearchResponse {
  data: Player[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface FilterOptions {
  search?: string;
  roles?: string[];
  banned?: boolean;
  sort_by?: string;
  sort_order?: 'ASC' | 'DESC';
}

export const searchApi = {
  search: async (query: string, page = 1, pageSize = 20): Promise<SearchResponse> => {
    const response = await axiosInstance.get('/web/players/search', {
      params: {
        q: query,
        page,
        page_size: pageSize,
      },
    });
    return response.data;
  },

  filter: async (options: FilterOptions, page = 1, pageSize = 20): Promise<SearchResponse> => {
    const params = new URLSearchParams();

    if (options.search) params.append('search', options.search);
    if (options.roles?.length) params.append('roles', options.roles.join(','));
    if (options.banned !== undefined) params.append('banned', options.banned.toString());
    if (options.sort_by) params.append('sort_by', options.sort_by);
    if (options.sort_order) params.append('sort_order', options.sort_order);

    params.append('page', page.toString());
    params.append('page_size', pageSize.toString());

    const response = await axiosInstance.get(`/web/players/filter?${params.toString()}`);
    return response.data;
  },
};
