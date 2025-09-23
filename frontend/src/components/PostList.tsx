import React, { useEffect, useState } from 'react';
import PostItem from './PostItem';
import type { Post } from '../types/post';
import api from '../lib/api';

const PostList: React.FC = () => {
  const [posts, setPosts] = useState<Post[] | null>(null);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    let mounted = true;
    api.getThreads().then((data: any) => {
      if (!mounted) return;
      // assume data is array or { threads }
      const list = Array.isArray(data) ? data : data.threads || [];
      setPosts(list);
    }).catch((e: any) => setErr(e.message || 'failed'));
    return () => { mounted = false; };
  }, []);

  if (err) return <div className="text-red-500">{err}</div>;
  if (!posts) return <div>กำลังโหลด...</div>;

  return (
    <div className="space-y-4">
      {posts.map((post) => (
        <div key={post.id} className="max-w-3xl mx-auto">
          <PostItem post={post} />
        </div>
      ))}
    </div>
  );
};

export default PostList;