import { create } from "zustand";
import { persist } from "zustand/middleware";

export interface Menu {
  id: string;
  name: string;
  description?: string;
  price: number;
  stock: number;
  image_url?: string;
  is_available?: boolean;
  tenant_name?: string;
}

export interface CartItem extends Menu {
  quantity: number;
}

interface CartState {
  items: CartItem[];
  add: (menu: Menu) => void;
  remove: (id: string) => void;
  setQty: (id: string, qty: number) => void;
  clear: () => void;
  total: () => number;
  count: () => number;
}

export const useCart = create<CartState>()(
  persist(
    (set, get) => ({
      items: [],
      add: (menu) =>
        set((s) => {
          const ex = s.items.find((i) => i.id === menu.id);
          if (ex) {
            return {
              items: s.items.map((i) =>
                i.id === menu.id ? { ...i, quantity: Math.min(i.quantity + 1, menu.stock) } : i,
              ),
            };
          }
          return { items: [...s.items, { ...menu, quantity: 1 }] };
        }),
      remove: (id) => set((s) => ({ items: s.items.filter((i) => i.id !== id) })),
      setQty: (id, qty) =>
        set((s) => ({
          items: s.items.map((i) =>
            i.id === id ? { ...i, quantity: Math.max(1, Math.min(qty, i.stock)) } : i,
          ),
        })),
      clear: () => set({ items: [] }),
      total: () => get().items.reduce((s, i) => s + i.price * i.quantity, 0),
      count: () => get().items.reduce((s, i) => s + i.quantity, 0),
    }),
    { name: "fc_cart" },
  ),
);

export const formatRupiah = (n: number) =>
  new Intl.NumberFormat("id-ID", { style: "currency", currency: "IDR", maximumFractionDigits: 0 }).format(n);
