import React, { useEffect, useState } from 'react';
import { useLocation, Link } from 'react-router-dom';
import api from '../lib/api';

function useQuery() {
  return new URLSearchParams(useLocation().search);
}

const SearchResults: React.FC = () => {
  const q = useQuery().get('q') || '';
  const [results, setResults] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    const doSearch = async () => {
      setLoading(true);
      setErr(null);
      try {
        // simple approach: fetch all threads and filter client-side
        const resp: any = await api.getThreads();
        const list = Array.isArray(resp) ? resp : (resp.threads || []);
        const qlc = q.trim().toLowerCase();
        if (!qlc) {
          setResults(list);
        } else {
          const filtered = list.filter((t: any) => {
            return (t.title || '').toLowerCase().includes(qlc) || (t.body || '').toLowerCase().includes(qlc);
          });
          setResults(filtered);
        }
      } catch (e: any) {
        setErr(e.message || 'failed to search');
      } finally {
        setLoading(false);
      }
    };
    doSearch();
  }, [q]);

  return (
    <div className="max-w-4xl mx-auto p-6">
      <h2 className="text-xl font-semibold mb-4">Search results for “{q}”</h2>
      {loading && <div>Loading…</div>}
      {err && <div className="text-red-500">{err}</div>}
      {!loading && results.length === 0 && <div className="text-gray-600">No results</div>}
      <div className="space-y-4 mt-4">
        {results.map((t: any) => (
          <article key={t.id} className="p-4 bg-white rounded shadow">
            <h3 className="text-lg font-bold"><Link to={`/threads/${t.id}`}>{t.title}</Link></h3>
            <p className="text-sm text-gray-600">by {t.username || t.user_id}</p>
            <p className="mt-2 text-gray-700">{(t.body || '').slice(0, 240)}{(t.body||'').length>240? '…':''}</p>
          </article>
        ))}
      </div>
    </div>
  );
};

export default SearchResults;
