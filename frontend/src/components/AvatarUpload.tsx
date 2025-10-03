import React, { useState } from 'react';
import api from '../lib/api';

const AvatarUpload: React.FC<{ userId: string; onUploaded?: (meta: any) => void }> = ({ userId, onUploaded }) => {
  const [file, setFile] = useState<File | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!file) return setError('กรุณาเลือกไฟล์');
    if (file.size > 5 * 1024 * 1024) return setError('ไฟล์ใหญ่เกินไป (max 5MB)');
    setLoading(true);
    setError(null);
    try {
      const resp = await api.uploadAvatar(file, userId);
      if (onUploaded) onUploaded(resp);
    } catch (err: any) {
      setError(err.message || 'upload failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex items-center gap-2">
      <input type="file" accept="image/*" onChange={(e) => setFile(e.target.files ? e.target.files[0] : null)} />
      <button className="px-3 py-1 bg-blue-600 text-white rounded" disabled={loading} type="submit">{loading ? 'อัพโหลด...' : 'อัพโหลด'}</button>
      {error && <div className="text-red-500 text-sm">{error}</div>}
    </form>
  );
};

export default AvatarUpload;
