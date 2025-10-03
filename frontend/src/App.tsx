import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import HomePage from "./pages/Home";
import AboutPage from "./pages/About";
import Feed from "./pages/Feed";
import CreatePostPage from "./pages/CreatePostPage";
import RegisterPageAndLogin from "./pages/RegisterPageAndLogin";
import RegisterPage from "./pages/RegisterPage";
import ForgotPassword from "./pages/ForgotPassword";
import ResetPassword from "./pages/ResetPassword";
import ThreadPage from "./pages/ThreadPage";
import Profile from "./pages/Profile";
import AdminDashboard from "./pages/AdminDashboard";
import ThreadsReport from "./pages/ThreadsReport";
import UsersReport from "./pages/UsersReport";
import SearchResults from "./pages/SearchResults";
import Header from "./components/Header";

const App: React.FC = () => (
  <Router>
    <Header />
    <Routes>
      <Route path="/" element={<HomePage />} />
        <Route path="/feed" element={<Feed />} />
      <Route path="/about" element={<AboutPage />} />
        <Route path="/admin" element={<AdminDashboard />} />
        <Route path="/admin/threads-report" element={<ThreadsReport />} />
        <Route path="/admin/users-report" element={<UsersReport />} />
  <Route path="/search" element={<SearchResults />} />
      <Route path="/create-post" element={<CreatePostPage />} />
      <Route path="/threads/:id" element={<ThreadPage />} />
      <Route path="/users/:id" element={<Profile />} />
      <Route path="/register-login" element={<RegisterPageAndLogin />} />
      <Route path="/register" element={<RegisterPage />} />
  <Route path="/forgot" element={<ForgotPassword />} />
  <Route path="/reset" element={<ResetPassword />} />
    </Routes>
  </Router>
);

export default App;
