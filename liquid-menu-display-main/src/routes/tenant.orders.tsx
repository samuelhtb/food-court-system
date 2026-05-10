import { useEffect, useState } from "react";
import { createFileRoute } from "@tanstack/react-router";
import { DashboardShell } from "@/components/DashboardShell";
import { UtensilsCrossed, ClipboardList, Wallet } from "lucide-react";
import { api } from "@/lib/api";
import { formatRupiah } from "@/stores/cart";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";

const NAV = [
  { to: "/tenant/menus", label: "Menus", icon: <UtensilsCrossed className="h-4 w-4" /> },
  { to: "/tenant/orders", label: "Orders", icon: <ClipboardList className="h-4 w-4" /> },
  { to: "/tenant/earnings", label: "Earnings", icon: <Wallet className="h-4 w-4" /> },
];

const STATUSES = ["menunggu_pembayaran", "diproses", "selesai"];

export const Route = createFileRoute("/tenant/orders")({
  component: TenantOrders,
});

function TenantOrders() {
  const [orders, setOrders] = useState<any[] | null>(null);

  const load = () => api<any>("/tenant/orders").then((d) => setOrders(d.data || d.orders || d)).catch((e) => toast.error(e.message));
  useEffect(() => { load(); }, []);

  const updateStatus = async (id: string, status: string) => {
    try { await api(`/tenant/orders/${id}/status`, { method: "PUT", body: JSON.stringify({ status }) }); toast.success("Status updated"); load(); }
    catch (e: any) { toast.error(e.message); }
  };

  return (
    <DashboardShell title="Tenant" role="tenant" nav={NAV}>
      <h1 className="mb-8 font-display text-3xl font-bold">Incoming Orders</h1>
      {!orders ? <Skeleton className="h-64 rounded-2xl" /> :
        orders.length === 0 ? <div className="glass rounded-3xl p-12 text-center text-muted-foreground">No orders yet.</div> :
        <div className="space-y-3">
          {orders.map((o) => {
            const items = o.order_items || [];
            const customerName = o.parent_order?.customer_name || "Customer";
            const total = items.reduce((acc: number, curr: any) => acc + (curr.price_at_order * curr.quantity), 0);

            return (
            <div key={o.id} className="glass flex flex-col gap-4 rounded-3xl p-6 sm:flex-row sm:items-start sm:justify-between">
              <div className="flex-1">
                <div className="mb-3 flex items-center gap-3">
                  <div className="font-display text-xl font-bold">Order #{o.id?.slice(0, 8)}</div>
                  <span className="rounded-full bg-burgundy/10 px-3 py-1 text-[10px] font-bold text-burgundy uppercase tracking-widest">
                    {customerName}
                  </span>
                </div>
                
                <div className="space-y-2 border-l-2 border-burgundy/20 pl-4">
                  {items.map((item: any) => (
                    <div key={item.id} className="flex items-center justify-between text-sm">
                      <span className="font-medium">
                        {item.menu?.name || "Unknown Menu"} 
                        <span className="ml-2 text-muted-foreground text-xs">×{item.quantity}</span>
                      </span>
                    </div>
                  ))}
                </div>

                <div className="mt-4 font-bold text-burgundy">{formatRupiah(total)}</div>
              </div>

              <div className="flex flex-col items-end gap-2">
                <span className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest">Order Status</span>
                <Select defaultValue={o.tenant_status || o.status} onValueChange={(v) => updateStatus(o.id, v)}>
                  <SelectTrigger className="w-48 rounded-xl border-none bg-muted/50 font-medium"><SelectValue /></SelectTrigger>
                  <SelectContent className="rounded-xl border-none">
                    {STATUSES.map((s) => (
                      <SelectItem key={s} value={s} className="rounded-lg">
                        {s.replace("_", " ").toUpperCase()}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
            );
          })}
        </div>
      }
    </DashboardShell>
  );
}
