import { useState } from 'react';
import { BannerInfo, UploadBannerRequest, UpdateBannerRequest } from '../api/banners';
import { uploadBanner, updateBanner, deleteBanner } from '../api/banners';

interface BannerManagerProps {
  banners: BannerInfo[];
  onBannerUpload: (banner: BannerInfo) => void;
  onBannerUpdate: (banner: BannerInfo) => void;
  onBannerDelete: (id: number) => void;
}

const BANNER_TYPES = ['banner', 'poster', 'avatar'];

export function BannerManager({ banners, onBannerUpload, onBannerUpdate, onBannerDelete }: BannerManagerProps) {
  const [showUploadForm, setShowUploadForm] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [formData, setFormData] = useState({
    image: null as File | null,
    type: 'banner',
    title: '',
    description: '',
  });

  const [editFormData, setEditFormData] = useState({
    title: '',
    description: '',
    display_order: 0,
    is_active: true,
  });

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setFormData(prev => ({ ...prev, image: file }));
    }
  };

  const handleUpload = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.image) {
      setError('Please select an image');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const request: UploadBannerRequest = {
        image: formData.image,
        type: formData.type,
        title: formData.title,
        description: formData.description,
      };

      const banner = await uploadBanner(request);
      onBannerUpload(banner);

      setFormData({
        image: null,
        type: 'banner',
        title: '',
        description: '',
      });
      setShowUploadForm(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to upload banner');
    } finally {
      setLoading(false);
    }
  };

  const handleEdit = (banner: BannerInfo) => {
    setEditingId(banner.id);
    setEditFormData({
      title: banner.title,
      description: banner.description,
      display_order: banner.display_order,
      is_active: banner.is_active,
    });
  };

  const handleSaveEdit = async (bannerId: number) => {
    setLoading(true);
    setError(null);

    try {
      const request: UpdateBannerRequest = editFormData;
      await updateBanner(bannerId, request);

      const updatedBanner = banners.find(b => b.id === bannerId);
      if (updatedBanner) {
        onBannerUpdate({
          ...updatedBanner,
          ...editFormData,
        });
      }

      setEditingId(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update banner');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (bannerId: number) => {
    if (!confirm('Are you sure you want to delete this banner?')) return;

    setLoading(true);
    setError(null);

    try {
      await deleteBanner(bannerId);
      onBannerDelete(bannerId);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete banner');
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

      <div>
        <button
          onClick={() => setShowUploadForm(!showUploadForm)}
          className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
        >
          {showUploadForm ? 'Cancel' : 'Upload New Banner'}
        </button>
      </div>

      {showUploadForm && (
        <form onSubmit={handleUpload} className="bg-slate-700/50 border border-slate-600 rounded-lg p-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-slate-200 mb-2">Image</label>
            <input
              type="file"
              accept="image/*"
              onChange={handleFileChange}
              disabled={loading}
              className="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded-lg text-white"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-200 mb-2">Type</label>
            <select
              value={formData.type}
              onChange={(e) => setFormData(prev => ({ ...prev, type: e.target.value }))}
              disabled={loading}
              className="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded-lg text-white"
            >
              {BANNER_TYPES.map(type => (
                <option key={type} value={type}>{type}</option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-200 mb-2">Title</label>
            <input
              type="text"
              value={formData.title}
              onChange={(e) => setFormData(prev => ({ ...prev, title: e.target.value }))}
              disabled={loading}
              className="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded-lg text-white"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-200 mb-2">Description</label>
            <textarea
              value={formData.description}
              onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
              disabled={loading}
              rows={3}
              className="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded-lg text-white"
            />
          </div>

          <button
            type="submit"
            disabled={loading || !formData.image}
            className="w-full px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-slate-600 text-white rounded-lg transition-colors"
          >
            {loading ? 'Uploading...' : 'Upload Banner'}
          </button>
        </form>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {banners.map(banner => (
          <div key={banner.id} className="bg-slate-700/50 border border-slate-600 rounded-lg overflow-hidden">
            <div className="aspect-video bg-slate-800 overflow-hidden">
              <img
                src={banner.image_url}
                alt={banner.title}
                className="w-full h-full object-cover"
              />
            </div>

            <div className="p-4 space-y-3">
              {editingId === banner.id ? (
                <div className="space-y-3">
                  <div>
                    <label className="block text-sm font-medium text-slate-200 mb-1">Title</label>
                    <input
                      type="text"
                      value={editFormData.title}
                      onChange={(e) => setEditFormData(prev => ({ ...prev, title: e.target.value }))}
                      disabled={loading}
                      className="w-full px-3 py-1 bg-slate-800 border border-slate-600 rounded text-white text-sm"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-slate-200 mb-1">Description</label>
                    <textarea
                      value={editFormData.description}
                      onChange={(e) => setEditFormData(prev => ({ ...prev, description: e.target.value }))}
                      disabled={loading}
                      rows={2}
                      className="w-full px-3 py-1 bg-slate-800 border border-slate-600 rounded text-white text-sm"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-slate-200 mb-1">Display Order</label>
                    <input
                      type="number"
                      value={editFormData.display_order}
                      onChange={(e) => setEditFormData(prev => ({ ...prev, display_order: parseInt(e.target.value) }))}
                      disabled={loading}
                      className="w-full px-3 py-1 bg-slate-800 border border-slate-600 rounded text-white text-sm"
                    />
                  </div>

                  <label className="flex items-center gap-2 text-sm text-slate-200">
                    <input
                      type="checkbox"
                      checked={editFormData.is_active}
                      onChange={(e) => setEditFormData(prev => ({ ...prev, is_active: e.target.checked }))}
                      disabled={loading}
                      className="w-4 h-4"
                    />
                    Active
                  </label>

                  <div className="flex gap-2">
                    <button
                      onClick={() => handleSaveEdit(banner.id)}
                      disabled={loading}
                      className="flex-1 px-3 py-1 bg-green-600 hover:bg-green-700 disabled:bg-slate-600 text-white rounded text-sm"
                    >
                      {loading ? 'Saving...' : 'Save'}
                    </button>
                    <button
                      onClick={() => setEditingId(null)}
                      disabled={loading}
                      className="flex-1 px-3 py-1 bg-slate-600 hover:bg-slate-500 text-white rounded text-sm"
                    >
                      Cancel
                    </button>
                  </div>
                </div>
              ) : (
                <>
                  <div>
                    <h3 className="font-medium text-white">{banner.title}</h3>
                    <p className="text-sm text-slate-400">{banner.type}</p>
                  </div>

                  {banner.description && (
                    <p className="text-sm text-slate-300 line-clamp-2">{banner.description}</p>
                  )}

                  <div className="flex gap-2 text-sm text-slate-400">
                    <span>{banner.width}x{banner.height}</span>
                    <span>•</span>
                    <span>{(banner.file_size / 1024).toFixed(2)} KB</span>
                  </div>

                  <div className="flex gap-2">
                    <button
                      onClick={() => handleEdit(banner)}
                      disabled={loading}
                      className="flex-1 px-3 py-1 bg-blue-600 hover:bg-blue-700 disabled:bg-slate-600 text-white rounded text-sm"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => handleDelete(banner.id)}
                      disabled={loading}
                      className="flex-1 px-3 py-1 bg-red-600 hover:bg-red-700 disabled:bg-slate-600 text-white rounded text-sm"
                    >
                      Delete
                    </button>
                  </div>
                </>
              )}
            </div>
          </div>
        ))}
      </div>

      {banners.length === 0 && !showUploadForm && (
        <div className="text-center py-12">
          <p className="text-slate-400">No banners yet. Upload your first banner!</p>
        </div>
      )}
    </div>
  );
}
