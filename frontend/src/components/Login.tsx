import React, { useState } from 'react';
import './Login.css'; // You'll need to create this CSS file

const Login: React.FC = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        console.log('Logging in with', { username, password });
    }

    return (
        <div className="login-container">
            <form onSubmit={handleSubmit} className="login-form">
                <h2>เข้าสู่ระบบ</h2>
                <div className="form-group">
                    <input
                        type="text"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        placeholder="ชื่อผู้ใช้"
                        required
                    />
                </div>
                <div className="form-group">
                    <input
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        placeholder="รหัสผ่าน"
                        required
                    />
                </div>
                <button type="submit">เข้าสู่ระบบ</button>
            </form>
        </div>
    );
}

export default Login;