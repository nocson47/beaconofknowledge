import React from 'react';
import { Link } from 'react-router-dom';
import PostList from '../components/PostList';

const HomePage: React.FC = () => (
  <div className="p-4">
    <header className="bg-gray-100 p-4 mb-6 rounded shadow">
      <h1 className="text-3xl font-bold text-center">BeaconOfKnowledge!</h1>
      <div className="text-center mt-4">
        <Link to="/create-post" className="bg-blue-600 text-white px-4 py-2 rounded">
          Create Post
        </Link>
      </div>
    </header>
    <PostList/>
  </div>
);

export default HomePage;