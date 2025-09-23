    import React from "react";
    import { Link, useNavigate } from "react-router-dom";
    import SearchBar from "./SearchComponents";
    import { getUserLocal, logout } from "../lib/api";
    import AvatarUpload from "./AvatarUpload";

    const Header: React.FC = () => {
    const navigate = useNavigate();

    const handleSearch = (query: string) => {
        if (!query.trim()) return;
        navigate(`/search?q=${encodeURIComponent(query)}`);
    };
    const user = getUserLocal();
    const handleLogout = () => { logout(); navigate('/'); };

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
            <Link
                to="/create-post"
                className="hover:text-blue-300 transition-colors"
            >
                ตั้งกระทู้
            </Link>
            <Link
                to="/about"
                className="hover:text-blue-300 transition-colors"
            >
                คอมมูนิตี้
            </Link>
            {!user && (
              <Link to="/register-login" className="hover:text-blue-300 transition-colors">เข้าสู่ระบบ / สมัครสมาชิก</Link>
            )}
            {user && (
              <div className="flex items-center gap-4">
                <img src={user.avatar_url || '/vite.svg'} className="h-8 w-8 rounded-full" alt="avatar" />
                <span className="text-sm">{user.username || user.name}</span>
                <AvatarUpload userId={String(user.id)} onUploaded={(m) => console.log('uploaded', m)} />
                <button onClick={handleLogout} className="text-sm text-red-400">Logout</button>
              </div>
            )}
            </div>
        </nav>
        </header>
    );
    };

    export default Header;
