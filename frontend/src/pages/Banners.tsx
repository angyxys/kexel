import { useState, useEffect } from 'react';
import { Navigation } from '../components/Navigation';
import { BannerManager } from '../components/BannerManager';
import { BannerInfo, getUserBanners } from '../api/banners';

export function Banners() {
  const [banners, setBanners] = useState<BannerInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchBanners();
  }, []);

  const fetchBanners = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await getUserBanners();
      setBanners(response.data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch banners');
    } finally {
      setLoading(false);
    }
  };

  const handleBannerUpload = (banner: BannerInfo) => {
    setBanners(prev => [banner, ...prev]);
  };

  const handleBannerUpdate = (updatedBanner: BannerInfo) => {
    setBanners(prev =>
      prev.map(b => b.id === updatedBanner.id ? updatedBanner : b)
    );
  };

  const handleBannerDelete = (id: number) => {
    setBanners(prev => prev.filter(b => b.id !== id));
  };

  if (loading) {
    return (
      <>
        <Navigation />
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="text-center">
            <p className="text-slate-400">Loading banners...</p>
          </div>
        </div>
      </>
    );
  }

  return (
    <>
      <Navigation />
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-white">Banners & Posters</h1>
          <p className="text-slate-400 mt-2">Manage your custom banners, posters, and avatars</p>
        </div>

        <BannerManager
          banners={banners}
          onBannerUpload={handleBannerUpload}
          onBannerUpdate={handleBannerUpdate}
          onBannerDelete={handleBannerDelete}
        />
      </div>
    </>
  );
}
