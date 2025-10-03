import React from 'react';
import { Link } from 'react-router-dom';
import api from '../lib/api';
import { useNavigate } from 'react-router-dom';
import { useState } from 'react';

const LoginForm: React.FC = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        try {
            await api.login(username, password);
            navigate('/');
        } catch (err: any) {
            setError(err.message || 'Login failed');
        }
    };

    return (
        <form className="space-y-6" onSubmit={handleSubmit}>
            <div>
                <label htmlFor="member_email" className="block text-sm font-medium text-gray-700">Username / Email</label>
                <input autoFocus id="member_email" name="member[email]" type="text" placeholder="Username or email" value={username} onChange={(e) => setUsername(e.target.value)} className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500" />
            </div>
            <div>
                <label htmlFor="member_password" className="block text-sm font-medium text-gray-700">Password</label>
                <div className="relative">
                    <input id="member_password" name="member[crypted_password]" type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500" />
                    <button type="button" tabIndex={-1} className="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400" aria-label="Toggle password visibility">
                        <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.875 18.825A10.05 10.05 0 0112 19c-5.523 0-10-4.477-10-10 0-1.657.336-3.234.938-4.675M15 12a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
                    </button>
                </div>
            </div>
            <div className="flex justify-between items-center">
                <a href="#" className="text-sm text-primary-600 hover:underline" role="button">Forgot password</a>
            </div>
            <button type="submit" className="w-full py-2 px-4 bg-primary-600 hover:bg-primary-700 text-blue font-semibold rounded-md shadow focus:outline-none">Login</button>
            {error && <div className="text-red-500 mt-2">{error}</div>}
        </form>
    );
};

const RegisterPageAndLogin: React.FC = () => (
    <div className="flex flex-col md:flex-row items-center justify-center min-h-screen bg-gray-50">
        {/* Login Form (preserve original style but handled client-side) */}
        <div className="w-full max-w-md bg-white rounded-lg shadow-md p-8 md:mr-8">
            <LoginForm />
        </div>

        {/* Divider */}
        <div className="my-8 md:my-0 md:mx-8 flex flex-col items-center">
            <span className="text-gray-400 font-semibold">หรือ</span>
        </div>

        {/* Social Login & Register */}
        <div className="w-full max-w-md bg-white rounded-lg shadow-md p-8">
            <div className="space-y-4">
                <button
                    type="button"
                    className="w-full flex items-center justify-center py-2 px-4 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-md"
                >
                        Login with Facebook
                </button>
                <button
                    type="button"
                    className="w-full flex items-center justify-center py-2 px-4 bg-red-500 hover:bg-red-600 text-white font-semibold rounded-md"
                >
                        Login with Google
                </button>
                <button
                    type="button"
                    className="w-full flex items-center justify-center py-2 px-4 bg-green-500 hover:bg-green-600 text-white font-semibold rounded-md"
                >
                    <span className="mr-2">
                        <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                            {/* Line icon */}
                            <circle cx="12" cy="12" r="10" />
                        </svg>
                    </span>
                        Login with Line
                </button>
                <button
                    type="button"
                    className="w-full flex items-center justify-center py-2 px-4 bg-black hover:bg-gray-800 text-white font-semibold rounded-md"
                >
                        Continue with Apple
                </button>
            </div>
            <div className="mt-6 text-center">
                <span className="text-gray-700">
                        Not a member yet?
                    <Link
                        to="/register"
                        className="ml-2 text-primary-600 hover:underline font-semibold"
                    >
                        Register
                    </Link>
                </span>
            </div>
        </div>
    </div>
);

export default RegisterPageAndLogin;