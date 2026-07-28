export interface User {
  id: string;
  name: string;
  avatar?: string;
}

export interface ExpenseParticipant {
  userId: string;
  userName: string;
  share: number;
}

export interface Expense {
  id: string;
  description: string;
  amount: number;
  payerId: string;
  payerName: string;
  category: 'food' | 'housing' | 'transport' | 'entertainment' | 'utilities' | 'other';
  date: string;
  participants: ExpenseParticipant[];
}

export interface Balance {
  userId: string;
  userName: string;
  amount: number; // positive = to receive, negative = owes
}

export interface DebtSettlement {
  id: string;
  fromUserId: string;
  fromUserName: string;
  toUserId: string;
  toUserName: string;
  amount: number;
}
