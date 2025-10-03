import React, { useEffect, useMemo, useState } from 'react';
import PostItem from './PostItem';
import type { Post } from '../types/post';
import api from '../lib/api';

const PostList: React.FC = () => {
  const [posts, setPosts] = useState<Post[] | null>(null);
  const [err, setErr] = useState<string | null>(null);
  const [selectedTag, setSelectedTag] = useState<string>('All');
  const [sortOrder, setSortOrder] = useState<'newest' | 'oldest'>('newest');

  useEffect(() => {
    let mounted = true;
    api.getThreads().then((data: any) => {
      if (!mounted) return;
      // assume data is array or { threads }
      const list = Array.isArray(data) ? data : data.threads || [];
      // defensive: filter out soft-deleted items if backend hasn't already
      const visible = (list || []).filter((p: any) => !(p.is_deleted === true || p.isDeleted === true));
      setPosts(visible);
    }).catch((e: any) => setErr(e.message || 'failed'));
    return () => { mounted = false; };
  }, []);

  // compute unique tags from posts
  const allTags = useMemo(() => {
    if (!posts) return [] as string[];
    const set = new Set<string>();
    posts.forEach((p: any) => {
      const tags = Array.isArray(p.tags) ? p.tags : (Array.isArray(p.Tags) ? p.Tags : []);
      (tags || []).forEach((t: string) => { if (t) set.add(t); });
    });
    return Array.from(set).sort();
  }, [posts]);

  // filtered + sorted view
  const visiblePosts = useMemo(() => {
    if (!posts) return null;
    let list = posts.slice();
    if (selectedTag && selectedTag !== 'All') {
      list = list.filter((p: any) => {
        const tags = Array.isArray(p.tags) ? p.tags : (Array.isArray(p.Tags) ? p.Tags : []);
        return (tags || []).some((t: string) => t === selectedTag);
      });
    }
    list.sort((a: any, b: any) => {
      const da = a.createdAt ? new Date(a.createdAt).getTime() : 0;
      const db = b.createdAt ? new Date(b.createdAt).getTime() : 0;
      return sortOrder === 'newest' ? db - da : da - db;
    });
    return list;
  }, [posts, selectedTag, sortOrder]);

  if (err) return <div className="text-red-500">{err}</div>;
  if (!posts) return <div>กำลังโหลด...</div>;

  return (
    <div className="space-y-4">
      <div className="max-w-3xl mx-auto flex items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <label className="text-sm text-gray-600">Sort:</label>
          <select value={sortOrder} onChange={e => setSortOrder(e.target.value as any)} className="border rounded p-1">
            <option value="newest">Newest</option>
            <option value="oldest">Oldest</option>
          </select>
        </div>

        <div className="flex items-center gap-2 overflow-x-auto">
          <button onClick={() => setSelectedTag('All')} className={`text-sm px-2 py-1 rounded-full ${selectedTag === 'All' ? 'bg-gray-800 text-white' : 'bg-white border'}`}>All</button>
          {allTags.map(t => (
            <button key={t} onClick={() => setSelectedTag(t)} title={t} className={`text-sm px-2 py-1 rounded-full ${selectedTag === t ? 'bg-gray-800 text-white' : 'bg-white border'}`}>#{t}</button>
          ))}
        </div>
      </div>

      {visiblePosts && visiblePosts.length === 0 ? (
        <div className="max-w-3xl mx-auto text-gray-600">ไม่มีโพสต์สำหรับตัวกรองนี้</div>
      ) : null}

      {visiblePosts && visiblePosts.map((post: any) => (
        <div key={post.id} className="max-w-3xl mx-auto">
          <PostItem post={post} />
        </div>
      ))}
    </div>
  );
};

export default PostList;