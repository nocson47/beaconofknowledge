import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import api from '../lib/api';

const Profile: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [user, setUser] = useState<Record<string, any> | null>(null);
  const [posts, setPosts] = useState<Array<Record<string, any>> | null>(null);
  const [err, setErr] = useState<string | null>(null);
  const [editing, setEditing] = useState(false);
  const [bio, setBio] = useState('');
  const [social, setSocial] = useState('');
  const [avatarFile, setAvatarFile] = useState<File | null>(null);
  const currentUser = JSON.parse(localStorage.getItem('auth_user') || 'null');

  useEffect(() => {
    if (!id) return;
    let mounted = true;
  api.getUserByID(id).then((u: Record<string, any>) => { if (mounted) { setUser(u); setBio(u.bio || ''); setSocial(u.social || ''); } }).catch((e: any) => setErr(e.message || 'failed to load user'));
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
      <div className="flex items-center gap-4 mb-4">
  <img src={user.avatar_url || '/person.svg'} alt="avatar" className="w-20 h-20 rounded-full object-cover" />
        <div>
          <h1 className="text-2xl font-bold">{user.username}</h1>
          {/* show email only to owner or admin */}
          {currentUser && (String(currentUser.id) === String(id) || currentUser.role === 'admin') ? (
            <div className="text-gray-600">{user.email}</div>
          ) : null}
        </div>
      </div>

      {/* report profile button for non-owner users */}
      {(currentUser && String(currentUser.id) !== String(id)) ? (
        <div className="mb-4">
          <button className="px-3 py-1 bg-red-600 text-white rounded" onClick={async () => {
            if (!confirm('Report this profile to moderators?')) return;
            try {
              await api.report('user', Number(id), 'User profile reported from UI');
              alert('Profile report submitted. Thank you.');
            } catch (e: any) { alert(e.message || 'report failed'); }
          }}>Report profile</button>
        </div>
      ) : null}

      {/* show edit button for owner */}
      {(currentUser && String(currentUser.id) === String(id)) ? (
        <div className="mb-4">
          <button className="px-3 py-1 bg-gray-100 rounded" onClick={() => setEditing(!editing)}>{editing ? 'Cancel' : 'Edit profile'}</button>
        </div>
      ) : null}

      {/* display bio and social when not editing */}
      {!editing ? (
        <div className="bg-white rounded shadow p-4 mb-4">
          {user.bio ? <p className="text-gray-800 mb-2">{user.bio}</p> : null}
          {user.social ? (
            <div className="flex flex-wrap gap-2">
              {(String(user.social) || '').split(',').map((s: string, i: number) => {
                const trimmed = s.trim();
                if (!trimmed) return null;
                const href = trimmed.startsWith('http') ? trimmed : `https://${trimmed}`;
                return <a key={i} href={href} className="text-blue-600 hover:underline" target="_blank" rel="noreferrer">{trimmed}</a>;
              })}
            </div>
          ) : null}
        </div>
      ) : null}

      {editing ? (
        <div className="bg-white rounded shadow p-4 mb-4">
          <label className="block mb-2">Avatar</label>
          <input type="file" accept="image/*" onChange={e => setAvatarFile(e.target.files ? e.target.files[0] : null)} />
          <label className="block mt-3">Bio</label>
          <textarea className="w-full border rounded p-2" value={bio} onChange={e => setBio(e.target.value)} />
          <label className="block mt-3">Social links (comma separated)</label>
          <input className="w-full border rounded p-2" value={social} onChange={e => setSocial(e.target.value)} />
          <div className="mt-3">
            <button className="px-4 py-2 bg-blue-600 text-white rounded" onClick={async () => {
              try {
                let uploadedUrl: string | null = null;
                // upload avatar first if present
                if (avatarFile) {
                  const resp: any = await api.uploadAvatar(avatarFile, String(id));
                  if (resp && resp.url) {
                    uploadedUrl = resp.url;
                    setUser((u:any) => ({ ...(u||{}), avatar_url: resp.url }));
                  }
                }
                // update profile fields (include avatar_url if we uploaded one so it's persisted)
                const payload: any = { bio, social };
                if (uploadedUrl) payload.avatar_url = uploadedUrl;
                await api.updateUser(id!, payload);
                const refreshed: any = await api.getUserByID(id!);
                setUser(refreshed);
                setEditing(false);
                alert('Profile updated');
              } catch (e: any) { alert(e.message || 'update failed'); }
            }}>Save</button>
          </div>
        </div>
      ) : null}

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
