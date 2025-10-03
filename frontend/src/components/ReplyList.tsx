import React, { useEffect, useState } from 'react';
import api from '../lib/api';
import ReplyItem from './ReplyItem';
import { getUserLocal } from '../lib/api';

const ReplyList: React.FC<{ threadId: number | string }> = ({ threadId }) => {
  const [replies, setReplies] = useState<any[] | null>(null);
  const [err, setErr] = useState<string | null>(null);
  const [newBody, setNewBody] = useState('');
  const currentUser = getUserLocal();

  const load = () => {
    api.getRepliesByThread(threadId).then((data:any) => {
      if (!data) return setReplies([]);
      if (Array.isArray(data)) return setReplies(data);
      if (data.replies && Array.isArray(data.replies)) return setReplies(data.replies);
      // fallback: try common shapes
      return setReplies([]);
    }).catch((e:any) => setErr(e?.message || 'failed to load replies'));
  };

  useEffect(() => { load(); }, [threadId]);

  if (err) return <div className="text-red-500">{err}</div>;
  if (replies === null) return <div className="text-gray-500">Loading replies...</div>;

  return (
    <div className="mt-6">
      <h3 className="text-lg font-medium mb-3">Replies</h3>
      {currentUser ? (
        <div className="mb-4">
          <textarea className="w-full border p-2 rounded" rows={3} value={newBody} onChange={e => setNewBody(e.target.value)} />
          <div className="mt-2">
            <button className="px-3 py-1 bg-blue-600 text-white rounded" onClick={async () => {
              if (!newBody.trim()) return alert('reply cannot be empty');
              try {
                await api.createReply({ thread_id: Number(threadId), body: newBody });
                setNewBody('');
                load();
              } catch (e:any) { alert(e.message || 'create failed'); }
            }}>Reply</button>
          </div>
        </div>
      ) : (
        <div className="text-sm text-gray-500">Log in to reply</div>
      )}

      <div className="space-y-2">
        {replies.map(r => (
          <ReplyItem key={r.id} reply={r} onDeleted={() => load()} onUpdated={() => load()} />
        ))}
      </div>
    </div>
  );
}

export default ReplyList;
