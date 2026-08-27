import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { playerSchema, PlayerFormData } from '../schemas';
import { playersApi } from '../api/players';
import { Player } from '../types';

interface PlayerModalProps {
  player: Player | null;
  onClose: () => void;
  onSaved: () => void;
}

export function PlayerModal({ player, onClose, onSaved }: PlayerModalProps) {
  const [error, setError] = useState<string>('');

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    watch,
  } = useForm<PlayerFormData>({
    resolver: zodResolver(playerSchema),
    defaultValues: player || {
      vrchat_id: '',
      roles: ['user'],
      is_banned: false,
    },
  });

  const roles = watch('roles');
  const is_banned = watch('is_banned');

  const onSubmit = async (data: PlayerFormData) => {
    try {
      setError('');

      if (player) {
        // Update
        await playersApi.updatePlayer(player.vrchat_id, {
          roles: data.roles,
          is_banned: data.is_banned,
        });
      } else {
        // Create
        await playersApi.createPlayer(data);
      }

      onSaved();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to save player');
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
      <div className="bg-slate-800 rounded-lg shadow-xl max-w-md w-full p-6">
        <h2 className="text-xl font-bold text-white mb-4">
          {player ? 'Edit Player' : 'Add New Player'}
        </h2>

        {error && (
          <div className="mb-4 p-3 bg-red-900/30 border border-red-700 rounded text-red-200 text-sm">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-slate-300 mb-2">VRChat ID</label>
            <input
              {...register('vrchat_id')}
              type="text"
              disabled={!!player}
              className="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
              placeholder="usr_..."
            />
            {errors.vrchat_id && (
              <p className="mt-1 text-xs text-red-400">{errors.vrchat_id.message}</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-300 mb-2">Roles</label>
            <div className="space-y-2">
              {(['owner', 'mod', 'vip', 'user'] as const).map((role) => (
                <label key={role} className="flex items-center text-slate-300 cursor-pointer">
                  <input
                    type="checkbox"
                    value={role}
                    {...register('roles')}
                    className="w-4 h-4 rounded border-slate-600 text-blue-600 focus:ring-blue-500"
                  />
                  <span className="ml-2 capitalize">{role}</span>
                </label>
              ))}
            </div>
            {errors.roles && (
              <p className="mt-1 text-xs text-red-400">{errors.roles.message}</p>
            )}
          </div>

          <div>
            <label className="flex items-center text-slate-300 cursor-pointer">
              <input
                type="checkbox"
                {...register('is_banned')}
                className="w-4 h-4 rounded border-slate-600 text-red-600 focus:ring-red-500"
              />
              <span className="ml-2">Ban this player</span>
            </label>
          </div>

          <div className="flex gap-2 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isSubmitting}
              className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-slate-600 text-white rounded-lg transition-colors"
            >
              {isSubmitting ? 'Saving...' : 'Save'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
