import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import HomePage from "./pages/Home";
import AboutPage from "./pages/About";
import CreatePostPage from "./pages/CreatePostPage";
import RegisterPageAndLogin from "./pages/RegisterPageAndLogin";
import RegisterPage from "./pages/RegisterPage";
import ThreadPage from "./pages/ThreadPage";
import Profile from "./pages/Profile";
import Header from "./components/Header";

const App: React.FC = () => (
  <Router>
    <Header />
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/about" element={<AboutPage />} />
      <Route path="/search" element={<div>Search Page</div>} />
      <Route path="/create-post" element={<CreatePostPage />} />
      <Route path="/threads/:id" element={<ThreadPage />} />
      <Route path="/users/:id" element={<Profile />} />
      <Route path="/register-login" element={<RegisterPageAndLogin />} />
      <Route path="/register" element={<RegisterPage />} />
    </Routes>
  </Router>
);

export default App;
