import { Expense, Balance, DebtSettlement, User } from '@/types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

export const MOCK_USERS: User[] = [
  { id: 'u1', name: 'Rafael (Você)', avatar: 'R' },
  { id: 'u2', name: 'Ana Silva', avatar: 'A' },
  { id: 'u3', name: 'Carlos Eduardo', avatar: 'C' },
  { id: 'u4', name: 'Mariana Lima', avatar: 'M' },
];

export const MOCK_EXPENSES: Expense[] = [
  {
    id: 'e1',
    description: 'Jantar no Restô Gourmet',
    amount: 240.0,
    payerId: 'u1',
    payerName: 'Rafael (Você)',
    category: 'food',
    date: '2026-07-27',
    participants: [
      { userId: 'u1', userName: 'Rafael', share: 60.0 },
      { userId: 'u2', userName: 'Ana Silva', share: 60.0 },
      { userId: 'u3', userName: 'Carlos Eduardo', share: 60.0 },
      { userId: 'u4', userName: 'Mariana Lima', share: 60.0 },
    ],
  },
  {
    id: 'e2',
    description: 'Supermercado Mensal da Casa',
    amount: 520.5,
    payerId: 'u2',
    payerName: 'Ana Silva',
    category: 'housing',
    date: '2026-07-25',
    participants: [
      { userId: 'u1', userName: 'Rafael', share: 260.25 },
      { userId: 'u2', userName: 'Ana Silva', share: 260.25 },
    ],
  },
  {
    id: 'e3',
    description: 'Uber Viagem Show',
    amount: 85.0,
    payerId: 'u3',
    payerName: 'Carlos Eduardo',
    category: 'transport',
    date: '2026-07-24',
    participants: [
      { userId: 'u1', userName: 'Rafael', share: 42.5 },
      { userId: 'u3', userName: 'Carlos Eduardo', share: 42.5 },
    ],
  },
];

export const MOCK_BALANCES: Balance[] = [
  { userId: 'u1', userName: 'Rafael (Você)', amount: -80.25 },
  { userId: 'u2', userName: 'Ana Silva', amount: 180.25 },
  { userId: 'u3', userName: 'Carlos Eduardo', amount: -40.0 },
  { userId: 'u4', userName: 'Mariana Lima', amount: -60.0 },
];

export const MOCK_SETTLEMENTS: DebtSettlement[] = [
  { id: 's1', fromUserId: 'u1', fromUserName: 'Rafael (Você)', toUserId: 'u2', toUserName: 'Ana Silva', amount: 80.25 },
  { id: 's2', fromUserId: 'u3', fromUserName: 'Carlos Eduardo', toUserId: 'u2', toUserName: 'Ana Silva', amount: 40.0 },
  { id: 's3', fromUserId: 'u4', fromUserName: 'Mariana Lima', toUserId: 'u2', toUserName: 'Ana Silva', amount: 60.0 },
];

export async function fetchExpenses(): Promise<Expense[]> {
  try {
    const res = await fetch(`${API_BASE_URL}/expenses`);
    if (!res.ok) throw new Error('API offline');
    return await res.json();
  } catch {
    return MOCK_EXPENSES;
  }
}
