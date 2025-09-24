import React from 'react';
import { Link } from 'react-router-dom';
import type { Post } from '../types/post';
import { getUserLocal } from '../lib/api';
import api from '../lib/api';

const PostItem: React.FC<{ post: Post }> = ({ post }) => {
  const currentUser = getUserLocal();
  const isOwner = currentUser && (currentUser.id === post.user_id || currentUser.role === 'admin');
  return (
    <div className="bg-white rounded-lg shadow p-6 flex flex-col gap-4 max-w-2xl mx-auto">
      <div>
        <Link to={`/threads/${post.id}`} className="text-2xl font-bold text-black hover:underline block mb-2">{post.title}</Link>
        <div className="text-sm text-gray-500 mb-3 flex items-center gap-2">
          <img src={(post as any).avatar_url || '/person.svg'} alt="avatar" className="w-6 h-6 rounded-full object-cover" />
          <div>by <Link to={`/users/${post.user_id}`} className="text-gray-800 font-medium hover:underline">{post.author}</Link> • {post.createdAt ? new Date(post.createdAt).toLocaleDateString() : ''}</div>
        </div>
        {(() => {
          const tags = Array.isArray((post as any).tags) ? (post as any).tags : (Array.isArray((post as any).Tags) ? (post as any).Tags : []);
          if (!tags || tags.length === 0) return null;
          return (
            <div className="mt-2 flex flex-wrap gap-2">
              {tags.filter((t: string) => !!t).map((t: string) => (
                <span key={t} title={t} className="text-sm text-gray-700 bg-white border border-gray-200 px-2 py-0.5 rounded-full">#{t}</span>
              ))}
            </div>
          );
        })()}
        <div className="text-gray-700" style={{ display: '-webkit-box', WebkitLineClamp: 3 as any, WebkitBoxOrient: 'vertical' as any, overflow: 'hidden' }}>{post.body}</div>
      </div>

      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4 text-gray-700">
          <div className="flex items-center gap-2 text-lg">
            <span>👏</span>
            <span className="font-semibold">{post.upvotes ?? 0}</span>
          </div>
          <div className="flex items-center gap-2 text-lg">
            <span>👎</span>
            <span className="font-semibold">{post.downvotes ?? 0}</span>
          </div>
        </div>

        {isOwner ? (
          <div className="relative">
            <div className="absolute top-0 right-0">
              <Link to={`/threads/${post.id}`} className="text-sm text-gray-600 hover:underline">Edit</Link>
            </div>
          </div>
        ) : null}
        {/* report button for non-owners */}
        {!isOwner ? (
          <div className="ml-auto">
            <button className="text-sm text-red-500 hover:underline" onClick={async () => {
              if (!confirm('Report this thread to moderators?')) return;
              try {
                await api.report('thread', post.id, 'User reported from UI');
                alert('Report submitted. Thank you.');
              } catch (e: any) {
                alert(e.message || 'report failed');
              }
            }}>Report</button>
          </div>
        ) : null}
      </div>
    </div>
  );
}

export default PostItem;