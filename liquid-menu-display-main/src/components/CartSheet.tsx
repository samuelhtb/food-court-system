import { useState } from "react";
import { useCart, formatRupiah } from "@/stores/cart";
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetFooter } from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Minus, Plus, Trash2, ShoppingBag, Copy, Check } from "lucide-react";
import { api } from "@/lib/api";
import { toast } from "sonner";

export function CartSheet({ open, onOpenChange }: { open: boolean; onOpenChange: (v: boolean) => void }) {
  const { items, setQty, remove, total, clear } = useCart();
  const [checkout, setCheckout] = useState(false);
  const [name, setName] = useState("");
  const [method, setMethod] = useState<"cash" | "qris">("cash");
  const [loading, setLoading] = useState(false);
  const [orderId, setOrderId] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  const submit = async () => {
    if (!name.trim()) return toast.error("Please enter your name");
    setLoading(true);
    try {
      const res: any = await api("/orders", {
        method: "POST",
        body: JSON.stringify({
          customer_name: name,
          payment_method: method,
          items: items.map((i) => ({ menu_id: i.id, quantity: i.quantity })),
        }),
      });
      const id = res?.id || res?.data?.id || res?.order?.id || "";
      setOrderId(id);
      clear();
      toast.success("Order placed!");
    } catch (e: any) {
      toast.error(e.message);
    } finally {
      setLoading(false);
    }
  };

  const reset = () => {
    setCheckout(false);
    setOrderId(null);
    setName("");
    setMethod("cash");
  };

  return (
    <Sheet open={open} onOpenChange={(v) => { onOpenChange(v); if (!v) reset(); }}>
      <SheetContent className="flex w-full flex-col gap-0 sm:max-w-md">
        <SheetHeader>
          <SheetTitle className="font-display text-2xl">
            {orderId ? "Order Confirmed" : checkout ? "Checkout" : "Your Cart"}
          </SheetTitle>
        </SheetHeader>

        {orderId ? (
          <div className="flex flex-1 flex-col items-center justify-center gap-4 p-6 text-center">
            <div className="rounded-full bg-primary/10 p-4">
              <Check className="h-10 w-10 text-burgundy" />
            </div>
            <p className="text-muted-foreground">Save your Order ID to track status.</p>
            <div className="glass flex w-full items-center justify-between gap-2 rounded-xl p-4">
              <code className="truncate text-sm font-medium">{orderId}</code>
              <Button
                size="sm"
                variant="ghost"
                onClick={() => { navigator.clipboard.writeText(orderId); setCopied(true); toast.success("Copied"); }}
              >
                {copied ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
              </Button>
            </div>
            <Button onClick={() => { onOpenChange(false); reset(); }} className="w-full">Done</Button>
          </div>
        ) : checkout ? (
          <div className="flex flex-1 flex-col gap-4 overflow-y-auto p-6">
            <div className="space-y-2">
              <Label>Name</Label>
              <Input value={name} onChange={(e) => setName(e.target.value)} placeholder="Your name" />
            </div>
            <div className="space-y-2">
              <Label>Payment Method</Label>
              <Select value={method} onValueChange={(v: any) => setMethod(v)}>
                <SelectTrigger><SelectValue /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="cash">Cash</SelectItem>
                  <SelectItem value="qris">QRIS</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="glass mt-2 rounded-xl p-4">
              <div className="mb-2 text-sm font-semibold">Order Summary</div>
              {items.map((i) => (
                <div key={i.id} className="flex justify-between py-1 text-sm">
                  <span>{i.name} × {i.quantity}</span>
                  <span>{formatRupiah(i.price * i.quantity)}</span>
                </div>
              ))}
              <div className="mt-3 flex justify-between border-t pt-3 font-bold">
                <span>Total</span><span>{formatRupiah(total())}</span>
              </div>
            </div>
            <SheetFooter className="mt-auto flex-col gap-2 sm:flex-col">
              <Button onClick={submit} disabled={loading} className="w-full">
                {loading ? "Placing..." : "Place Order"}
              </Button>
              <Button onClick={() => setCheckout(false)} variant="outline" className="w-full">Back</Button>
            </SheetFooter>
          </div>
        ) : items.length === 0 ? (
          <div className="flex flex-1 flex-col items-center justify-center gap-3 p-6 text-center text-muted-foreground">
            <ShoppingBag className="h-12 w-12 opacity-30" />
            <p>Your cart is empty</p>
          </div>
        ) : (
          <>
            <div className="flex-1 space-y-3 overflow-y-auto p-6">
              {items.map((i) => (
                <div key={i.id} className="glass flex gap-3 rounded-xl p-3">
                  {i.image_url && (
                    <img src={i.image_url} alt={i.name} className="h-16 w-16 rounded-lg object-cover" />
                  )}
                  <div className="flex flex-1 flex-col justify-between">
                    <div className="flex items-start justify-between">
                      <div>
                        <div className="text-sm font-semibold">{i.name}</div>
                        <div className="text-xs text-burgundy">{formatRupiah(i.price)}</div>
                      </div>
                      <button onClick={() => remove(i.id)} className="text-muted-foreground hover:text-destructive">
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </div>
                    <div className="flex items-center gap-2">
                      <button onClick={() => setQty(i.id, i.quantity - 1)} className="rounded-full border p-1 hover:bg-muted">
                        <Minus className="h-3 w-3" />
                      </button>
                      <span className="w-6 text-center text-sm font-medium">{i.quantity}</span>
                      <button onClick={() => setQty(i.id, i.quantity + 1)} className="rounded-full border p-1 hover:bg-muted">
                        <Plus className="h-3 w-3" />
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
            <SheetFooter className="border-t bg-background/50 p-6 backdrop-blur">
              <div className="flex w-full flex-col gap-3">
                <div className="flex justify-between text-lg font-bold">
                  <span>Total</span><span>{formatRupiah(total())}</span>
                </div>
                <Button onClick={() => setCheckout(true)} className="w-full" size="lg">Checkout</Button>
              </div>
            </SheetFooter>
          </>
        )}
      </SheetContent>
    </Sheet>
  );
}
