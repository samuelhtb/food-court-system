import { useEffect, useState } from "react";
import { createFileRoute } from "@tanstack/react-router";
import { DashboardShell } from "@/components/DashboardShell";
import { ClipboardList, Users, Plus, Loader2 } from "lucide-react";
import { api } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger, DialogFooter } from "@/components/ui/dialog";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";

const NAV = [
  { to: "/admin/orders", label: "All Orders", icon: <ClipboardList className="h-4 w-4" /> },
  { to: "/admin/tenants", label: "Tenants", icon: <Users className="h-4 w-4" /> },
];

export const Route = createFileRoute("/admin/tenants")({
  component: AdminTenants,
});

function AdminTenants() {
  const [tenants, setTenants] = useState<any[] | null>(null);
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [form, setForm] = useState({
    username: "",
    password: "",
    email: "",
    name: "",
    tenant_name: "",
  });

  const load = () => api<any>("/admin/tenants", { auth: true }).then((d) => setTenants(d.data || d)).catch((e) => toast.error(e.message));
  useEffect(() => { load(); }, []);

  const save = async () => {
    if (!form.username || !form.password || !form.email || !form.name || !form.tenant_name) {
      toast.error("All fields are required");
      return;
    }

    setLoading(true);
    try {
      await api("/admin/tenants", { method: "POST", body: JSON.stringify(form), auth: true });
      toast.success("Tenant added successfully");
      setOpen(false);
      setForm({ username: "", password: "", email: "", name: "", tenant_name: "" });
      load();
    } catch (e: any) {
      toast.error(e.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <DashboardShell title="Admin / Cashier" role="admin" nav={NAV}>
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="font-display text-3xl font-bold">Tenants</h1>
          <p className="text-sm text-muted-foreground">Manage your food court stands.</p>
        </div>
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogTrigger asChild>
            <Button className="rounded-xl"><Plus className="h-4 w-4 mr-2" /> Add Tenant</Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle>Add New Tenant</DialogTitle>
            </DialogHeader>
            <div className="space-y-4 py-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label>Username</Label>
                  <Input value={form.username} onChange={(e) => setForm({ ...form, username: e.target.value })} placeholder="johndoe" />
                </div>
                <div className="space-y-2">
                  <Label>Password</Label>
                  <Input type="password" value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })} placeholder="••••••••" />
                </div>
              </div>
              <div className="space-y-2">
                <Label>Email</Label>
                <Input type="email" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} placeholder="tenant@example.com" />
              </div>
              <div className="space-y-2">
                <Label>Full Name (Owner)</Label>
                <Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} placeholder="John Doe" />
              </div>
              <div className="space-y-2">
                <Label>Stand / Tenant Name</Label>
                <Input value={form.tenant_name} onChange={(e) => setForm({ ...form, tenant_name: e.target.value })} placeholder="Sate Padang Mak Syukur" />
              </div>
            </div>
            <DialogFooter>
              <Button onClick={save} disabled={loading} className="w-full bg-burgundy hover:bg-burgundy/90 rounded-xl">
                {loading ? <><Loader2 className="h-4 w-4 animate-spin mr-2" /> Processing...</> : "Create Tenant"}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      {!tenants ? <Skeleton className="h-96 rounded-2xl" /> :
        tenants.length === 0 ? <div className="glass rounded-3xl p-12 text-center text-muted-foreground">No tenants yet.</div> :
        <div className="glass-strong overflow-hidden rounded-3xl">
          <table className="w-full text-sm">
            <thead className="bg-burgundy text-burgundy-foreground">
              <tr className="text-left">
                <th className="p-4">Stand Name</th>
                <th className="p-4">Owner</th>
                <th className="p-4">Email</th>
                <th className="p-4">Username</th>
              </tr>
            </thead>
            <tbody>
              {tenants.map((t) => (
                <tr key={t.id} className="border-t transition-colors hover:bg-muted/50">
                  <td className="p-4 font-bold">{t.tenant_name}</td>
                  <td className="p-4">{t.name}</td>
                  <td className="p-4">{t.email}</td>
                  <td className="p-4 text-muted-foreground">@{t.username}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      }
    </DashboardShell>
  );
}
