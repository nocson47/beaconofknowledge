import React from 'react';
import CreatePost from '../components/CreatePost';
import api from '../lib/api';
import { useNavigate } from 'react-router-dom';

const CreatePostPage: React.FC = () => {
  const navigate = useNavigate();
  const handleSubmit = async (post: any) => {
    try {
  await api.createThread({ title: post.title, body: post.body, tags: post.tags });
      navigate('/');
    } catch (err) {
      console.error(err);
      alert('Create post failed');
    }
  };

  return (
    <div className="p-4">
      <CreatePost onSubmit={handleSubmit} />
    </div>
  );
};

export default CreatePostPage;