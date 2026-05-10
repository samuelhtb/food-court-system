import { useEffect, useState } from "react";
import { createFileRoute } from "@tanstack/react-router";
import { DashboardShell } from "@/components/DashboardShell";
import { ClipboardList, Users, Wallet, ShieldCheck, Plus, Loader2 } from "lucide-react";
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
  { to: "/admin/earnings", label: "Earnings", icon: <Wallet className="h-4 w-4" /> },
  { to: "/admin/manage-admins", label: "Manage Admins", icon: <ShieldCheck className="h-4 w-4" /> },
];

export const Route = createFileRoute("/admin/manage-admins")({
  component: ManageAdmins,
});

function ManageAdmins() {
  const [admins, setAdmins] = useState<any[] | null>(null);
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [form, setForm] = useState({
    username: "",
    password: "",
    email: "",
    name: "",
  });

  const load = () => api<any>("/admin/manage", { auth: true }).then((d) => setAdmins(d.data || d)).catch((e) => toast.error(e.message));
  
  useEffect(() => { 
    load(); 
  }, []);

  const save = async () => {
    if (!form.username || !form.password || !form.email || !form.name) {
      toast.error("Semua field wajib diisi");
      return;
    }

    setLoading(true);
    try {
      await api("/admin/manage", { 
        method: "POST", 
        body: JSON.stringify(form), 
        auth: true 
      });
      toast.success("Admin berhasil ditambahkan");
      setOpen(false);
      setForm({ username: "", password: "", email: "", name: "" });
      load();
    } catch (e: any) {
      toast.error(e.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <DashboardShell title="Admin" role="admin" nav={NAV}>
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="font-display text-3xl font-bold">Kelola Admin</h1>
          <p className="text-sm text-muted-foreground">Tambah atau lihat daftar administrator sistem.</p>
        </div>
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogTrigger asChild>
            <Button className="rounded-xl"><Plus className="h-4 w-4 mr-2" /> Tambah Admin</Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle>Tambah Admin Baru</DialogTitle>
            </DialogHeader>
            <div className="space-y-4 py-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label>Username</Label>
                  <Input value={form.username} onChange={(e) => setForm({ ...form, username: e.target.value })} placeholder="admin_baru" />
                </div>
                <div className="space-y-2">
                  <Label>Password</Label>
                  <Input type="password" value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })} placeholder="••••••••" />
                </div>
              </div>
              <div className="space-y-2">
                <Label>Email</Label>
                <Input type="email" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} placeholder="admin@foodcourt.com" />
              </div>
              <div className="space-y-2">
                <Label>Nama Lengkap</Label>
                <Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} placeholder="Super Admin" />
              </div>
            </div>
            <DialogFooter>
              <Button onClick={save} disabled={loading} className="w-full bg-burgundy hover:bg-burgundy/90 rounded-xl">
                {loading ? <><Loader2 className="h-4 w-4 animate-spin mr-2" /> Memproses...</> : "Buat Akun Admin"}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      {!admins ? <Skeleton className="h-96 rounded-2xl" /> :
        admins.length === 0 ? <div className="glass rounded-3xl p-12 text-center text-muted-foreground">Belum ada admin lain.</div> :
        <div className="glass-strong overflow-hidden rounded-3xl">
          <table className="w-full text-sm">
            <thead className="bg-burgundy text-burgundy-foreground">
              <tr className="text-left">
                <th className="p-4">Nama</th>
                <th className="p-4">Email</th>
                <th className="p-4">Username</th>
                <th className="p-4">Role</th>
              </tr>
            </thead>
            <tbody>
              {admins.map((a) => (
                <tr key={a.id} className="border-t transition-colors hover:bg-muted/50">
                  <td className="p-4 font-bold">{a.name}</td>
                  <td className="p-4">{a.email}</td>
                  <td className="p-4 text-muted-foreground">@{a.username}</td>
                  <td className="p-4">
                    <span className="inline-flex items-center rounded-full bg-burgundy/10 px-2.5 py-0.5 text-xs font-medium text-burgundy">
                      {a.role}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      }
    </DashboardShell>
  );
}
