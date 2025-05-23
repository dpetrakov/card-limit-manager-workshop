// Example for frontend/src/types/limitRequest.ts
export interface LimitRequestCreate {
  amount: number;
  currency: string;
  justification: string;
  desired_date: string; // Consider using Date type and formatting before API call
}

export interface LimitRequestView extends LimitRequestCreate {
  id: string;
  user_id: string;
  status: string; // Consider creating an enum for RequestStatus
  current_approver_id: string | null;
  created_at: string; // Consider Date
  updated_at: string; // Consider Date
}

export interface ApiError {
  // Based on CLM.json Error schema
  code: number;
  message: string;
}
