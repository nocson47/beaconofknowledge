import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import api from '../lib/api';

const Profile: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [user, setUser] = useState<Record<string, any> | null>(null);
  const [posts, setPosts] = useState<Array<Record<string, any>> | null>(null);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;
    let mounted = true;
    api.getUserByID(id).then((u: Record<string, any>) => { if (mounted) setUser(u); }).catch((e: any) => setErr(e.message || 'failed to load user'));
    api.getThreads().then((list: any) => {
      const arr: Array<Record<string, any>> = Array.isArray(list) ? list : list.threads || [];
      if (mounted) setPosts(arr.filter(p => String(p.user_id) === String(id)));
    }).catch(() => {});
    return () => { mounted = false; };
  }, [id]);

  if (err) return <div className="text-red-500">{err}</div>;
  if (!user) return <div>Loading user...</div>;

  return (
    <div className="max-w-3xl mx-auto p-4">
      <h1 className="text-2xl font-bold mb-2">{user.username}</h1>
      <div className="text-gray-600 mb-4">{user.email}</div>

      <h2 className="text-xl font-semibold mb-2">Posts by {user.username}</h2>
      <div className="space-y-4">
        {posts ? posts.map(p => (
          <div key={p.id} className="bg-white rounded shadow p-4">
            <Link to={`/threads/${p.id}`} className="text-lg font-semibold text-blue-700 hover:underline">{p.title}</Link>
            <div className="text-gray-600 text-sm line-clamp-3 max-h-20 overflow-hidden">{p.body}</div>
          </div>
        )) : <div>Loading posts...</div>}
      </div>
    </div>
  );
};

export default Profile;
