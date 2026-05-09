import { useEffect, useState } from "react";
import { createFileRoute } from "@tanstack/react-router";
import { DashboardShell } from "@/components/DashboardShell";
import { ClipboardList, Eye } from "lucide-react";
import { api } from "@/lib/api";
import { formatRupiah } from "@/stores/cart";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { toast } from "sonner";

const NAV = [
  { to: "/admin/orders", label: "All Orders", icon: <ClipboardList className="h-4 w-4" /> },
];

export const Route = createFileRoute("/admin/orders")({
  component: AdminOrders,
});

function AdminOrders() {
  const [orders, setOrders] = useState<any[] | null>(null);
  const [selectedOrder, setSelectedOrder] = useState<any | null>(null);
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
                    <td className="p-4">{formatRupiah(o.total_amount || 0)}</td>
                    <td className="p-4 capitalize">{o.payment_method}</td>
                    <td className="p-4">
                      <span className={`rounded-full px-2 py-0.5 text-xs font-semibold ${paid ? "bg-green-100 text-green-800" : "bg-amber-100 text-amber-800"}`}>
                        {paid ? "Paid" : "Pending"}
                      </span>
                    </td>
                    <td className="p-4 text-right flex justify-end gap-2">
                      <Button size="sm" variant="outline" onClick={() => setSelectedOrder(o)}>
                        <Eye className="h-4 w-4 mr-1" /> Details
                      </Button>
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

      <Dialog open={!!selectedOrder} onOpenChange={(open) => !open && setSelectedOrder(null)}>
        <DialogContent className="max-w-md data-[state=open]:zoom-in-100 data-[state=open]:slide-in-from-top-[50%] data-[state=open]:slide-in-from-left-[50%]">
          <DialogHeader>
            <DialogTitle>Order Details</DialogTitle>
          </DialogHeader>
          {selectedOrder && (
            <div className="space-y-4 text-sm">
              <div className="flex justify-between border-b pb-2">
                <span className="font-medium">Order ID:</span>
                <code className="text-xs bg-muted px-1 rounded">{selectedOrder.id}</code>
              </div>
              <div className="flex justify-between border-b pb-2">
                <span className="font-medium">Customer:</span>
                <span>{selectedOrder.customer_name}</span>
              </div>
              <div>
                <h4 className="font-semibold mb-2">Items:</h4>
                <ul className="space-y-2">
                  {(selectedOrder.sub_orders || []).map((so: any) => 
                    (so.order_items || []).map((item: any) => (
                      <li key={item.id} className="flex justify-between text-muted-foreground">
                        <span>{item.quantity}x {item.menu?.name || "Unknown Menu"}</span>
                        <span>{formatRupiah(item.price_at_order * item.quantity)}</span>
                      </li>
                    ))
                  )}
                </ul>
              </div>
              <div className="flex justify-between border-t pt-2 font-bold text-base">
                <span>Total</span>
                <span className="text-burgundy">{formatRupiah(selectedOrder.total_amount || 0)}</span>
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </DashboardShell>
  );
}
