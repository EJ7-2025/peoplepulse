'use client';
import { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';
export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const { login } = useAuth();
  const router = useRouter();
  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault(); setError('');
    try {
      const response = await fetch('http://localhost:8080/login', {
        method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ email, password }),
      });
      const data = await response.json();
      if (response.ok) { login(data.token); router.push('/'); } else { setError(data.error || 'Falha no login'); }
    } catch (err) { setError('Erro de conexão com o servidor.'); }
  };
  return (
    <main className="flex min-h-screen flex-col items-center justify-center bg-gray-100">
      <div className="w-full max-w-md p-8 space-y-6 bg-white rounded-xl shadow-lg">
        <h1 className="text-3xl font-bold text-center text-gray-800">Login PeoplePulse</h1>
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label htmlFor="email" className="text-sm font-bold text-gray-600 block">Email</label>
            <input id="email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} className="w-full p-2 border border-gray-300 rounded-md mt-1" required />
          </div>
          <div>
            <label htmlFor="password" className="text-sm font-bold text-gray-600 block">Senha</label>
            <input id="password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} className="w-full p-2 border border-gray-300 rounded-md mt-1" required />
          </div>
          {error && (<p className="text-sm text-red-600 text-center">{error}</p>)}
          <div><button type="submit" className="w-full py-2 px-4 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-md">Entrar</button></div>
        </form>
      </div>
    </main>
  );
}