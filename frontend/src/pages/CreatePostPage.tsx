import React from 'react';
import CreatePost from '../components/CreatePost';

const CreatePostPage: React.FC = () => (
  <div className="p-4">
    <CreatePost onSubmit={(post) => console.log('Post Created:', post)} />
  </div>
);

export default CreatePostPage;