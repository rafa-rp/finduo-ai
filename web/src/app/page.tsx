'use client';

import { useState } from 'react';
import { MOCK_EXPENSES, MOCK_BALANCES, MOCK_SETTLEMENTS, MOCK_USERS } from '@/services/api';
import { Expense } from '@/types';

export default function Home() {
  const [expenses, setExpenses] = useState<Expense[]>(MOCK_EXPENSES);
  const [showAddModal, setShowAddModal] = useState(false);
  const [newDesc, setNewDesc] = useState('');
  const [newAmount, setNewAmount] = useState('');
  const [newCategory, setNewCategory] = useState<Expense['category']>('food');
  const [newPayer, setNewPayer] = useState('u1');

  // Calculate totals
  const totalExpenses = expenses.reduce((acc, item) => acc + item.amount, 0);
  const userBalance = MOCK_BALANCES.find((b) => b.userId === 'u1')?.amount || 0;

  const handleAddExpense = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newDesc || !newAmount) return;

    const amountNum = parseFloat(newAmount);
    const payerObj = MOCK_USERS.find((u) => u.id === newPayer);

    const createdExpense: Expense = {
      id: `e-${Date.now()}`,
      description: newDesc,
      amount: amountNum,
      payerId: newPayer,
      payerName: payerObj ? payerObj.name : 'Você',
      category: newCategory,
      date: new Date().toISOString().split('T')[0],
      participants: MOCK_USERS.map((u) => ({
        userId: u.id,
        userName: u.name.split(' ')[0],
        share: amountNum / MOCK_USERS.length,
      })),
    };

    setExpenses([createdExpense, ...expenses]);
    setNewDesc('');
    setNewAmount('');
    setShowAddModal(false);
  };

  const getCategoryBadge = (category: Expense['category']) => {
    const categoriesMap: Record<string, { label: string; color: string }> = {
      food: { label: 'Alimentação', color: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' },
      housing: { label: 'Moradia', color: 'bg-blue-500/10 text-blue-400 border-blue-500/20' },
      transport: { label: 'Transporte', color: 'bg-amber-500/10 text-amber-400 border-amber-500/20' },
      entertainment: { label: 'Lazer', color: 'bg-purple-500/10 text-purple-400 border-purple-500/20' },
      utilities: { label: 'Contas', color: 'bg-cyan-500/10 text-cyan-400 border-cyan-500/20' },
      pet: { label: 'Pet / Dog', color: 'bg-pink-500/10 text-pink-400 border-pink-500/20' },
      travel: { label: 'Viagem', color: 'bg-indigo-500/10 text-indigo-400 border-indigo-500/20' },
      health: { label: 'Saúde', color: 'bg-rose-500/10 text-rose-400 border-rose-500/20' },
      other: { label: 'Outros', color: 'bg-slate-500/10 text-slate-400 border-slate-500/20' },
    };

    const key = (category || 'other').toLowerCase();
    const cat = categoriesMap[key] || categoriesMap.other;
    return (
      <span className={`px-2.5 py-0.5 text-xs rounded-full border ${cat.color} font-medium`}>
        {cat.label}
      </span>
    );
  };

  return (
    <div className="min-h-screen bg-[#0b0f19] text-slate-100 flex flex-col font-sans">
      {/* Header Navigation */}
      <header className="border-b border-slate-800 bg-[#0f172a]/80 backdrop-blur-md sticky top-0 z-30">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-tr from-emerald-500 to-teal-400 flex items-center justify-center shadow-lg shadow-emerald-500/20 font-bold text-slate-950 text-xl">
              F
            </div>
            <div>
              <span className="text-lg font-bold tracking-tight text-white flex items-center gap-2">
                finduo<span className="gradient-text-ai text-sm font-extrabold px-1.5 py-0.5 rounded bg-purple-950/50 border border-purple-800/40">AI</span>
              </span>
              <p className="text-xs text-slate-400">Divisão de Contas Inteligente</p>
            </div>
          </div>

          <div className="flex items-center gap-4">
            <div className="hidden sm:flex items-center gap-2 text-xs bg-slate-900 border border-slate-800 px-3 py-1.5 rounded-full">
              <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
              <span className="text-slate-300">Backend Go: <span className="font-semibold text-emerald-400">Conectado (v1.26.2)</span></span>
            </div>
            <button
              onClick={() => setShowAddModal(true)}
              id="add-expense-header-btn"
              className="bg-gradient-to-r from-emerald-500 to-teal-500 hover:from-emerald-400 hover:to-teal-400 text-slate-950 font-semibold px-4 py-2 rounded-lg shadow-md shadow-emerald-500/10 transition-all flex items-center gap-2 text-sm"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2.5" d="M12 4v16m8-8H4" />
              </svg>
              Nova Despesa
            </button>
          </div>
        </div>
      </header>

      {/* Main Content Container */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 flex-1 w-full space-y-8">

        {/* Welcome Section */}
        <section className="flex flex-col md:flex-row md:items-center justify-between gap-4 bg-gradient-to-r from-slate-900 via-purple-950/20 to-slate-900 p-6 rounded-2xl border border-slate-800 shadow-xl">
          <div>
            <h1 className="text-2xl sm:text-3xl font-extrabold text-white tracking-tight">
              Olá, Rafael 👋
            </h1>
            <p className="text-slate-400 text-sm mt-1">
              Aqui está o resumo das finanças compartilhadas do grupo neste mês.
            </p>
          </div>
          <div className="flex items-center gap-3">
            <span className="text-xs text-purple-300 bg-purple-900/40 border border-purple-700/50 px-3 py-1.5 rounded-lg flex items-center gap-1.5">
              <svg className="w-4 h-4 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
              RAG & Agentes Ativos
            </span>
          </div>
        </section>

        {/* Metric Cards Grid */}
        <section className="grid grid-cols-1 md:grid-cols-3 gap-5">
          {/* Card 1: Total em Despesas */}
          <div className="glass-card glass-card-hover p-5 rounded-xl flex flex-col justify-between">
            <div className="flex justify-between items-start">
              <div>
                <p className="text-xs font-medium text-slate-400 uppercase tracking-wider">Total em Despesas</p>
                <h2 className="text-2xl sm:text-3xl font-bold text-white mt-1">
                  R$ {totalExpenses.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}
                </h2>
              </div>
              <div className="p-2.5 rounded-lg bg-blue-500/10 text-blue-400 border border-blue-500/20">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
            </div>
            <p className="text-xs text-slate-400 mt-4 flex items-center gap-1">
              <span className="text-emerald-400 font-semibold">{expenses.length} transações</span> registradas no grupo
            </p>
          </div>

          {/* Card 2: Seu Saldo Atual */}
          <div className="glass-card glass-card-hover p-5 rounded-xl flex flex-col justify-between">
            <div className="flex justify-between items-start">
              <div>
                <p className="text-xs font-medium text-slate-400 uppercase tracking-wider">Seu Saldo Líquido</p>
                <h2 className={`text-2xl sm:text-3xl font-bold mt-1 ${userBalance >= 0 ? 'text-emerald-400' : 'text-rose-400'}`}>
                  {userBalance >= 0 ? '+' : ''}R$ {userBalance.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}
                </h2>
              </div>
              <div className={`p-2.5 rounded-lg border ${userBalance >= 0 ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' : 'bg-rose-500/10 text-rose-400 border-rose-500/20'}`}>
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M3 6l3 18h12l3-18H3z" />
                </svg>
              </div>
            </div>
            <p className="text-xs text-slate-400 mt-4">
              {userBalance < 0 ? (
                <span className="text-rose-400 font-medium">Você precisa pagar R$ {Math.abs(userBalance).toFixed(2)} para quitar</span>
              ) : (
                <span className="text-emerald-400 font-medium">Você tem a receber do grupo</span>
              )}
            </p>
          </div>

          {/* Card 3: AI Debt Optimizer */}
          <div className="glass-card glass-card-hover p-5 rounded-xl flex flex-col justify-between border-purple-500/20">
            <div className="flex justify-between items-start">
              <div>
                <p className="text-xs font-medium text-purple-300 uppercase tracking-wider flex items-center gap-1">
                  AI Settlement Optimizer
                </p>
                <h2 className="text-2xl sm:text-3xl font-bold text-white mt-1">
                  3 Transações
                </h2>
              </div>
              <div className="p-2.5 rounded-lg bg-purple-500/10 text-purple-400 border border-purple-500/20">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
              </div>
            </div>
            <p className="text-xs text-purple-300/80 mt-4">
              Otimização de grafo ativada: <span className="text-purple-200 font-semibold">2 transferências economizadas</span>
            </p>
          </div>
        </section>

        {/* Detailed Sections: Left (Expenses) / Right (Settlements & AI Insights) */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          
          {/* Left Column: Recent Expenses (2 cols wide on desktop) */}
          <section className="lg:col-span-2 space-y-4">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-bold text-white flex items-center gap-2">
                Despesas Recentes
                <span className="text-xs font-normal bg-slate-800 text-slate-300 px-2 py-0.5 rounded-full">
                  {expenses.length}
                </span>
              </h2>
            </div>

            <div className="space-y-3">
              {expenses.map((item) => (
                <div key={item.id} className="glass-card glass-card-hover p-4 rounded-xl flex items-center justify-between gap-4">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-lg bg-slate-800 border border-slate-700 flex items-center justify-center text-slate-300 font-bold text-sm">
                      {item.payerName[0]}
                    </div>
                    <div>
                      <div className="flex items-center gap-2">
                        <h3 className="font-semibold text-slate-100 text-sm">{item.description}</h3>
                        {getCategoryBadge(item.category)}
                      </div>
                      <p className="text-xs text-slate-400 mt-1">
                        Pago por <span className="text-slate-200 font-medium">{item.payerName}</span> em {item.date}
                      </p>
                    </div>
                  </div>

                  <div className="text-right">
                    <span className="text-base font-bold text-white">
                      R$ {item.amount.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}
                    </span>
                    <p className="text-xs text-slate-400 mt-0.5">
                      R$ {(item.amount / item.participants.length).toFixed(2)} / pessoa
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </section>

          {/* Right Column: Debt Settlements & AI Suggestions */}
          <section className="space-y-6">
            
            {/* Optimized Debt Settlements Card */}
            <div className="glass-card p-5 rounded-xl space-y-4">
              <h2 className="text-base font-bold text-white flex items-center gap-2">
                <svg className="w-5 h-5 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
                </svg>
                Acertos de Contas Sugeridos
              </h2>

              <div className="space-y-3">
                {MOCK_SETTLEMENTS.map((s) => (
                  <div key={s.id} className="p-3 rounded-lg bg-slate-900/60 border border-slate-800/80 flex items-center justify-between text-xs">
                    <div>
                      <p className="text-slate-300">
                        <span className="font-semibold text-rose-400">{s.fromUserName}</span> paga a{' '}
                        <span className="font-semibold text-emerald-400">{s.toUserName}</span>
                      </p>
                    </div>
                    <span className="font-bold text-white bg-slate-800 px-2.5 py-1 rounded border border-slate-700">
                      R$ {s.amount.toFixed(2)}
                    </span>
                  </div>
                ))}
              </div>

              <button 
                id="settle-up-btn"
                className="w-full text-xs font-semibold py-2.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-emerald-400 border border-emerald-500/30 transition-all text-center"
              >
                Registrar Quitação de Dívida
              </button>
            </div>

            {/* AI Insights Widget */}
            <div className="p-5 rounded-xl bg-gradient-to-br from-purple-950/40 via-slate-900 to-indigo-950/30 border border-purple-800/40 shadow-lg space-y-3">
              <div className="flex items-center gap-2">
                <div className="p-1.5 rounded-md bg-purple-500/20 text-purple-300">
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
                  </svg>
                </div>
                <h3 className="text-sm font-bold gradient-text-ai">Finduo AI Insights</h3>
              </div>
              <p className="text-xs text-slate-300 leading-relaxed">
                Detectamos que <span className="text-purple-300 font-medium">Rafael</span> e <span className="text-purple-300 font-medium">Ana Silva</span> concentram 85% dos pagamentos do grupo neste mês. O fechamento otimizado reduz 2 transferências intermediárias via Pix.
              </p>
            </div>

          </section>
        </div>
      </main>

      {/* Footer */}
      <footer className="border-t border-slate-900 bg-[#080c14] py-6 mt-12 text-xs text-slate-500 text-center">
        <div className="max-w-7xl mx-auto px-4">
          <p>© 2026 Finduo AI — Built with Next.js, React & Go 1.26.2</p>
        </div>
      </footer>

      {/* Add Expense Modal */}
      {showAddModal && (
        <div className="fixed inset-0 z-50 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-[#111827] border border-slate-800 rounded-2xl max-w-md w-full p-6 shadow-2xl space-y-5 animate-in fade-in zoom-in-95 duration-150">
            <div className="flex items-center justify-between">
              <h3 className="text-lg font-bold text-white">Adicionar Nova Despesa</h3>
              <button
                onClick={() => setShowAddModal(false)}
                className="text-slate-400 hover:text-white p-1"
              >
                ✕
              </button>
            </div>

            <form onSubmit={handleAddExpense} className="space-y-4">
              <div>
                <label className="block text-xs font-medium text-slate-300 mb-1">Descrição</label>
                <input
                  type="text"
                  placeholder="Ex: Jantar de Sexta"
                  value={newDesc}
                  onChange={(e) => setNewDesc(e.target.value)}
                  className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-emerald-500"
                  required
                />
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-medium text-slate-300 mb-1">Valor Total (R$)</label>
                  <input
                    type="number"
                    step="0.01"
                    placeholder="0.00"
                    value={newAmount}
                    onChange={(e) => setNewAmount(e.target.value)}
                    className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-emerald-500"
                    required
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-slate-300 mb-1">Categoria</label>
                  <select
                    value={newCategory}
                    onChange={(e) => setNewCategory(e.target.value as Expense['category'])}
                    className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-emerald-500"
                  >
                    <option value="food">Alimentação / Mercado (food)</option>
                    <option value="housing">Moradia / Casa (housing)</option>
                    <option value="transport">Transporte / Carro (transport)</option>
                    <option value="entertainment">Lazer (entertainment)</option>
                    <option value="utilities">Contas (utilities)</option>
                    <option value="pet">Pet / Dog (pet)</option>
                    <option value="travel">Viagem (travel)</option>
                    <option value="health">Saúde (health)</option>
                    <option value="other">Outros (other)</option>
                  </select>
                </div>
              </div>

              <div>
                <label className="block text-xs font-medium text-slate-300 mb-1">Quem pagou?</label>
                <select
                  value={newPayer}
                  onChange={(e) => setNewPayer(e.target.value)}
                  className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-emerald-500"
                >
                  {MOCK_USERS.map((u) => (
                    <option key={u.id} value={u.id}>
                      {u.name}
                    </option>
                  ))}
                </select>
              </div>

              <div className="pt-2 flex justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setShowAddModal(false)}
                  className="px-4 py-2 text-xs font-semibold rounded-lg bg-slate-800 text-slate-300 hover:bg-slate-700"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 text-xs font-semibold rounded-lg bg-emerald-500 text-slate-950 hover:bg-emerald-400 font-bold"
                >
                  Salvar Despesa
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
