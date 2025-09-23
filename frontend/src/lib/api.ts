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
  const resp = await request('/login', { method: 'POST', body });
  // assume resp: { token, expires_in, user }
  if (resp.token) setToken(resp.token);
  if (resp.user) localStorage.setItem('auth_user', JSON.stringify(resp.user));
  return resp;
}

export async function register(payload: any) {
  const body = JSON.stringify(payload);
  return request('/register', { method: 'POST', body });
}

export function logout() {
  setToken(null);
  localStorage.removeItem('auth_user');
}

export function getUserLocal() {
  try { return JSON.parse(localStorage.getItem('auth_user') || 'null'); } catch { return null; }
}

export async function getThreads() {
  return request('/threads');
}

export async function createThread(payload: any) {
  return request('/threads', { method: 'POST', body: JSON.stringify(payload) });
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

export default { login, register, logout, getThreads, createThread, uploadAvatar, getUserLocal };
