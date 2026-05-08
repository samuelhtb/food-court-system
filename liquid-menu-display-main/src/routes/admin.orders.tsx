import { useEffect, useState } from "react";
import { createFileRoute } from "@tanstack/react-router";
import { DashboardShell } from "@/components/DashboardShell";
import { ClipboardList } from "lucide-react";
import { api } from "@/lib/api";
import { formatRupiah } from "@/stores/cart";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";

const NAV = [
  { to: "/admin/orders", label: "All Orders", icon: <ClipboardList className="h-4 w-4" /> },
];

export const Route = createFileRoute("/admin/orders")({
  component: AdminOrders,
});

function AdminOrders() {
  const [orders, setOrders] = useState<any[] | null>(null);
  const load = () => api<any>("/admin/orders").then((d) => setOrders(d.data || d.orders || d)).catch((e) => toast.error(e.message));
  useEffect(() => { load(); }, []);

  const markPaid = async (id: string) => {
    try { await api(`/admin/orders/${id}/pay`, { method: "PUT" }); toast.success("Marked as paid"); load(); }
    catch (e: any) { toast.error(e.message); }
  };

  return (
    <DashboardShell title="Admin / Cashier" role="admin" nav={NAV}>
      <h1 className="mb-8 font-display text-3xl font-bold">All Orders</h1>
      {!orders ? <Skeleton className="h-96 rounded-2xl" /> :
        orders.length === 0 ? <div className="glass rounded-3xl p-12 text-center text-muted-foreground">No orders.</div> :
        <div className="glass-strong overflow-hidden rounded-2xl">
          <table className="w-full text-sm">
            <thead className="bg-burgundy text-burgundy-foreground">
              <tr className="text-left">
                <th className="p-4">Customer</th>
                <th className="p-4">Total</th>
                <th className="p-4">Method</th>
                <th className="p-4">Payment</th>
                <th className="p-4 text-right">Action</th>
              </tr>
            </thead>
            <tbody>
              {orders.map((o) => {
                const paid = (o.payment_status || "").toLowerCase() === "paid" || o.is_paid;
                return (
                  <tr key={o.id} className="border-t">
                    <td className="p-4 font-medium">{o.customer_name}</td>
                    <td className="p-4">{formatRupiah(o.total_price || o.total || 0)}</td>
                    <td className="p-4 capitalize">{o.payment_method}</td>
                    <td className="p-4">
                      <span className={`rounded-full px-2 py-0.5 text-xs font-semibold ${paid ? "bg-green-100 text-green-800" : "bg-amber-100 text-amber-800"}`}>
                        {paid ? "Paid" : "Pending"}
                      </span>
                    </td>
                    <td className="p-4 text-right">
                      {!paid && o.payment_method === "cash" && (
                        <Button size="sm" onClick={() => markPaid(o.id)}>Mark Paid</Button>
                      )}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      }
    </DashboardShell>
  );
}
