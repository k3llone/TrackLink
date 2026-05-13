import { http } from "@/api/http";
import type { User } from "@/entities/session/session.types";

export type { User, UserRole } from "@/entities/session/session.types";

export interface RegisterRequest {
  email: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterResponse {
  user: User;
}

export interface AuthResponse {
  user: User;
}

export const registerUser = (payload: RegisterRequest) => http.post<RegisterResponse>("/auth/register", payload);

export const loginUser = (payload: LoginRequest) => http.post<AuthResponse>("/auth/login", payload);

export const getCurrentUser = () => http.get<User>("/me");

export const logoutUser = () => http.post<void>("/auth/logout");
