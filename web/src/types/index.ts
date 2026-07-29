export interface User {
  id: string;
  name: string;
  salary: number;
}

export interface Expense {
  id: string;
  description: string;
  amount: number;
  date: string;
  category: 'food' | 'housing' | 'transport' | 'entertainment' | 'utilities' | 'other' | string;
  payer_id: string;
  is_shared: boolean;
}

export interface UserSummary {
  id: string;
  name: string;
  salary: number;
  proportion: number;
  total_paid_shared: number;
  total_paid_individual: number;
  fair_share: number;
  balance: number;
}

export interface MonthlySummary {
  month: string;
  is_settled: boolean;
  settled_by_id?: string | null;
  total_shared_expenses: number;
  settlement_message: string;
  users: UserSummary[];
  category_breakdown: Record<string, number>;
  expenses: Expense[];
}
