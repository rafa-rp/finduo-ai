import { Expense, User, MonthlySummary } from '@/types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

export async function fetchUsers(): Promise<User[]> {
  const res = await fetch(`${API_BASE_URL}/users`);
  if (!res.ok) throw new Error('Failed to fetch users');
  return await res.json();
}

export async function saveUser(user: { id?: string; name: string; salary: number }): Promise<User> {
  const res = await fetch(`${API_BASE_URL}/users`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(user),
  });
  if (!res.ok) throw new Error('Failed to save user');
  return await res.json();
}

export async function fetchSummary(monthStr: string): Promise<MonthlySummary> {
  const res = await fetch(`${API_BASE_URL}/summary?month=${monthStr}`);
  if (!res.ok) throw new Error('Failed to fetch monthly summary');
  return await res.json();
}

export async function createExpense(expense: {
  description: string;
  amount: number;
  date: string;
  category: string;
  payer_id: string;
  is_shared: boolean;
}): Promise<Expense> {
  const res = await fetch(`${API_BASE_URL}/expenses`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(expense),
  });
  if (!res.ok) throw new Error('Failed to create expense');
  return await res.json();
}

export async function deleteExpense(id: string): Promise<void> {
  const res = await fetch(`${API_BASE_URL}/expenses/${id}`, {
    method: 'DELETE',
  });
  if (!res.ok) throw new Error('Failed to delete expense');
}

export async function toggleSettlement(year: number, month: number, isSettled: boolean, settledByID?: string): Promise<void> {
  const res = await fetch(`${API_BASE_URL}/summary/settle`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      year,
      month,
      is_settled: isSettled,
      settled_by_id: settledByID,
    }),
  });
  if (!res.ok) throw new Error('Failed to update settlement status');
}
