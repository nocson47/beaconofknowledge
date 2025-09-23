import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../lib/api';

const Register: React.FC = () => {
	const [username, setUsername] = useState('');
	const [password, setPassword] = useState('');
	const [error, setError] = useState<string | null>(null);
	const navigate = useNavigate();

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		try {
			await api.register({ username, password });
			navigate('/register-login');
		} catch (err: any) {
			setError(err.message || 'Register failed');
		}
	};

	return (
		<div className="max-w-md mx-auto p-6 bg-white rounded shadow">
			<h2 className="text-2xl font-semibold mb-4">Register</h2>
			<form onSubmit={handleSubmit} className="space-y-4">
				<input value={username} onChange={(e) => setUsername(e.target.value)} placeholder="Username" required className="w-full p-2 border rounded" />
				<input value={password} onChange={(e) => setPassword(e.target.value)} placeholder="Password" type="password" required className="w-full p-2 border rounded" />
				<button type="submit" className="px-4 py-2 bg-blue-600 text-white rounded">Create Account</button>
				{error && <div className="text-red-500">{error}</div>}
			</form>
		</div>
	);
};

export default Register;
