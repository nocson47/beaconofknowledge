// src/types/post.ts
export interface Post {
  id: number;
  title: string;
  user_id?: number;
  // backend returns `body` field for the post content
  body: string;
  author: string;
  createdAt?: string;
  tags?: string[];
  // backend uses upvotes/downvotes
  upvotes?: number;
  downvotes?: number;
}