import { useEffect, useState } from "react";
import { createFileRoute } from "@tanstack/react-router";
import { DashboardShell } from "@/components/DashboardShell";
import { ClipboardList, Users, Wallet, TrendingUp, Receipt, Filter, Loader2, ShieldCheck } from "lucide-react";
import { api } from "@/lib/api";
import { formatRupiah } from "@/stores/cart";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";

const NAV = [
  { to: "/admin/orders", label: "All Orders", icon: <ClipboardList className="h-4 w-4" /> },
  { to: "/admin/tenants", label: "Tenants", icon: <Users className="h-4 w-4" /> },
  { to: "/admin/earnings", label: "Earnings", icon: <Wallet className="h-4 w-4" /> },
  { to: "/admin/manage-admins", label: "Manage Admins", icon: <ShieldCheck className="h-4 w-4" /> },
];

export const Route = createFileRoute("/admin/earnings")({
  component: AdminEarnings,
});

function AdminEarnings() {
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const [dates, setDates] = useState({ start: "", end: "" });

  const fetchReport = async () => {
    setLoading(true);
    try {
      let url = "/admin/reports/income";
      const params = new URLSearchParams();
      if (dates.start) params.append("start_date", dates.start);
      if (dates.end) params.append("end_date", dates.end);
      if (params.toString()) url += `?${params.toString()}`;

      const res = await api<any>(url, { auth: true });
      setData(res.data || res);
    } catch (e: any) {
      toast.error(e.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchReport();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <DashboardShell title="Admin" role="admin" nav={NAV}>
      <div className="mb-8 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="font-display text-3xl font-bold">Laporan Pemasukan Sistem</h1>
          <p className="text-sm text-muted-foreground">Monitor total pendapatan seluruh tenant.</p>
        </div>
        
        <div className="flex items-end gap-3 glass-strong p-3 rounded-2xl">
          <div className="space-y-1.5">
            <Label className="text-xs text-muted-foreground ml-1">Mulai Tanggal</Label>
            <input 
              type="date" 
              className="flex h-10 w-full rounded-xl border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-burgundy/50 disabled:cursor-not-allowed disabled:opacity-50"
              value={dates.start}
              onChange={(e) => setDates({ ...dates, start: e.target.value })}
            />
          </div>
          <div className="space-y-1.5">
            <Label className="text-xs text-muted-foreground ml-1">Sampai Tanggal</Label>
            <input 
              type="date" 
              className="flex h-10 w-full rounded-xl border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-burgundy/50 disabled:cursor-not-allowed disabled:opacity-50"
              value={dates.end}
              onChange={(e) => setDates({ ...dates, end: e.target.value })}
            />
          </div>
          <Button onClick={fetchReport} disabled={loading} className="rounded-xl h-10 bg-burgundy hover:bg-burgundy/90">
            {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Filter className="h-4 w-4" />}
          </Button>
        </div>
      </div>

      {!data && loading ? (
        <div className="grid gap-5 sm:grid-cols-2"><Skeleton className="h-40 rounded-3xl" /><Skeleton className="h-40 rounded-3xl" /></div>
      ) : (
        <div className="space-y-8">
          <div className="grid gap-5 sm:grid-cols-2">
            <Stat icon={<TrendingUp />} label="Total Pendapatan Sistem" value={formatRupiah(data?.total_revenue || 0)} />
            <Stat icon={<Receipt />} label="Total Transaksi Sistem" value={String(data?.total_orders || 0)} />
          </div>

          {data?.breakdown && data.breakdown.length > 0 && (
            <div className="glass-strong rounded-3xl p-8">
              <h2 className="mb-6 font-display text-2xl font-bold">Rincian Per Tenant</h2>
              <div className="overflow-x-auto">
                <table className="w-full text-left text-sm">
                  <thead className="text-muted-foreground">
                    <tr>
                      <th className="pb-4 font-medium">Nama Tenant</th>
                      <th className="pb-4 font-medium text-right">Total Pesanan</th>
                      <th className="pb-4 font-medium text-right">Total Pendapatan</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border/50">
                    {data.breakdown.map((item: any) => (
                      <tr key={item.tenant_id}>
                        <td className="py-4 font-medium">{item.tenant_name}</td>
                        <td className="py-4 text-right">{item.total_orders}</td>
                        <td className="py-4 text-right font-medium">{formatRupiah(item.revenue)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}
        </div>
      )}
    </DashboardShell>
  );
}

function Stat({ icon, label, value }: { icon: React.ReactNode; label: string; value: string }) {
  return (
    <div className="glass-strong rounded-3xl p-8 transition-all hover:scale-[1.01]">
      <div className="flex items-center gap-3 text-burgundy">{icon}<span className="text-sm font-bold uppercase tracking-widest">{label}</span></div>
      <div className="mt-4 font-display text-4xl font-bold">{value}</div>
    </div>
  );
}
