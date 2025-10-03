    import React from "react";
    import { Link, useNavigate } from "react-router-dom";
    import SearchBar from "./SearchComponents";
    import { getUserLocal, logout } from "../lib/api";

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
        return () => {
            window.removeEventListener('authChanged', handler);
            window.removeEventListener('storage', handler);
        };
    }, []);

    return (
        <header className="bg-[#2d2640] text-white">
        <nav className="max-w-7xl mx-auto flex items-center justify-between h-16 px-4 sm:px-6 lg:px-8">
            {/* ----- Left: Logo ----- */}
            <Link
            to="/"
            className="flex items-center gap-2 font-bold text-xl hover:opacity-90 transition-opacity"
            >
            {/* ถ้ามีโลโก้เป็นรูป ใส่ <img src="/logo.svg" alt="logo" className="h-6" /> */}
            BeaconOfKnowledge
            </Link>

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
        </header>
    );
    };

    export default Header;
