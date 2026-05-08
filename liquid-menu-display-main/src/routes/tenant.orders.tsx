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

const STATUSES = ["diproses", "dimasak", "siap_diambil", "selesai"];

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
            const total = items.reduce((acc: number, curr: any) => acc + (curr.price_at_order * curr.quantity), 0);
            const itemCount = items.reduce((acc: number, curr: any) => acc + curr.quantity, 0);

            return (
            <div key={o.id} className="glass flex flex-col gap-3 rounded-2xl p-5 sm:flex-row sm:items-center sm:justify-between">
              <div>
                <div className="font-semibold">Sub-Order #{o.id?.slice(0, 8)}</div>
                <div className="text-xs text-muted-foreground">
                   {itemCount} items · {formatRupiah(total)}
                </div>
              </div>
              <Select defaultValue={o.tenant_status || o.status} onValueChange={(v) => updateStatus(o.id, v)}>
                <SelectTrigger className="w-48"><SelectValue /></SelectTrigger>
                <SelectContent>
                  {STATUSES.map((s) => <SelectItem key={s} value={s}>{s.replace("_", " ")}</SelectItem>)}
                </SelectContent>
              </Select>
            </div>
            );
          })}
        </div>
      }
    </DashboardShell>
  );
}
