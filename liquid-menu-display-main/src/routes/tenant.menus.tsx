import { useEffect, useState } from "react";
import { createFileRoute } from "@tanstack/react-router";
import { DashboardShell } from "@/components/DashboardShell";
import { UtensilsCrossed, ClipboardList, Wallet, Pencil, Trash2, Plus } from "lucide-react";
import { api } from "@/lib/api";
import { formatRupiah, type Menu } from "@/stores/cart";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger, DialogFooter } from "@/components/ui/dialog";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";

const NAV = [
  { to: "/tenant/menus", label: "Menus", icon: <UtensilsCrossed className="h-4 w-4" /> },
  { to: "/tenant/orders", label: "Orders", icon: <ClipboardList className="h-4 w-4" /> },
  { to: "/tenant/earnings", label: "Earnings", icon: <Wallet className="h-4 w-4" /> },
];

export const Route = createFileRoute("/tenant/menus")({
  component: TenantMenus,
});

interface FormData { name: string; description: string; price: number; stock: number; image_url: string; is_available: boolean }
const empty: FormData = { name: "", description: "", price: 0, stock: 0, image_url: "", is_available: true };

function TenantMenus() {
  const [menus, setMenus] = useState<Menu[] | null>(null);
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<Menu | null>(null);
  const [form, setForm] = useState<FormData>(empty);

  const load = () => api<any>("/tenant/menus").then((d) => setMenus(d.data || d.menus || d)).catch((e) => toast.error(e.message));
  useEffect(() => { load(); }, []);

  const openCreate = () => { setEditing(null); setForm(empty); setOpen(true); };
  const openEdit = (m: Menu) => {
    setEditing(m);
    setForm({ name: m.name, description: m.description || "", price: m.price, stock: m.stock, image_url: m.image_url || "", is_available: m.is_available ?? true });
    setOpen(true);
  };

  const save = async () => {
    try {
      if (editing) {
        await api(`/menus/${editing.id}`, { method: "PUT", body: JSON.stringify(form) });
        toast.success("Menu updated");
      } else {
        const { is_available, ...payload } = form;
        await api(`/menus`, { method: "POST", body: JSON.stringify(payload) });
        toast.success("Menu created");
      }
      setOpen(false);
      load();
    } catch (e: any) { toast.error(e.message); }
  };

  const del = async (id: string) => {
    if (!confirm("Delete this menu?")) return;
    try { await api(`/menus/${id}`, { method: "DELETE" }); toast.success("Deleted"); load(); }
    catch (e: any) { toast.error(e.message); }
  };

  return (
    <DashboardShell title="Tenant" role="tenant" nav={NAV}>
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="font-display text-3xl font-bold">My Menus</h1>
          <p className="text-sm text-muted-foreground">Manage your stand&apos;s offerings.</p>
        </div>
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogTrigger asChild>
            <Button onClick={openCreate}><Plus className="h-4 w-4" /> New Menu</Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader><DialogTitle>{editing ? "Edit Menu" : "New Menu"}</DialogTitle></DialogHeader>
            <div className="space-y-3">
              <div><Label>Name</Label><Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} /></div>
              <div><Label>Description</Label><Textarea value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} /></div>
              <div className="grid grid-cols-2 gap-3">
                <div><Label>Price</Label><Input type="number" min={0} value={form.price} onChange={(e) => setForm({ ...form, price: +e.target.value })} /></div>
                <div><Label>Stock</Label><Input type="number" min={0} value={form.stock} onChange={(e) => setForm({ ...form, stock: +e.target.value })} /></div>
              </div>
              <div><Label>Image URL</Label><Input value={form.image_url} onChange={(e) => setForm({ ...form, image_url: e.target.value })} /></div>
              {editing && (
                <label className="flex items-center gap-2 text-sm">
                  <input type="checkbox" checked={form.is_available} onChange={(e) => setForm({ ...form, is_available: e.target.checked })} />
                  Available
                </label>
              )}
            </div>
            <DialogFooter><Button onClick={save}>Save</Button></DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      {!menus ? (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 3 }).map((_, i) => <Skeleton key={i} className="h-64 rounded-2xl" />)}
        </div>
      ) : menus.length === 0 ? (
        <div className="glass rounded-3xl p-12 text-center text-muted-foreground">No menus yet. Create your first one.</div>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {menus.map((m) => (
            <div key={m.id} className="glass overflow-hidden rounded-2xl">
              {m.image_url && <img src={m.image_url} alt={m.name} className="aspect-video w-full object-cover" />}
              <div className="p-4">
                <div className="flex items-start justify-between">
                  <h3 className="font-display text-lg font-bold">{m.name}</h3>
                  <span className="text-sm font-bold text-burgundy">{formatRupiah(m.price)}</span>
                </div>
                <p className="mt-1 line-clamp-2 text-xs text-muted-foreground">{m.description}</p>
                <div className="mt-3 flex items-center justify-between text-xs">
                  <span className="text-muted-foreground">Stock: {m.stock}</span>
                  <div className="flex gap-1">
                    <Button size="sm" variant="ghost" onClick={() => openEdit(m)}><Pencil className="h-3 w-3" /></Button>
                    <Button size="sm" variant="ghost" onClick={() => del(m.id)}><Trash2 className="h-3 w-3" /></Button>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </DashboardShell>
  );
}
