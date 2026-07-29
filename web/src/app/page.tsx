'use client';

import { useState, useEffect, useCallback } from 'react';
import { fetchUsers, saveUser, fetchSummary, createExpense, deleteExpense, toggleSettlement } from '@/services/api';
import { User, MonthlySummary, Expense } from '@/types';

export default function Home() {
  const [currentMonth, setCurrentMonth] = useState('2026-07');
  const [summary, setSummary] = useState<MonthlySummary | null>(null);
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Modals
  const [showExpenseModal, setShowExpenseModal] = useState(false);
  const [showUserModal, setShowUserModal] = useState(false);

  // New Expense form state
  const [newDesc, setNewDesc] = useState('');
  const [newAmount, setNewAmount] = useState('');
  const [newCategory, setNewCategory] = useState<string>('food');
  const [newPayerId, setNewPayerId] = useState('');
  const [newIsShared, setNewIsShared] = useState(true);

  // New User form state
  const [newUserName, setNewUserName] = useState('');
  const [newUserSalary, setNewUserSalary] = useState('');

  const loadData = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const [uData, sData] = await Promise.all([
        fetchUsers(),
        fetchSummary(currentMonth),
      ]);
      setUsers(uData);
      setSummary(sData);
      if (uData.length > 0 && !newPayerId) {
        setNewPayerId(uData[0].id);
      }
    } catch (err: unknown) {
      console.error(err);
      setError('Não foi possível conectar ao backend Go. Verifique se a API está rodando na porta 8080.');
    } finally {
      setLoading(false);
    }
  }, [currentMonth, newPayerId]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const handleAddUser = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newUserName) return;

    try {
      const salaryNum = parseFloat(newUserSalary) || 0;
      await saveUser({ name: newUserName, salary: salaryNum });
      setNewUserName('');
      setNewUserSalary('');
      setShowUserModal(false);
      await loadData();
    } catch (err) {
      alert('Erro ao salvar participante: ' + err);
    }
  };

  const handleAddExpense = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newDesc || !newAmount || !newPayerId) return;

    try {
      await createExpense({
        description: newDesc,
        amount: parseFloat(newAmount),
        date: `${currentMonth}-15`, // Default to mid-month
        category: newCategory,
        payer_id: newPayerId,
        is_shared: newIsShared,
      });

      setNewDesc('');
      setNewAmount('');
      setShowExpenseModal(false);
      await loadData();
    } catch (err) {
      alert('Erro ao salvar despesa: ' + err);
    }
  };

  const handleDeleteExpense = async (id: string) => {
    if (!confirm('Deseja realmente remover esta despesa?')) return;
    try {
      await deleteExpense(id);
      await loadData();
    } catch (err) {
      alert('Erro ao remover despesa: ' + err);
    }
  };

  const handleToggleSettle = async () => {
    if (!summary) return;
    const [year, month] = currentMonth.split('-').map(Number);
    try {
      await toggleSettlement(year, month, !summary.is_settled);
      await loadData();
    } catch (err) {
      alert('Erro ao atualizar status de quitação: ' + err);
    }
  };

  const getCategoryBadge = (category: string) => {
    const categoriesMap: Record<string, { label: string; color: string }> = {
      food: { label: 'Alimentação', color: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' },
      housing: { label: 'Moradia', color: 'bg-blue-500/10 text-blue-400 border-blue-500/20' },
      transport: { label: 'Transporte', color: 'bg-amber-500/10 text-amber-400 border-amber-500/20' },
      entertainment: { label: 'Lazer', color: 'bg-purple-500/10 text-purple-400 border-purple-500/20' },
      utilities: { label: 'Contas', color: 'bg-cyan-500/10 text-cyan-400 border-cyan-500/20' },
      other: { label: 'Outros', color: 'bg-slate-500/10 text-slate-400 border-slate-500/20' },
    };

    const cat = categoriesMap[category] || categoriesMap.other;
    return (
      <span className={`px-2.5 py-0.5 text-xs rounded-full border ${cat.color} font-medium`}>
        {cat.label}
      </span>
    );
  };

  const getUserName = (id: string) => {
    const u = users.find((item) => item.id === id);
    return u ? u.name : 'Desconhecido';
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
              <p className="text-xs text-slate-400">Divisão de Contas Inteligente (SQLite)</p>
            </div>
          </div>

          <div className="flex items-center gap-3">
            <div className="hidden sm:flex items-center gap-2 text-xs bg-slate-900 border border-slate-800 px-3 py-1.5 rounded-full">
              <span className={`w-2 h-2 rounded-full ${error ? 'bg-rose-500 animate-ping' : 'bg-emerald-400 animate-pulse'}`}></span>
              <span className="text-slate-300">Backend Go: <span className={`font-semibold ${error ? 'text-rose-400' : 'text-emerald-400'}`}>{error ? 'Offline' : 'Conectado'}</span></span>
            </div>

            <button
              onClick={() => setShowUserModal(true)}
              className="bg-slate-800 hover:bg-slate-700 text-slate-200 font-semibold px-3 py-2 rounded-lg border border-slate-700 transition-all text-xs flex items-center gap-1.5"
            >
              <svg className="w-4 h-4 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
              </svg>
              Participantes ({users.length})
            </button>

            <button
              onClick={() => setShowExpenseModal(true)}
              disabled={users.length === 0}
              id="add-expense-header-btn"
              className="bg-gradient-to-r from-emerald-500 to-teal-500 hover:from-emerald-400 hover:to-teal-400 disabled:opacity-50 text-slate-950 font-semibold px-4 py-2 rounded-lg shadow-md shadow-emerald-500/10 transition-all flex items-center gap-2 text-sm"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2.5" d="M12 4v16m8-8H4" />
              </svg>
              Nova Despesa
            </button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 flex-1 w-full space-y-8">
        {error && (
          <div className="bg-rose-950/40 border border-rose-800/50 rounded-xl p-4 text-rose-300 text-sm flex items-center gap-3">
            <svg className="w-6 h-6 text-rose-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <div>
              <p className="font-semibold">{error}</p>
              <p className="text-xs text-rose-400/80 mt-0.5">Certifique-se de executar `go run main.go` no terminal.</p>
            </div>
          </div>
        )}

        {/* Month Selector Banner */}
        <section className="flex flex-col md:flex-row md:items-center justify-between gap-4 bg-gradient-to-r from-slate-900 via-purple-950/20 to-slate-900 p-6 rounded-2xl border border-slate-800 shadow-xl">
          <div>
            <h1 className="text-2xl sm:text-3xl font-extrabold text-white tracking-tight">
              Divisão Proporcional de Despesas
            </h1>
            <p className="text-slate-400 text-sm mt-1">
              {users.length === 0
                ? 'Cadastre os participantes e salários para calcular a divisão proporcional.'
                : `Resumo financeiro calculado com base na renda dos ${users.length} participantes.`}
            </p>
          </div>

          <div className="flex items-center gap-3">
            <label className="text-xs font-medium text-slate-400">Mês:</label>
            <input
              type="month"
              value={currentMonth}
              onChange={(e) => setCurrentMonth(e.target.value)}
              className="bg-slate-900 border border-slate-700 text-white rounded-lg px-3 py-1.5 text-sm focus:outline-none focus:border-emerald-500"
            />
          </div>
        </section>

        {loading ? (
          <div className="py-16 text-center text-slate-400 text-sm animate-pulse">
            Carregando dados do SQLite...
          </div>
        ) : (
          <>
            {/* Metric Cards Grid */}
            <section className="grid grid-cols-1 md:grid-cols-3 gap-5">
              {/* Card 1: Total Compartilhado */}
              <div className="glass-card p-5 rounded-xl flex flex-col justify-between">
                <div className="flex justify-between items-start">
                  <div>
                    <p className="text-xs font-medium text-slate-400 uppercase tracking-wider">Total em Despesas Compartilhadas</p>
                    <h2 className="text-2xl sm:text-3xl font-bold text-white mt-1">
                      R$ {(summary?.total_shared_expenses || 0).toLocaleString('pt-BR', { minimumFractionDigits: 2 })}
                    </h2>
                  </div>
                  <div className="p-2.5 rounded-lg bg-blue-500/10 text-blue-400 border border-blue-500/20">
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>
                </div>
                <p className="text-xs text-slate-400 mt-4">
                  <span className="text-emerald-400 font-semibold">{summary?.expenses?.length || 0} despesas</span> salvas em {currentMonth}
                </p>
              </div>

              {/* Card 2: Acerto / Status das Dívidas */}
              <div className="glass-card p-5 rounded-xl flex flex-col justify-between border-purple-500/20">
                <div className="flex justify-between items-start">
                  <div>
                    <p className="text-xs font-medium text-purple-300 uppercase tracking-wider">Resultado da Divisão</p>
                    <h2 className="text-lg font-bold text-emerald-400 mt-1">
                      {summary?.settlement_message || 'Nenhuma pendência'}
                    </h2>
                  </div>
                  <div className="p-2.5 rounded-lg bg-purple-500/10 text-purple-400 border border-purple-500/20">
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>
                </div>

                <div className="mt-4 flex items-center justify-between">
                  <span className={`text-xs px-2.5 py-1 rounded-full border ${summary?.is_settled ? 'bg-emerald-500/20 text-emerald-300 border-emerald-500/30' : 'bg-amber-500/20 text-amber-300 border-amber-500/30'}`}>
                    {summary?.is_settled ? '✓ Mês Quitado' : 'Em Aberto'}
                  </span>
                  <button
                    onClick={handleToggleSettle}
                    className="text-xs text-purple-300 hover:text-purple-100 underline font-medium"
                  >
                    {summary?.is_settled ? 'Marcar em Aberto' : 'Marcar como Quitado'}
                  </button>
                </div>
              </div>

              {/* Card 3: Participantes & Proporções */}
              <div className="glass-card p-5 rounded-xl flex flex-col justify-between">
                <div className="flex justify-between items-start">
                  <div>
                    <p className="text-xs font-medium text-slate-400 uppercase tracking-wider">Proporções do Casal</p>
                    <div className="mt-2 space-y-1">
                      {summary?.users?.map((u) => (
                        <div key={u.id} className="text-xs flex items-center justify-between gap-2">
                          <span className="text-slate-300 font-medium">{u.name}:</span>
                          <span className="text-purple-300 font-bold">{(u.proportion * 100).toFixed(1)}%</span>
                        </div>
                      ))}
                      {(!summary?.users || summary.users.length === 0) && (
                        <p className="text-xs text-slate-500 italic">Cadastre os participantes</p>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            </section>

            {/* Participantes Breakdown Table */}
            {summary?.users && summary.users.length > 0 && (
              <section className="glass-card p-6 rounded-xl space-y-4">
                <h2 className="text-base font-bold text-white flex items-center gap-2">
                  <svg className="w-5 h-5 text-teal-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                  </svg>
                  Detalhamento por Participante (Proporcional à Renda)
                </h2>

                <div className="overflow-x-auto">
                  <table className="w-full text-xs text-left text-slate-300">
                    <thead className="bg-slate-900/80 text-slate-400 uppercase text-[10px] tracking-wider border-b border-slate-800">
                      <tr>
                        <th className="p-3">Nome</th>
                        <th className="p-3">Salário / Renda</th>
                        <th className="p-3">Proporção</th>
                        <th className="p-3">Pago (Compartilhado)</th>
                        <th className="p-3">Justo (Sua Cota)</th>
                        <th className="p-3">Saldo Final</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-800/60">
                      {summary.users.map((u) => (
                        <tr key={u.id} className="hover:bg-slate-900/40">
                          <td className="p-3 font-semibold text-white">{u.name}</td>
                          <td className="p-3">R$ {u.salary.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}</td>
                          <td className="p-3 font-bold text-purple-300">{(u.proportion * 100).toFixed(1)}%</td>
                          <td className="p-3">R$ {u.total_paid_shared.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}</td>
                          <td className="p-3">R$ {u.fair_share.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}</td>
                          <td className={`p-3 font-bold ${u.balance > 0 ? 'text-rose-400' : u.balance < 0 ? 'text-emerald-400' : 'text-slate-300'}`}>
                            {u.balance > 0
                              ? `Deve R$ ${u.balance.toFixed(2)}`
                              : u.balance < 0
                              ? `Recebe R$ ${Math.abs(u.balance).toFixed(2)}`
                              : 'R$ 0,00 (Quittado)'}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </section>
            )}

            {/* Expenses List Section */}
            <section className="space-y-4">
              <div className="flex items-center justify-between">
                <h2 className="text-lg font-bold text-white flex items-center gap-2">
                  Despesas Registradas no SQLite ({summary?.expenses?.length || 0})
                </h2>
              </div>

              {(!summary?.expenses || summary.expenses.length === 0) ? (
                <div className="glass-card p-12 text-center rounded-xl space-y-3">
                  <div className="w-12 h-12 rounded-full bg-slate-800 flex items-center justify-center mx-auto text-slate-500">
                    💸
                  </div>
                  <p className="text-slate-300 font-medium">Nenhuma despesa cadastrada neste mês.</p>
                  <p className="text-xs text-slate-500">Clique em "Nova Despesa" acima para adicionar a primeira despesa no seu banco de dados!</p>
                </div>
              ) : (
                <div className="space-y-3">
                  {summary.expenses.map((item) => (
                    <div key={item.id} className="glass-card p-4 rounded-xl flex items-center justify-between gap-4">
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-lg bg-slate-800 border border-slate-700 flex items-center justify-center text-slate-300 font-bold text-sm">
                          {getUserName(item.payer_id)[0] || 'D'}
                        </div>
                        <div>
                          <div className="flex items-center gap-2">
                            <h3 className="font-semibold text-slate-100 text-sm">{item.description}</h3>
                            {getCategoryBadge(item.category)}
                            {!item.is_shared && (
                              <span className="text-[10px] px-2 py-0.5 rounded bg-amber-500/10 text-amber-400 border border-amber-500/20">
                                Individual
                              </span>
                            )}
                          </div>
                          <p className="text-xs text-slate-400 mt-1">
                            Pago por <span className="text-slate-200 font-medium">{getUserName(item.payer_id)}</span> em {item.date}
                          </p>
                        </div>
                      </div>

                      <div className="flex items-center gap-4">
                        <div className="text-right">
                          <span className="text-base font-bold text-white">
                            R$ {item.amount.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}
                          </span>
                        </div>

                        <button
                          onClick={() => handleDeleteExpense(item.id)}
                          className="text-slate-500 hover:text-rose-400 p-1 transition-colors"
                          title="Remover despesa"
                        >
                          ✕
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </section>
          </>
        )}
      </main>

      {/* Modal: Adicionar Participante */}
      {showUserModal && (
        <div className="fixed inset-0 z-50 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-[#111827] border border-slate-800 rounded-2xl max-w-md w-full p-6 shadow-2xl space-y-5">
            <div className="flex items-center justify-between">
              <h3 className="text-lg font-bold text-white">Cadastrar Participante</h3>
              <button onClick={() => setShowUserModal(false)} className="text-slate-400 hover:text-white">✕</button>
            </div>

            <form onSubmit={handleAddUser} className="space-y-4">
              <div>
                <label className="block text-xs font-medium text-slate-300 mb-1">Nome</label>
                <input
                  type="text"
                  placeholder="Ex: Rafael"
                  value={newUserName}
                  onChange={(e) => setNewUserName(e.target.value)}
                  className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-emerald-500"
                  required
                />
              </div>

              <div>
                <label className="block text-xs font-medium text-slate-300 mb-1">Salário / Renda Mensal (R$)</label>
                <input
                  type="number"
                  step="0.01"
                  placeholder="Ex: 5000.00"
                  value={newUserSalary}
                  onChange={(e) => setNewUserSalary(e.target.value)}
                  className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-emerald-500"
                  required
                />
                <p className="text-[11px] text-slate-400 mt-1">Usado para calcular a proporção da divisão de contas.</p>
              </div>

              <div className="pt-2 flex justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setShowUserModal(false)}
                  className="px-4 py-2 text-xs font-semibold rounded-lg bg-slate-800 text-slate-300 hover:bg-slate-700"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 text-xs font-semibold rounded-lg bg-emerald-500 text-slate-950 hover:bg-emerald-400 font-bold"
                >
                  Salvar Participante
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Modal: Adicionar Despesa */}
      {showExpenseModal && (
        <div className="fixed inset-0 z-50 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-[#111827] border border-slate-800 rounded-2xl max-w-md w-full p-6 shadow-2xl space-y-5">
            <div className="flex items-center justify-between">
              <h3 className="text-lg font-bold text-white">Adicionar Nova Despesa</h3>
              <button onClick={() => setShowExpenseModal(false)} className="text-slate-400 hover:text-white">✕</button>
            </div>

            <form onSubmit={handleAddExpense} className="space-y-4">
              <div>
                <label className="block text-xs font-medium text-slate-300 mb-1">Descrição</label>
                <input
                  type="text"
                  placeholder="Ex: Supermercado"
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
                    onChange={(e) => setNewCategory(e.target.value)}
                    className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-emerald-500"
                  >
                    <option value="food">Alimentação</option>
                    <option value="housing">Moradia</option>
                    <option value="transport">Transporte</option>
                    <option value="entertainment">Lazer</option>
                    <option value="utilities">Contas</option>
                    <option value="other">Outros</option>
                  </select>
                </div>
              </div>

              <div>
                <label className="block text-xs font-medium text-slate-300 mb-1">Quem Pagou?</label>
                <select
                  value={newPayerId}
                  onChange={(e) => setNewPayerId(e.target.value)}
                  className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-emerald-500"
                  required
                >
                  {users.map((u) => (
                    <option key={u.id} value={u.id}>
                      {u.name} (R$ {u.salary.toFixed(2)})
                    </option>
                  ))}
                </select>
              </div>

              <div className="flex items-center gap-2 pt-1">
                <input
                  type="checkbox"
                  id="is_shared"
                  checked={newIsShared}
                  onChange={(e) => setNewIsShared(e.target.checked)}
                  className="rounded border-slate-700 text-emerald-500 focus:ring-0 bg-slate-900"
                />
                <label htmlFor="is_shared" className="text-xs text-slate-300">
                  Despesa compartilhada pelo casal (proporcional à renda)
                </label>
              </div>

              <div className="pt-2 flex justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setShowExpenseModal(false)}
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
