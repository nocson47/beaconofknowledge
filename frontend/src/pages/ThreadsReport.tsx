import React, { useEffect, useState } from 'react';
import api from '../lib/api';

const ThreadsReport: React.FC = () => {
  const [reports, setReports] = useState<any[]>([]);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    (async () => {
      try {
        const resp: any = await api.getReports ? await api.getReports('thread') : { reports: [] };
        setReports(resp.reports || resp || []);
      } catch (e: any) { setErr(e.message || 'failed to load reports'); }
    })();
  }, []);

  if (err) return <div className="text-red-500">{err}</div>;

  return (
    <div className="max-w-4xl mx-auto p-6">
      <h1 className="text-2xl font-bold mb-4">Threads Report</h1>
      {reports.length === 0 ? (
        <div className="text-gray-600">No thread reports (backend endpoint /reports not implemented)</div>
      ) : (
        <ul className="space-y-3">
          {reports.map(r => (
            <li key={r.id} className="p-3 bg-white rounded shadow">
              <div>Thread ID: {r.target_id}</div>
              <div>Reason: {r.reason}</div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default ThreadsReport;
