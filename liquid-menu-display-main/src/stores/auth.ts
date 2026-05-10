import { create } from "zustand";
import { persist } from "zustand/middleware";

export type Role = "admin" | "tenant" | null;

interface AuthState {
  token: string | null;
  role: Role;
  email: string | null;
  _hasHydrated: boolean;
  setAuth: (token: string, role: Role, email: string) => void;
  logout: () => void;
  setHasHydrated: (state: boolean) => void;
}

export const useAuth = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      role: null,
      email: null,
      _hasHydrated: false,
      setAuth: (token, role, email) => {
        if (typeof window !== "undefined") localStorage.setItem("fc_token", token);
        set({ token, role, email });
      },
      logout: () => {
        if (typeof window !== "undefined") localStorage.removeItem("fc_token");
        set({ token: null, role: null, email: null });
      },
      setHasHydrated: (state) => set({ _hasHydrated: state }),
    }),
    {
      name: "fc_auth",
      onRehydrateStorage: (state) => {
        return () => state.setHasHydrated(true);
      },
    },
  ),
);
