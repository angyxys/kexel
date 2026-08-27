import { useState } from 'react';
import { FilterOptions } from '../api/search';

interface FilterPanelProps {
  onFilterChange: (filters: FilterOptions) => void;
  isOpen: boolean;
  onClose: () => void;
}

export function FilterPanel({ onFilterChange, isOpen, onClose }: FilterPanelProps) {
  const [search, setSearch] = useState('');
  const [selectedRoles, setSelectedRoles] = useState<string[]>([]);
  const [banned, setBanned] = useState<boolean | undefined>(undefined);
  const [sortBy, setSortBy] = useState('vrchat_id');
  const [sortOrder, setSortOrder] = useState<'ASC' | 'DESC'>('ASC');

  const roles = ['user', 'vip', 'mod', 'owner'];

  const handleRoleToggle = (role: string) => {
    setSelectedRoles((prev) =>
      prev.includes(role) ? prev.filter((r) => r !== role) : [...prev, role]
    );
  };

  const handleApplyFilters = () => {
    onFilterChange({
      search: search || undefined,
      roles: selectedRoles.length > 0 ? selectedRoles : undefined,
      banned: banned !== undefined ? banned : undefined,
      sort_by: sortBy,
      sort_order: sortOrder,
    });
    onClose();
  };

  const handleResetFilters = () => {
    setSearch('');
    setSelectedRoles([]);
    setBanned(undefined);
    setSortBy('vrchat_id');
    setSortOrder('ASC');
    onFilterChange({});
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
      <div className="bg-slate-800 rounded-lg shadow-xl max-w-md w-full p-6">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-bold text-white">Advanced Filters</h2>
          <button
            onClick={onClose}
            className="text-slate-400 hover:text-white text-2xl leading-none"
          >
            ×
          </button>
        </div>

        <div className="space-y-4">
          {/* Search */}
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-2">Search</label>
            <input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Search by VRChat ID..."
              className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          {/* Roles */}
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-2">Roles</label>
            <div className="space-y-2">
              {roles.map((role) => (
                <label key={role} className="flex items-center text-slate-300 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={selectedRoles.includes(role)}
                    onChange={() => handleRoleToggle(role)}
                    className="w-4 h-4 rounded border-slate-600 text-blue-600 focus:ring-blue-500"
                  />
                  <span className="ml-2 capitalize">{role}</span>
                </label>
              ))}
            </div>
          </div>

          {/* Ban Status */}
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-2">Ban Status</label>
            <div className="space-y-2">
              <label className="flex items-center text-slate-300 cursor-pointer">
                <input
                  type="radio"
                  name="banned"
                  checked={banned === undefined}
                  onChange={() => setBanned(undefined)}
                  className="w-4 h-4"
                />
                <span className="ml-2">All</span>
              </label>
              <label className="flex items-center text-slate-300 cursor-pointer">
                <input
                  type="radio"
                  name="banned"
                  checked={banned === false}
                  onChange={() => setBanned(false)}
                  className="w-4 h-4"
                />
                <span className="ml-2">Active Only</span>
              </label>
              <label className="flex items-center text-slate-300 cursor-pointer">
                <input
                  type="radio"
                  name="banned"
                  checked={banned === true}
                  onChange={() => setBanned(true)}
                  className="w-4 h-4"
                />
                <span className="ml-2">Banned Only</span>
              </label>
            </div>
          </div>

          {/* Sort */}
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-2">Sort By</label>
            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value)}
              className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="vrchat_id">VRChat ID</option>
              <option value="created_at">Created Date</option>
              <option value="updated_at">Updated Date</option>
            </select>
          </div>

          {/* Sort Order */}
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-2">Order</label>
            <div className="flex gap-2">
              <label className="flex-1 flex items-center text-slate-300 cursor-pointer">
                <input
                  type="radio"
                  name="order"
                  checked={sortOrder === 'ASC'}
                  onChange={() => setSortOrder('ASC')}
                  className="w-4 h-4"
                />
                <span className="ml-2">Ascending</span>
              </label>
              <label className="flex-1 flex items-center text-slate-300 cursor-pointer">
                <input
                  type="radio"
                  name="order"
                  checked={sortOrder === 'DESC'}
                  onChange={() => setSortOrder('DESC')}
                  className="w-4 h-4"
                />
                <span className="ml-2">Descending</span>
              </label>
            </div>
          </div>
        </div>

        {/* Buttons */}
        <div className="mt-6 flex gap-2">
          <button
            onClick={handleResetFilters}
            className="flex-1 px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg transition-colors"
          >
            Reset
          </button>
          <button
            onClick={handleApplyFilters}
            className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
          >
            Apply Filters
          </button>
        </div>
      </div>
    </div>
  );
}
