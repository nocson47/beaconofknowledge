import React, { useState, useEffect } from 'react';
import api, { getUserLocal } from '../lib/api';

const ReplyItem: React.FC<{ reply: any; onDeleted?: () => void; onUpdated?: (r:any) => void }> = ({ reply, onDeleted, onUpdated }) => {
  const currentUser = getUserLocal();
  const isOwner = currentUser && (currentUser.id === reply.user_id || currentUser.role === 'admin');
  const [editing, setEditing] = useState(false);
  const [body, setBody] = useState(reply.body || '');
  const [username, setUsername] = useState<string | null>(reply.author || null);

  useEffect(() => {
    let mounted = true;
    if (username) return;
    // fallback: fetch user by id to get username
    (async () => {
      try {
        const u: any = await api.getUserByID(reply.user_id);
        if (!mounted) return;
        setUsername((u && u.username) ? u.username : String(reply.user_id));
      } catch (e) {
        if (!mounted) return;
        setUsername(String(reply.user_id));
      }
    })();
    return () => { mounted = false; };
  }, [reply.user_id, reply.author, username]);

  return (
    <div className="border-b py-4">
      <div className="flex items-start justify-between">
        <div>
          <div className="text-sm text-gray-600">{username ?? reply.user_id} • <span className="text-gray-400">{reply.createdAt ? new Date(reply.createdAt).toLocaleString() : ''}</span></div>
          {!editing ? (
            <div className="mt-2 text-gray-800">{reply.body}</div>
          ) : (
            <div className="mt-2">
              <textarea className="w-full border p-2 rounded" value={body} onChange={e => setBody(e.target.value)} />
              <div className="mt-2 flex gap-2">
                <button className="px-3 py-1 bg-yellow-400 text-white rounded" onClick={async () => {
                  try {
                    await api.updateReply(reply.id, { body });
                    setEditing(false);
                    if (onUpdated) onUpdated({ ...reply, body });
                  } catch (e:any) { alert(e.message || 'update failed'); }
                }}>Save</button>
                <button className="px-3 py-1 bg-gray-100 rounded" onClick={() => { setEditing(false); setBody(reply.body); }}>Cancel</button>
              </div>
            </div>
          )}
        </div>
        {isOwner ? (
          <div className="flex gap-2 ml-4">
            {!editing && <button className="text-sm text-yellow-600" onClick={() => setEditing(true)}>Edit</button>}
            <button className="text-sm text-red-600" onClick={async () => {
              if (!confirm('Delete reply?')) return;
              try {
                await api.deleteReply(reply.id);
                if (onDeleted) onDeleted();
              } catch (e:any) { alert(e.message || 'delete failed'); }
            }}>Delete</button>
          </div>
        ) : null}
      </div>
    </div>
  );
}

export default ReplyItem;
