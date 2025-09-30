import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import api from '../lib/api';

const AdminDashboard: React.FC = () => {
  const [threadsCount, setThreadsCount] = useState<number | null>(null);
  const [usersCount, setUsersCount] = useState<number | null>(null);
  const [reportsCount, setReportsCount] = useState<number | null>(null);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    (async () => {
      try {
  const threads: any = await api.getThreads();
        // get all users endpoint exists at /users (may return array)
        let usersList: any[] = [];
        try { const ul = await fetch((import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000') + '/users'); if (ul.ok) usersList = await ul.json(); } catch {}
        const reps: any = await api.getReports().catch(()=>({reports:[]}));
        setThreadsCount(Array.isArray(threads) ? threads.length : (threads.threads ? threads.threads.length : null));
        setUsersCount(Array.isArray(usersList) ? usersList.length : null);
        const rlist = reps.reports || reps || [];
        setReportsCount(Array.isArray(rlist) ? rlist.length : null);
      } catch (e: any) {
        setErr(e.message || 'failed to load stats');
      }
    })();
  }, []);

  return (
    <div className="max-w-4xl mx-auto p-6">
      <h1 className="text-2xl font-bold mb-4">Admin Dashboard</h1>
      {err && <div className="text-red-500">{err}</div>}
      <div className="grid grid-cols-3 gap-4 mb-6">
        <div className="p-4 bg-white rounded shadow">
          <div className="text-sm text-gray-500">Threads</div>
          <div className="text-2xl font-bold">{threadsCount ?? '—'}</div>
        </div>
        <div className="p-4 bg-white rounded shadow">
          <div className="text-sm text-gray-500">Users</div>
          <div className="text-2xl font-bold">{usersCount ?? '—'}</div>
        </div>
        <div className="p-4 bg-white rounded shadow">
          <div className="text-sm text-gray-500">Reports</div>
          <div className="text-2xl font-bold">{reportsCount ?? '—'}</div>
        </div>
      </div>

      <div className="space-y-3">
        <Link to="/admin/threads-report" className="block px-4 py-3 bg-white rounded shadow">Threads report</Link>
        <Link to="/admin/users-report" className="block px-4 py-3 bg-white rounded shadow">Users report</Link>
      </div>
    </div>
  );
};

export default AdminDashboard;
