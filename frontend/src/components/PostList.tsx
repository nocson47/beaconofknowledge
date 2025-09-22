import React from 'react';
import PostItem from './PostItem';
import type { Post } from '../types/post';

const mockPosts: Post[] = [
  {
    id: 1,
    title: 'ตัวอย่างโพสต์แรก',
    content: 'เนื้อหาของโพสต์แรก',
    author: 'User1',
    createdAt: new Date().toISOString(),
    tags: ['react', 'webboard'],
    clap: 2,
    down: 0,
  },
  {
    id: 2,
    title: 'ตัวอย่างโพสต์ที่สอง',
    content: 'เนื้อหาของโพสต์ที่สอง',
    author: 'User2',
    createdAt: new Date().toISOString(),
    tags: ['typescript', 'frontend'],
    clap: 5,
    down: 1,
  },
];

const PostList: React.FC = () => (
  <div className="space-y-4">
    {mockPosts.map(post => (
      <PostItem key={post.id} post={post} />
    ))}
  </div>
);

export default PostList;