import React from 'react';
import { Link } from 'react-router-dom';

const AdminDashboard: React.FC = () => {
  return (
    <div className="max-w-4xl mx-auto p-6">
      <h1 className="text-2xl font-bold mb-4">Admin Dashboard</h1>
      <div className="space-y-3">
        <Link to="/admin/threads-report" className="block px-4 py-3 bg-white rounded shadow">Threads report</Link>
        <Link to="/admin/users-report" className="block px-4 py-3 bg-white rounded shadow">Users report</Link>
      </div>
    </div>
  );
};

export default AdminDashboard;
