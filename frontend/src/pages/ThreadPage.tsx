import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import api from '../lib/api';
import { getUserLocal } from '../lib/api';
import ReplyList from '../components/ReplyList';

const ThreadPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [thread, setThread] = useState<Record<string, any> | null>(null);
  const [err, setErr] = useState<string | null>(null);
  const [counts, setCounts] = useState<{ upvotes?: number; downvotes?: number }>({});
  const [editing, setEditing] = useState(false);
  const [editBody, setEditBody] = useState('');
  const [editTitle, setEditTitle] = useState('');
  const currentUser = getUserLocal();

  useEffect(() => {
    if (!id) return;
    let mounted = true;
    api.getThread(id).then((data: any) => {
      if (!mounted) return;
      setThread(data);
      setEditBody(data.body || '');
      setEditTitle(data.title || '');
    }).catch((e: any) => setErr(e.message || 'failed to load thread'));
    api.getThreadCounts(id).then((c: any) => { if (mounted) setCounts(c); }).catch(() => {});
    return () => { mounted = false; };
  }, [id]);

  if (err) return <div className="text-red-500">{err}</div>;
  if (!thread) return <div>Loading thread...</div>;

  return (
    <div className="max-w-3xl mx-auto p-6">
      <Link to="/" className="text-sm text-blue-600 hover:underline">← Back</Link>
      <header className="mt-4 mb-4">
        <h1 className="text-4xl font-extrabold mb-2">{thread.title}</h1>
        <div className="text-sm text-gray-500">by <span className="font-medium text-gray-700">{thread.author}</span> • {thread.createdAt ? new Date(thread.createdAt).toLocaleString() : ''}</div>
        {(() => {
          const tags = Array.isArray((thread as any).tags) ? (thread as any).tags : (Array.isArray((thread as any).Tags) ? (thread as any).Tags : []);
          if (!tags || tags.length === 0) return null;
          return (
            <div className="mt-3 flex flex-wrap gap-2">
              {tags.filter((t: string) => !!t).map((t: string) => (
                <span key={t} title={t} className="text-sm text-gray-700 bg-white border border-gray-200 px-2 py-0.5 rounded-full">#{t}</span>
              ))}
            </div>
          );
        })()}
      </header>

  <section className="bg-white p-6 rounded shadow-sm mb-4 relative">
        {editing ? (
          <div className="flex flex-col gap-3">
            <input className="border p-2 rounded" value={editTitle} onChange={e => setEditTitle(e.target.value)} />
            <textarea className="border p-2 rounded" rows={6} value={editBody} onChange={e => setEditBody(e.target.value)} />
            <div className="flex gap-3">
              <button className="px-4 py-2 bg-yellow-400 text-white rounded shadow" onClick={async () => {
                try {
                  await api.updateThread(thread.id, { title: editTitle, body: editBody });
                  setEditing(false);
                  const refreshed: any = await api.getThread(id!);
                  setThread(refreshed);
                } catch (e: any) { alert(e.message || 'update failed'); }
              }}>Save</button>
              <button className="px-4 py-2 bg-gray-100 rounded" onClick={() => { setEditing(false); setEditBody(thread.body); setEditTitle(thread.title); }}>Cancel</button>
            </div>
          </div>
        ) : (
          <div className="prose max-w-none text-lg text-gray-800">{thread.body}</div>
        )}

        {(currentUser && (currentUser.id === thread.user_id || currentUser.role === 'admin')) ? (
          <div className="absolute top-4 right-4 flex items-center gap-3">
            <button className="text-sm text-gray-600 hover:underline bg-transparent p-0 m-0" onClick={() => setEditing(!editing)}>{editing ? 'Editing' : 'Edit'}</button>
            <button className="text-sm text-red-600 hover:underline bg-transparent p-0 m-0" onClick={async () => {
              if (!confirm('Delete this thread?')) return;
              try {
                await api.deleteThread(thread.id);
                window.location.href = '/';
              } catch (e: any) { alert(e.message || 'delete failed'); }
            }}>Delete</button>
          </div>
        ) : null}
      </section>

  <ReplyList threadId={thread.id} />

  <footer className="flex items-center gap-6">
        <div className="flex items-center gap-4 text-lg text-gray-700">
          <button aria-label="Clap" className="flex items-center gap-2" onClick={async () => {
            try {
              const res: any = await api.voteThread(thread.id, 'up');
              if (res.upvotes !== undefined) setCounts({ upvotes: res.upvotes, downvotes: res.downvotes });
            } catch (e: any) { alert(e.message || 'vote failed'); }
          }}>
            <span>👏</span>
            <span className="ml-2 font-medium">{counts.upvotes ?? thread.upvotes ?? 0}</span>
          </button>

          <button aria-label="Down" className="flex items-center gap-2" onClick={async () => {
            try {
              const res: any = await api.voteThread(thread.id, 'down');
              if (res.downvotes !== undefined) setCounts({ upvotes: res.upvotes, downvotes: res.downvotes });
            } catch (e: any) { alert(e.message || 'vote failed'); }
          }}>
            <span>👎</span>
            <span className="ml-2 font-medium">{counts.downvotes ?? thread.downvotes ?? 0}</span>
          </button>
        </div>
      </footer>
    </div>
  );
};

export default ThreadPage;
