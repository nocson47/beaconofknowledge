// src/types/post.ts
export interface Post {
  id: number;
  title: string;
  content: string;
  author: string;
  createdAt: string;
  tags: string[];
  clap: number;
  down: number;
}