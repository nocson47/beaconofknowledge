import React from 'react';
import PostList from '../components/PostList';

const HomePage: React.FC = () => (
  <div className="p-4">
    <h1 className="text-3xl font-bold mb-6">Welcome to BeaconOfKnowledge</h1>
    <p className="mb-6">Explore posts, learn, and share knowledge!</p>
    <div>
      <h2 className="text-2xl font-semibold mb-4"></h2>
      <PostList />  
    </div>
  </div>
);

export default HomePage;
