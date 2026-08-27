import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';
import { playersApi } from '../api/players';
import { searchApi, FilterOptions } from '../api/search';
import { Player } from '../types';
import { PlayerModal } from '../components/PlayerModal';
import { SearchBar } from '../components/SearchBar';
import { FilterPanel } from '../components/FilterPanel';
import { BanModal } from '../components/BanModal';
import { Navigation } from '../components/Navigation';
import { bansApi } from '../api/bans';

export function Dashboard() {
  const navigate = useNavigate();
  const { user, refreshToken, logout } = useAuthStore();
  const [players, setPlayers] = useState<Player[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [showModal, setShowModal] = useState(false);
  const [selectedPlayer, setSelectedPlayer] = useState<Player | null>(null);

  // Search & Filter
  const [searchQuery, setSearchQuery] = useState('');
  const [showFilters, setShowFilters] = useState(false);
  const [currentFilters, setCurrentFilters] = useState<FilterOptions>({});
  const [page, setPage] = useState(1);
  const [pageSize] = useState(20);
  const [totalPages, setTotalPages] = useState(1);

  // Ban modal
  const [showBanModal, setShowBanModal] = useState(false);
  const [playerToBan, setPlayerToBan] = useState<string | null>(null);

  useEffect(() => {
    if (!user) {
      navigate('/login');
      return;
    }
    loadData();
  }, [user, navigate, page, searchQuery, currentFilters]);

  const loadData = async () => {
    try {
      setLoading(true);
      setError('');

      let data;
      if (searchQuery) {
        // Use search API
        data = await searchApi.search(searchQuery, page, pageSize);
      } else if (Object.keys(currentFilters).length > 0) {
        // Use filter API
        data = await searchApi.filter(currentFilters, page, pageSize);
      } else {
        // Load all players
        const players = await playersApi.listPlayers();
        setPlayers(players);
        setTotalPages(1);
        return;
      }

      setPlayers(data.data);
      setTotalPages(data.total_pages);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load players');
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = async () => {
    try {
      if (refreshToken) {
        await authApi.logout(refreshToken);
      }
    } catch {
      // Ignore error, just logout
    }
    logout();
    navigate('/login');
  };

  const handleDeletePlayer = async (id: string) => {
    if (!confirm('Are you sure you want to delete this player?')) return;

    try {
      await playersApi.deletePlayer(id);
      setPlayers(players.filter((p) => p.vrchat_id !== id));
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to delete player');
    }
  };

  const handlePlayerSaved = () => {
    loadData();
    setShowModal(false);
    setSelectedPlayer(null);
  };

  const handleSearch = (query: string) => {
    setSearchQuery(query);
    setPage(1);
    setCurrentFilters({}); // Clear filters when searching
  };

  const handleFilterChange = (filters: FilterOptions) => {
    setCurrentFilters(filters);
    setPage(1);
    setSearchQuery(''); // Clear search when filtering
  };

  const handleResetSearch = () => {
    setSearchQuery('');
    setCurrentFilters({});
    setPage(1);
  };

  const handleBanClick = (playerID: string) => {
    setPlayerToBan(playerID);
    setShowBanModal(true);
  };

  const handleUnbanClick = async (playerID: string) => {
    if (!confirm('Are you sure you want to unban this player?')) return;

    try {
      await bansApi.unbanPlayer(playerID);
      loadData();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to unban player');
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800">
      <Navigation />

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Error Message */}
        {error && (
          <div className="mb-6 p-4 bg-red-900/30 border border-red-700 rounded-lg text-red-200 text-sm">
            {error}
          </div>
        )}

        {/* Actions */}
        <div className="mb-8">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-semibold text-white">Players</h2>
            <div className="flex gap-2">
              <button
                onClick={loadData}
                className="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg transition-colors"
              >
                Refresh
              </button>
              <button
                onClick={() => {
                  setSelectedPlayer(null);
                  setShowModal(true);
                }}
                className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
              >
                Add Player
              </button>
            </div>
          </div>

          {/* Search Bar */}
          <div className="mb-4">
            <SearchBar onSearch={handleSearch} loading={loading} />
          </div>

          {/* Filter & Reset */}
          <div className="flex gap-2">
            <button
              onClick={() => setShowFilters(true)}
              className="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg text-sm transition-colors"
            >
              🔍 Advanced Filters
            </button>
            {(searchQuery || Object.keys(currentFilters).length > 0) && (
              <button
                onClick={handleResetSearch}
                className="px-4 py-2 bg-red-700 hover:bg-red-600 text-white rounded-lg text-sm transition-colors"
              >
                Clear Filters
              </button>
            )}
          </div>

          {/* Active Filters Display */}
          {(searchQuery || Object.keys(currentFilters).length > 0) && (
            <div className="mt-3 text-sm text-slate-300 bg-slate-800/30 p-3 rounded-lg">
              {searchQuery && <span>Search: <strong>{searchQuery}</strong></span>}
              {currentFilters.roles?.length && (
                <span className="ml-3">Roles: <strong>{currentFilters.roles.join(', ')}</strong></span>
              )}
              {currentFilters.banned !== undefined && (
                <span className="ml-3">Status: <strong>{currentFilters.banned ? 'Banned' : 'Active'}</strong></span>
              )}
            </div>
          )}
        </div>

        {/* Players Table */}
        {loading ? (
          <div className="text-center py-8">
            <p className="text-slate-300">Loading players...</p>
          </div>
        ) : players.length === 0 ? (
          <div className="text-center py-8 bg-slate-800/30 rounded-lg border border-slate-700">
            <p className="text-slate-300">No players yet. Add your first player!</p>
          </div>
        ) : (
          <div className="overflow-x-auto bg-slate-800/50 rounded-lg border border-slate-700">
            <table className="w-full">
              <thead className="bg-slate-800 border-b border-slate-700">
                <tr>
                  <th className="px-6 py-3 text-left text-sm font-semibold text-slate-300">VRChat ID</th>
                  <th className="px-6 py-3 text-left text-sm font-semibold text-slate-300">Roles</th>
                  <th className="px-6 py-3 text-left text-sm font-semibold text-slate-300">Status</th>
                  <th className="px-6 py-3 text-left text-sm font-semibold text-slate-300">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-700">
                {players.map((player) => (
                  <tr key={player.vrchat_id} className="hover:bg-slate-700/50 transition-colors">
                    <td className="px-6 py-4 text-sm text-slate-200 font-mono">{player.vrchat_id}</td>
                    <td className="px-6 py-4 text-sm text-slate-200">
                      <div className="flex gap-2">
                        {player.roles.map((role) => (
                          <span
                            key={role}
                            className={`px-2 py-1 rounded text-xs font-medium ${
                              role === 'owner'
                                ? 'bg-purple-900 text-purple-200'
                                : role === 'mod'
                                  ? 'bg-blue-900 text-blue-200'
                                  : role === 'vip'
                                    ? 'bg-green-900 text-green-200'
                                    : 'bg-slate-700 text-slate-200'
                            }`}
                          >
                            {role}
                          </span>
                        ))}
                      </div>
                    </td>
                    <td className="px-6 py-4 text-sm">
                      {player.is_banned ? (
                        <div>
                          <span className="px-2 py-1 bg-red-900 text-red-200 rounded text-xs font-medium block mb-1">
                            Banned
                          </span>
                          <span className="text-xs text-red-300">⏱️ Expires soon</span>
                        </div>
                      ) : (
                        <span className="px-2 py-1 bg-green-900 text-green-200 rounded text-xs font-medium">
                          Active
                        </span>
                      )}
                    </td>
                    <td className="px-6 py-4 text-sm">
                      <div className="flex gap-2 flex-wrap">
                        <button
                          onClick={() => {
                            setSelectedPlayer(player);
                            setShowModal(true);
                          }}
                          className="text-blue-400 hover:text-blue-300 text-xs"
                        >
                          Edit
                        </button>
                        {player.is_banned ? (
                          <button
                            onClick={() => handleUnbanClick(player.vrchat_id)}
                            className="text-green-400 hover:text-green-300 text-xs"
                          >
                            Unban
                          </button>
                        ) : (
                          <button
                            onClick={() => handleBanClick(player.vrchat_id)}
                            className="text-yellow-400 hover:text-yellow-300 text-xs"
                          >
                            Ban
                          </button>
                        )}
                        <button
                          onClick={() => handleDeletePlayer(player.vrchat_id)}
                          className="text-red-400 hover:text-red-300 text-xs"
                        >
                          Delete
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Pagination */}
          <div className="mt-6 flex justify-between items-center">
            <div className="text-sm text-slate-400">
              Page {page} of {totalPages} ({players.length} entries)
            </div>
            <div className="flex gap-2">
              <button
                onClick={() => setPage(Math.max(1, page - 1))}
                disabled={page === 1}
                className="px-4 py-2 bg-slate-700 hover:bg-slate-600 disabled:opacity-50 text-white rounded-lg transition-colors"
              >
                Previous
              </button>
              <button
                onClick={() => setPage(Math.min(totalPages, page + 1))}
                disabled={page === totalPages}
                className="px-4 py-2 bg-slate-700 hover:bg-slate-600 disabled:opacity-50 text-white rounded-lg transition-colors"
              >
                Next
              </button>
            </div>
          </div>
        )}
      </main>

      {/* Filter Panel Modal */}
      <FilterPanel isOpen={showFilters} onClose={() => setShowFilters(false)} onFilterChange={handleFilterChange} />

      {/* Player Modal */}
      {showModal && (
        <PlayerModal
          player={selectedPlayer}
          onClose={() => {
            setShowModal(false);
            setSelectedPlayer(null);
          }}
          onSaved={handlePlayerSaved}
        />
      )}

      {/* Ban Modal */}
      {showBanModal && playerToBan && (
        <BanModal
          playerID={playerToBan}
          onClose={() => {
            setShowBanModal(false);
            setPlayerToBan(null);
          }}
          onSuccess={() => {
            loadData();
          }}
        />
      )}
    </div>
  );
}

// Import authApi
import { authApi } from '../api/auth';
