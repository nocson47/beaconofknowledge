import React from 'react';
import PostList from '../components/PostList';

const HomePage: React.FC = () => (
  <div className="p-6">
    <div className="max-w-5xl mx-auto">
      <h1 className="text-3xl font-bold mb-4 text-center">Welcome to BeaconOfKnowledge</h1>
      <p className="mb-6 text-center text-gray-600">Explore posts, learn, and share knowledge!</p>
      <div>
        <PostList />
      </div>
    </div>
  </div>
);

export default HomePage;
