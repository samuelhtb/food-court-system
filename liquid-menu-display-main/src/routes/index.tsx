import { useEffect, useState } from "react";
import { createFileRoute } from "@tanstack/react-router";
import { api } from "@/lib/api";
import { useCart, formatRupiah, type Menu } from "@/stores/cart";
import { Button } from "@/components/ui/button";
import { PublicNav } from "@/components/PublicNav";
import { CartSheet } from "@/components/CartSheet";
import { Skeleton } from "@/components/ui/skeleton";
import { Plus, Sparkles } from "lucide-react";
import { toast } from "sonner";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  const [menus, setMenus] = useState<Menu[] | null>(null);
  const [cartOpen, setCartOpen] = useState(false);
  const add = useCart((s) => s.add);

  useEffect(() => {
    api<any>("/menus", { auth: false })
      .then((d) => setMenus(d.data || d.menus || d))
      .catch((e) => { toast.error(e.message); setMenus([]); });
  }, []);

  return (
    <div className="min-h-screen">
      <PublicNav onCart={() => setCartOpen(true)} />

      {/* Hero */}
      <section className="mx-auto max-w-6xl px-6 pt-16 pb-12 text-center">
        <div className="glass mx-auto mb-6 inline-flex items-center gap-2 rounded-full px-4 py-2 text-xs font-medium text-burgundy">
          <Sparkles className="h-3 w-3" /> Curated tenants · Lightning fast
        </div>
        <h1 className="font-display text-5xl font-bold leading-[1.05] sm:text-7xl">
          Taste the food court,<br />
          <span className="text-burgundy italic">redefined.</span>
        </h1>
        <p className="mx-auto mt-5 max-w-xl text-base text-muted-foreground">
          Order from the best stalls in one place. Fresh menus, real-time tracking,
          zero queueing.
        </p>
      </section>

      {/* Menu grid */}
      <section className="mx-auto max-w-6xl px-6 pb-24">
        <div className="mb-6 flex items-end justify-between">
          <h2 className="font-display text-3xl font-bold">Today&apos;s Menu</h2>
          <span className="text-sm text-muted-foreground">{menus?.length ?? 0} items</span>
        </div>

        {!menus ? (
          <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
            {Array.from({ length: 6 }).map((_, i) => (
              <Skeleton key={i} className="h-80 rounded-3xl" />
            ))}
          </div>
        ) : menus.length === 0 ? (
          <div className="glass rounded-3xl p-12 text-center text-muted-foreground">
            No menus available right now. Please check back later.
          </div>
        ) : (
          <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
            {menus.map((m) => (
              <article key={m.id} className="glass group flex flex-col overflow-hidden rounded-3xl transition hover:-translate-y-1 hover:shadow-2xl">
                <div className="relative aspect-[4/3] overflow-hidden bg-muted">
                  {m.image_url ? (
                    <img src={m.image_url} alt={m.name} loading="lazy" className="h-full w-full object-cover transition duration-500 group-hover:scale-105" />
                  ) : (
                    <div className="flex h-full items-center justify-center text-muted-foreground">No image</div>
                  )}
                  <div className="absolute right-3 top-3 glass-strong rounded-full px-3 py-1 text-xs font-semibold">
                    Stock: {m.stock}
                  </div>
                </div>
                <div className="flex flex-1 flex-col p-5">
                  <div className="flex items-center justify-between">
                    <h3 className="font-display text-xl font-bold">{m.name}</h3>
                    {m.tenant_name && (
                      <span className="text-[10px] font-bold uppercase tracking-tighter text-muted-foreground/60 bg-muted px-2 py-0.5 rounded-full">
                        {m.tenant_name}
                      </span>
                    )}
                  </div>
                  <p className="mt-1 line-clamp-2 text-sm text-muted-foreground">{m.description}</p>
                  <div className="mt-auto flex items-center justify-between pt-4">
                    <span className="text-lg font-bold text-burgundy">{formatRupiah(m.price)}</span>
                    <Button
                      onClick={() => { add(m); toast.success(`${m.name} added`); }}
                      disabled={m.stock <= 0}
                      size="sm"
                      className="rounded-full"
                    >
                      <Plus className="h-4 w-4" /> Add
                    </Button>
                  </div>
                </div>
              </article>
            ))}
          </div>
        )}
      </section>

      <footer className="border-t py-8 text-center text-xs text-muted-foreground">
        © {new Date().getFullYear()} Saji Food Court
      </footer>

      <CartSheet open={cartOpen} onOpenChange={setCartOpen} />
    </div>
  );
}
