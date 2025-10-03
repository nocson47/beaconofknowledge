import React, { useState } from 'react';
import api from '../lib/api';
import { useNavigate } from 'react-router-dom';

const ForgotPassword: React.FC = () => {
  const [email, setEmail] = useState('');
  const [message, setMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setMessage(null);
    try {
      await api.requestPasswordReset(email, window.location.origin);
      // Always show a generic success message to avoid leaking user existence
      setMessage('If an account exists for that email, a password reset link has been sent.');
      setEmail('');
    } catch (err: any) {
      setError(err.message || 'Request failed');
    }
  };

  return (
    <div className="auth-page">
      <form onSubmit={handleSubmit} className="auth-form">
        <h2>Forgot Password</h2>
        <p>Enter the email associated with your account. We'll send a reset link if an account exists.</p>
        <div className="form-group">
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="Your email"
            required
          />
        </div>
        <button type="submit">Send reset link</button>
        {message && <div className="text-green-600 mt-2">{message}</div>}
        {error && <div className="text-red-600 mt-2">{error}</div>}
        <div className="mt-4">
          <button type="button" onClick={() => navigate('/register-login')} className="link">Back to login</button>
        </div>
      </form>
    </div>
  );
};

export default ForgotPassword;
