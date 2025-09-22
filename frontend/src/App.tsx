import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import HomePage from './pages/Home';
// import CreatePostPage from './pages/CreatePostPage';

const App: React.FC = () => (
  <Router>
    <nav className="bg-gray-100 p-4">
      <Link to="/" className="mr-4 text-blue-600">Home</Link>
      <Link to="/create-post" className="text-blue-600">Create Post</Link>
    </nav>
    <Routes>
      <Route path="/" element={<HomePage />} />
      {/* <Route path="/create-post" element={<CreatePostPage />} /> */}
    </Routes>
  </Router>
);

export default App;