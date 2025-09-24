    import React from "react";
    import { Link, useNavigate } from "react-router-dom";
    import SearchBar from "./SearchComponents";
    import { getUserLocal, logout, getUserByID, getMe, getToken } from "../lib/api";

    const Header: React.FC = () => {
    const navigate = useNavigate();

    const handleSearch = (query: string) => {
        if (!query.trim()) return;
        navigate(`/search?q=${encodeURIComponent(query)}`);
    };
    const [user, setUser] = React.useState(getUserLocal());
    const handleLogout = () => { logout(); navigate('/'); };

    React.useEffect(() => {
        const handler = () => setUser(getUserLocal());
        window.addEventListener('authChanged', handler);
        window.addEventListener('storage', handler);

        // if we have a cached user id, refresh from API to ensure we have up-to-date role info
        const cached = getUserLocal();
        if (cached && cached.id) {
            (async () => {
                try {
                    const fresh = await getMe();
                    if (fresh) {
                        try { localStorage.setItem('auth_user', JSON.stringify(fresh)); } catch {}
                        setUser(fresh);
                        try { window.dispatchEvent(new Event('authChanged')); } catch {}
                    }
                } catch (e) {
                    // fallback: try fetch by id
                    try {
                        const byId = await getUserByID(cached.id);
                        if (byId) { localStorage.setItem('auth_user', JSON.stringify(byId)); setUser(byId); try { window.dispatchEvent(new Event('authChanged')); } catch {} }
                    } catch {}
                }
            })();
        }

        return () => {
            window.removeEventListener('authChanged', handler);
            window.removeEventListener('storage', handler);
        };
    }, []);

    const [menuOpen, setMenuOpen] = React.useState(false);

    // helper: try to parse role claim from JWT payload as fallback
    const parseRoleFromToken = (token: string | null) => {
        if (!token) return null;
        try {
            const parts = token.split('.');
            if (parts.length < 2) return null;
            let payload = parts[1];
            // Add padding if necessary
            payload = payload.replace(/-/g, '+').replace(/_/g, '/');
            while (payload.length % 4) payload += '=';
            const decoded = atob(payload);
            const obj = JSON.parse(decoded);
            return obj.role || null;
        } catch (e) {
            return null;
        }
    };

    return (
        <header className="bg-[#2d2640] text-white">
        <nav className="max-w-7xl mx-auto flex items-center justify-between h-16 px-4 sm:px-6 lg:px-8">
            {/* ----- Left: Hamburger + Logo ----- */}
            <div className="flex items-center gap-3">
            <button
                aria-label="Open menu"
                onClick={() => setMenuOpen(true)}
                className="p-2 rounded-md hover:bg-white/10 transition-colors"
            >
                {/* simple hamburger icon */}
                <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                </svg>
            </button>

            <Link
            to="/"
            className="flex items-center gap-2 font-bold text-xl hover:opacity-90 transition-opacity"
            >
            BeaconOfKnowledge
            </Link>
            </div>

            {/* ----- Center: Search Bar ----- */}
            <div className="flex-1 mx-6 hidden sm:flex justify-center">
            <div className="w-full max-w-md">
                <SearchBar
                onSearch={handleSearch}
                />
            </div>
            </div>

            {/* ----- Right: Menu Links ----- */}
                        <div className="flex items-center gap-6 text-sm sm:text-base font-medium">
            <button onClick={() => { const u = getUserLocal(); if (u) { navigate('/create-post'); } else { navigate('/register-login'); } }} className="hover:text-blue-300 transition-colors">
                Create Post
            </button>
                        {(() => {
                            const token = getToken();
                            const tokenRole = parseRoleFromToken(token);
                            const isAdmin = (user && user.role === 'admin') || tokenRole === 'admin';
                            return isAdmin ? (<Link to="/admin" className="hover:text-blue-300 transition-colors">Dashboard</Link>) : null;
                        })()}
            <Link
                to="/about"
                className="hover:text-blue-300 transition-colors"
            >
                Community
            </Link>
                        {/* Keep login spot stable: show login link when not authenticated, otherwise show avatar + name + logout inline */}
                        {!user ? (
                            <Link to="/register-login" className="hover:text-blue-300 transition-colors">Login / Register</Link>
                        ) : (
                            <div className="flex items-center gap-3 hover:text-blue-300 transition-colors">
                                <Link to={`/users/${user.id}`} className="text-sm">{user.username || user.name}</Link>
                                <button onClick={handleLogout} className="text-sm text-red-400 ml-2">Logout</button>
                            </div>
                        )}
                        </div>
        </nav>

        {/* Slide-in menu overlay */}
        <div className={`fixed inset-0 z-40 ${menuOpen ? '' : 'pointer-events-none'}`} aria-hidden={!menuOpen}>
            {/* background dim */}
            <div
                className={`absolute inset-0 bg-black/50 transition-opacity ${menuOpen ? 'opacity-100' : 'opacity-0'}`}
                onClick={() => setMenuOpen(false)}
            />
            {/* side panel */}
            <aside className={`absolute left-0 top-0 h-full w-80 max-w-full transform transition-transform ${menuOpen ? 'translate-x-0' : '-translate-x-full'}`}>
                <div className="h-full bg-[#2d2640] text-white flex flex-col">
                    <div className="flex items-center justify-between p-4 border-b border-white/10">
                        <div className="flex items-center gap-3">
                            <svg className="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M3 12h18"></path></svg>
                            <span className="font-bold">Menu</span>
                        </div>
                        <button aria-label="Close menu" onClick={() => setMenuOpen(false)} className="p-2 rounded-md hover:bg-white/10">
                            <svg className="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M6 18L18 6M6 6l12 12"></path></svg>
                        </button>
                    </div>

                    <nav className="p-4 overflow-auto">
                        {/* Top items (icons + text) */}
                        <ul className="space-y-3">
                            <li className="flex items-center gap-3 px-2 py-3 rounded hover:bg-white/5">
                                {/* Home icon */}
                                <svg className="w-5 h-5 opacity-90" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M3 9.5L12 3l9 6.5V20a1 1 0 0 1-1 1h-5v-6H9v6H4a1 1 0 0 1-1-1V9.5z"></path></svg>
                                <Link to="/" className="text-sm">หน้าแรก</Link>
                            </li>
                            <li className="flex items-center gap-3 px-2 py-3 rounded hover:bg-white/5">
                                {/* Profile icon (uses same RSS placeholder icon) */}
                                <svg className="w-5 h-5 opacity-90" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M4 11a9 9 0 0 1 9 9M4 4a16 16 0 0 1 16 16M5 19a1 1 0 1 0 0 2 1 1 0 0 0 0-2z"></path></svg>
                                {user ? (
                                    <Link to={`/users/${user.id}`} className="text-sm">Profile</Link>
                                ) : (
                                    <Link to="/register-login" className="text-sm">Profile</Link>
                                )}
                            </li>
                            {/* Admin dashboard links */}
                            {user && user.role === 'admin' ? (
                                <li className="mt-2 border-t border-white/10 pt-3">
                                    <div className="text-xs text-white/80 mb-2">Admin</div>
                                    <ul className="space-y-2">
                                            <li><Link to="/admin" className="text-sm hover:underline">Dashboard</Link></li>
                                            <li><Link to="/admin/threads-report" className="text-sm hover:underline">Threads report</Link></li>
                                            <li><Link to="/admin/users-report" className="text-sm hover:underline">Users report</Link></li>
                                        </ul>
                                </li>
                            ) : null}
                        </ul>

                        <hr className="my-4 border-white/10" />

                        <ul className="space-y-2 text-sm opacity-90">
                            <li className="flex items-center gap-3 px-2 py-2 rounded hover:bg-white/5">
                                {/* Gift / points icon */}
                                <svg className="w-5 h-5 opacity-90" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M20 12v7a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2v-7"></path><path d="M12 12V7"></path><path d="M7 7h10"></path></svg>
                                <span>แลกพอยต์</span>
                            </li>
            
                        </ul>

                        <div className="mt-auto p-4 text-xs opacity-80">
                            <p>© 2025 Internet Marketing co., ltd</p>
                        </div>
                    </nav>
                </div>
            </aside>
        </div>
        </header>
    );
    };

    export default Header;
