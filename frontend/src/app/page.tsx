'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';

interface Kpi {
  id: number;
  title: string;
  value: number;
}

export default function HomePage() {
  const [kpis, setKpis] = useState<Kpi[]>([]);
  const [loading, setLoading] = useState(true);
  const { token, logout } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!token) {
      router.push('/login');
      return;
    }

    fetch('http://localhost:8080/kpis', {
      method: 'GET',
      headers: {
        // A linha crucial que envia o token
        'Authorization': `Bearer ${token}`,
      },
    })
      .then(response => {
        if (!response.ok) {
          throw new Error('Falha ao buscar KPIs ou token inválido');
        }
        return response.json();
      })
      .then(data => {
        setKpis(data || []); // Garante que kpis seja sempre um array
        setLoading(false);
      })
      .catch(error => {
        console.error("Erro:", error);
        logout();
        router.push('/login');
      });
  }, [token, router, logout]);

  if (!token || loading) {
    return <div className="flex min-h-screen items-center justify-center">Carregando...</div>;
  }

  return (
    <main className="flex min-h-screen flex-col items-center p-24 bg-gray-50">
      <div className="z-10 w-full max-w-5xl">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-4xl font-bold text-gray-800">Meu Painel</h1>
          <button 
            onClick={() => { logout(); router.push('/login'); }}
            className="py-2 px-4 bg-red-600 hover:bg-red-700 text-white font-semibold rounded-md">
            Sair
          </button>
        </div>
        <div className="bg-white p-6 rounded-xl shadow-lg border w-full">
          <h3 className="font-semibold text-lg mb-4 text-gray-700">Meus KPIs</h3>
          <div className="space-y-4">
            {kpis.length > 0 ? kpis.map(kpi => (
              <div key={kpi.id}>
                <div className="flex justify-between mb-1">
                  <span className="text-base font-medium text-gray-600">{kpi.title}</span>
                  <span className="text-sm font-medium text-gray-800">{kpi.value}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-4">
                  <div className="bg-blue-600 h-4 rounded-full" style={{ width: `${kpi.value}%` }}></div>
                </div>
              </div>
            )) : <p className="text-gray-500">Nenhum KPI encontrado para este usuário.</p>}
          </div>
        </div>
      </div>
    </main>
  );
}