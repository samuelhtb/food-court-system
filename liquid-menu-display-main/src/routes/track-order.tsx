import { useState } from "react";
import { createFileRoute } from "@tanstack/react-router";
import { api } from "@/lib/api";
import { PublicNav } from "@/components/PublicNav";
import { CartSheet } from "@/components/CartSheet";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { formatRupiah } from "@/stores/cart";
import { toast } from "sonner";
import { Search } from "lucide-react";

export const Route = createFileRoute("/track-order")({
  component: TrackOrder,
});

function TrackOrder() {
  const [id, setId] = useState("");
  const [order, setOrder] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const [cartOpen, setCartOpen] = useState(false);

  const search = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!id.trim()) return;
    setLoading(true);
    try {
      const res: any = await api(`/orders/${id}`, { auth: false });
      setOrder(res.data || res);
    } catch (err: any) {
      toast.error(err.message);
      setOrder(null);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen">
      <PublicNav onCart={() => setCartOpen(true)} />
      <section className="mx-auto max-w-2xl px-6 py-16">
        <h1 className="font-display text-4xl font-bold">Track your order</h1>
        <p className="mt-2 text-muted-foreground">Enter your Order ID to see live status.</p>
        <form onSubmit={search} className="glass mt-8 flex flex-col gap-3 rounded-2xl p-5 sm:flex-row">
          <Input value={id} onChange={(e) => setId(e.target.value)} placeholder="Order ID" className="flex-1 bg-white/60" />
          <Button type="submit" disabled={loading}>
            <Search className="h-4 w-4" /> {loading ? "Searching..." : "Track"}
          </Button>
        </form>

        {order && (
          <div className="glass-strong mt-8 rounded-2xl p-6">
            <div className="flex items-center justify-between">
              <div>
                <div className="text-xs text-muted-foreground">Customer</div>
                <div className="font-semibold">{order.customer_name}</div>
              </div>
              <span className="rounded-full bg-burgundy px-3 py-1 text-xs font-bold uppercase text-burgundy-foreground">
                {order.status || order.payment_status || "pending"}
              </span>
            </div>
            <div className="mt-4 space-y-2">
              {(order.items || []).map((it: any, i: number) => (
                <div key={i} className="flex justify-between text-sm">
                  <span>{it.menu?.name || it.name} × {it.quantity}</span>
                  <span>{formatRupiah((it.price || it.menu?.price || 0) * it.quantity)}</span>
                </div>
              ))}
            </div>
            <div className="mt-4 flex justify-between border-t pt-4 font-bold">
              <span>Total</span><span>{formatRupiah(order.total_amount || order.total_price || order.total || 0)}</span>
            </div>
            <div className="mt-2 text-xs text-muted-foreground">
              Payment: {order.payment_method} · {order.payment_status}
            </div>
          </div>
        )}
      </section>
      <CartSheet open={cartOpen} onOpenChange={setCartOpen} />
    </div>
  );
}
