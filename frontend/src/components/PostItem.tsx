import React from 'react';
import type { Post } from '../types/post';

const PostItem: React.FC<{ post: Post }> = ({ post }) => (
  <div className="bg-white rounded shadow p-4 flex flex-col gap-2">
    <div className="flex items-center gap-2">
      {post.tags.map(tag => (
        <span key={tag} className="bg-yellow-200 text-yellow-800 text-xs px-2 py-1 rounded">{tag}</span>
      ))}
    </div>
    <a href="#" className="text-lg font-bold text-blue-700 hover:underline">{post.title}</a>
    <div className="text-gray-600 text-sm">{post.content}</div>
    <div className="flex items-center gap-4 text-xs text-gray-400">
      <span>โดย {post.author}</span>
      <span>{new Date(post.createdAt).toLocaleString()}</span>
      <span>👏 {post.clap} 👎 {post.down}</span>
    </div>
  </div>
);

export default PostItem;