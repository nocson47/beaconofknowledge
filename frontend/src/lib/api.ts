// Minimal API client used by the frontend PoC
const BASE = (import.meta.env.VITE_API_BASE_URL as string) || 'http://localhost:3000';

type LoginResp = { token: string; expires_in?: number; user?: any };

export function setToken(token: string | null) {
  if (token) localStorage.setItem('auth_token', token);
  else localStorage.removeItem('auth_token');
}

export function getToken(): string | null {
  return localStorage.getItem('auth_token');
}

function authHeaders(): Record<string, string> {
  const token = getToken();
  return token ? { Authorization: `Bearer ${token}` } : {} as Record<string,string>;
}

async function request(path: string, opts: RequestInit = {}) {
  const headers = Object.assign({ 'Content-Type': 'application/json' }, opts.headers || {}, authHeaders());
  const res = await fetch(`${BASE}${path}`, Object.assign({}, opts, { headers }));
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || res.statusText);
  }
  const ct = res.headers.get('content-type') || '';
  if (ct.includes('application/json')) return res.json();
  return res.text();
}

export async function login(username: string, password: string): Promise<LoginResp> {
  const body = JSON.stringify({ username, password });
  const resp = await request('/users/login', { method: 'POST', body });
  // assume resp: { token, expires_in, user }
  if (resp.token) setToken(resp.token);
  // If backend did not return user object, fetch it by username
  if (resp.user) {
    localStorage.setItem('auth_user', JSON.stringify(resp.user));
  } else {
    try {
      const user = await request(`/users/username/${encodeURIComponent(username)}`);
      if (user) localStorage.setItem('auth_user', JSON.stringify(user));
    } catch (e) {
      // ignore fetch user failure; token is set so future requests can be authorized
      console.warn('failed to fetch user after login', e);
    }
  }
  // notify listeners in the same tab that auth changed
  try { window.dispatchEvent(new Event('authChanged')); } catch {}
  return resp;
}

export async function register(payload: any) {
  const body = JSON.stringify(payload);
  return request('/users', { method: 'POST', body });
}

export function logout() {
  setToken(null);
  localStorage.removeItem('auth_user');
  try { window.dispatchEvent(new Event('authChanged')); } catch {}
}

export function getUserLocal() {
  try { return JSON.parse(localStorage.getItem('auth_user') || 'null'); } catch { return null; }
}

export async function getThreads() {
  return request('/threads');
}

export async function getThread(id: number | string) {
  return request(`/threads/${id}`);
}

export async function voteThread(thread_id: number | string, value: 'up' | 'down' | 1 | -1) {
  const body = JSON.stringify({ thread_id: Number(thread_id), value });
  return request('/votes', { method: 'POST', body });
}

export async function getThreadCounts(id: number | string) {
  return request(`/threads/${id}/votes`);
}

export async function createThread(payload: any) {
  return request('/threads', { method: 'POST', body: JSON.stringify(payload) });
}

export async function updateThread(id: number | string, payload: any) {
  return request(`/threads/${id}`, { method: 'PUT', body: JSON.stringify(payload) });
}

export async function deleteThread(id: number | string) {
  return request(`/threads/${id}`, { method: 'DELETE' });
}

export async function getUserByID(id: number | string) {
  return request(`/users/${id}`);
}

export async function updateUser(id: number | string, payload: any) {
  return request(`/users/${id}`, { method: 'PUT', body: JSON.stringify(payload) });
}

export async function getRepliesByThread(thread_id: number | string) {
  return request(`/replies/thread/${thread_id}`);
}

export async function createReply(payload: any) {
  return request('/replies', { method: 'POST', body: JSON.stringify(payload) });
}

export async function updateReply(id: number | string, payload: any) {
  return request(`/replies/${id}`, { method: 'PUT', body: JSON.stringify(payload) });
}

export async function deleteReply(id: number | string) {
  return request(`/replies/${id}`, { method: 'DELETE' });
}

export async function uploadAvatar(file: File, user_id: string) {
  const form = new FormData();
  form.append('avatar', file);
  form.append('user_id', user_id);
  // fetch with FormData must not set Content-Type; include Authorization header separately
  const headers = authHeaders();
  const res = await fetch(`${BASE}/debug/avatar`, { method: 'POST', body: form, headers });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function report(kind: 'thread' | 'user', target_id: number | string, reason?: string) {
  const body = JSON.stringify({ kind, target_id: Number(target_id), reason });
  return request('/reports', { method: 'POST', body });
}

export async function getReports(kind?: 'thread' | 'user') {
  const q = kind ? `?kind=${encodeURIComponent(kind)}` : '';
  return request(`/reports${q}`);
}

export async function getMe() {
  return request('/users/me');
}

export default { login, register, logout, getThreads, getThread, createThread, updateThread, deleteThread, voteThread, getThreadCounts, getUserByID, updateUser, uploadAvatar, report, getReports, getMe, getUserLocal, getRepliesByThread, createReply, updateReply, deleteReply };
