import { create } from "zustand";
import { persist } from "zustand/middleware";

export type Role = "admin" | "tenant" | null;

interface AuthState {
  token: string | null;
  role: Role;
  email: string | null;
  setAuth: (token: string, role: Role, email: string) => void;
  logout: () => void;
}

export const useAuth = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      role: null,
      email: null,
      setAuth: (token, role, email) => {
        if (typeof window !== "undefined") localStorage.setItem("fc_token", token);
        set({ token, role, email });
      },
      logout: () => {
        if (typeof window !== "undefined") localStorage.removeItem("fc_token");
        set({ token: null, role: null, email: null });
      },
    }),
    { name: "fc_auth" },
  ),
);
