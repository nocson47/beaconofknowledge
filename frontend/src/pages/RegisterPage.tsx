import React, { useState } from "react";

export default function Register() {
  const [form, setForm] = useState({
    username: "",
    email: "",
    password: "",
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // 👉 Send form data to your backend here
    console.log(form);
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50">
      <div className="w-full max-w-sm rounded-2xl bg-white p-6 shadow-sm">
        <h1 className="mb-6 text-center text-2xl font-semibold text-gray-800">
          Register
        </h1>

        <form onSubmit={handleSubmit} className="space-y-4">
          <input
            name="username"
            type="text"
            placeholder="Username"
            value={form.username}
            onChange={handleChange}
            required
            className="w-full rounded-lg border border-gray-300 bg-gray-50 p-3 text-gray-900 placeholder-gray-400
                       focus:border-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-400"
          />

          <input
            name="email"
            type="email"
            placeholder="Email"
            value={form.email}
            onChange={handleChange}
            required
            className="w-full rounded-lg border border-gray-300 bg-gray-50 p-3 text-gray-900 placeholder-gray-400
                       focus:border-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-400"
          />

          <input
            name="password"
            type="password"
            placeholder="Password"
            value={form.password}
            onChange={handleChange}
            required
            className="w-full rounded-lg border border-gray-300 bg-gray-50 p-3 text-gray-900 placeholder-gray-400
                       focus:border-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-400"
          />

          <button
            type="submit"
            className="w-full rounded-lg bg-blue-600 px-4 py-3 text-white transition hover:bg-blue-700
                       focus:outline-none focus:ring-2 focus:ring-blue-400"
          >
            Create Account
          </button>
        </form>
      </div>
    </div>
  );
}
