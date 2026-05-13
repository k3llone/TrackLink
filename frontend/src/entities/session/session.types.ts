export type UserRole = "customer" | "admin";

export interface User {
  id: string;
  email: string;
  role: UserRole;
  createdAt: string;
}

export type SessionStatus = "idle" | "loading" | "authenticated" | "guest";

export interface SessionError {
  status: number | null;
  code: string;
  message: string;
}

