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
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  const [menus, setMenus] = useState<Menu[] | null>(null);
  const [cartOpen, setCartOpen] = useState(false);
  const [selectedMenu, setSelectedMenu] = useState<Menu | null>(null);
  const add = useCart((s) => s.add);

  useEffect(() => {
    api<any>("/menus", { auth: false })
      .then((d) => setMenus(d.data || d.menus || (Array.isArray(d) ? d : [])))
      .catch((e) => { 
        toast.error(e.message); 
        setMenus([]); 
      });
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
              <article 
                key={m.id} 
                className="glass group flex flex-col overflow-hidden rounded-3xl transition hover:-translate-y-1 hover:shadow-2xl cursor-pointer"
                onClick={() => setSelectedMenu(m)}
              >
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
                      onClick={(e) => { 
                        e.stopPropagation();
                        add(m); 
                        toast.success(`${m.name} added`); 
                      }}
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

      <Dialog open={!!selectedMenu} onOpenChange={(open) => !open && setSelectedMenu(null)}>
        <DialogContent className="sm:max-w-[425px] overflow-hidden rounded-3xl p-0 gap-0 border-none bg-background">
          {selectedMenu && (
            <>
              {selectedMenu.image_url ? (
                <div className="relative h-64 w-full bg-muted">
                  <img src={selectedMenu.image_url} alt={selectedMenu.name} className="absolute inset-0 h-full w-full object-cover" />
                  <div className="absolute top-4 left-4 glass-strong rounded-full px-3 py-1 text-xs font-semibold">
                    Stock: {selectedMenu.stock}
                  </div>
                </div>
              ) : (
                <div className="relative flex h-64 w-full items-center justify-center bg-muted text-muted-foreground">
                  No image
                  <div className="absolute top-4 left-4 glass-strong rounded-full px-3 py-1 text-xs font-semibold">
                    Stock: {selectedMenu.stock}
                  </div>
                </div>
              )}
              
              <div className="p-6">
                <DialogHeader className="text-left space-y-1 mb-4">
                  <div className="flex items-center justify-between pr-8">
                    <DialogTitle className="font-display text-2xl font-bold">{selectedMenu.name}</DialogTitle>
                  </div>
                  <DialogDescription className="mt-1">
                    {selectedMenu.tenant_name && (
                      <span className="text-[10px] font-bold uppercase tracking-tighter text-muted-foreground/60 bg-muted px-2 py-0.5 rounded-full inline-block mb-2">
                        {selectedMenu.tenant_name}
                      </span>
                    )}
                  </DialogDescription>
                </DialogHeader>
                
                <div className="max-h-[120px] overflow-y-auto mb-6 pr-2">
                  <p className="text-sm text-muted-foreground leading-relaxed">{selectedMenu.description}</p>
                </div>
                
                <div className="flex justify-between items-center pt-4 border-t border-border/50">
                  <div>
                    <span className="text-sm text-muted-foreground block mb-1">Price</span>
                    <span className="text-2xl font-bold text-burgundy">{formatRupiah(selectedMenu.price)}</span>
                  </div>
                  <Button
                    onClick={() => {
                      add(selectedMenu);
                      toast.success(`${selectedMenu.name} added`);
                      setSelectedMenu(null);
                    }}
                    disabled={selectedMenu.stock <= 0}
                    className="rounded-full px-6"
                    size="lg"
                  >
                    <Plus className="h-5 w-5 mr-2" /> {selectedMenu.stock > 0 ? "Add to Cart" : "Out of Stock"}
                  </Button>
                </div>
              </div>
            </>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}
