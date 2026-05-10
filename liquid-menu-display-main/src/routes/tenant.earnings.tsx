import { useEffect, useState } from "react";
import { createFileRoute } from "@tanstack/react-router";
import { DashboardShell } from "@/components/DashboardShell";
import { UtensilsCrossed, ClipboardList, Wallet, TrendingUp, Receipt, Filter, Loader2, Download } from "lucide-react";
import { api } from "@/lib/api";
import { formatRupiah } from "@/stores/cart";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";

const NAV = [
  { to: "/tenant/menus", label: "Menus", icon: <UtensilsCrossed className="h-4 w-4" /> },
  { to: "/tenant/orders", label: "Orders", icon: <ClipboardList className="h-4 w-4" /> },
  { to: "/tenant/earnings", label: "Earnings", icon: <Wallet className="h-4 w-4" /> },
];

export const Route = createFileRoute("/tenant/earnings")({
  component: Earnings,
});

function Earnings() {
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const [dates, setDates] = useState({ start: "", end: "" });

  const fetchReport = async () => {
    setLoading(true);
    try {
      let url = "/tenant/reports/income";
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
  }, []); // Initial load without dates

  return (
    <DashboardShell title="Tenant" role="tenant" nav={NAV}>
      <div className="mb-8 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="font-display text-3xl font-bold">Laporan Pemasukan</h1>
          <p className="text-sm text-muted-foreground">Monitor pendapatan stand Anda.</p>
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
            <Stat icon={<TrendingUp />} label="Total Pendapatan" value={formatRupiah(data?.total_revenue || 0)} />
            <Stat icon={<Receipt />} label="Total Pesanan" value={String(data?.total_orders || 0)} />
          </div>

          {data?.history && data.history.length > 0 && (
            <div className="glass-strong rounded-3xl p-8 bg-card">
              <div className="mb-6 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                <h2 className="font-display text-2xl font-bold">Riwayat Penjualan Menu</h2>
                <Button onClick={() => window.print()} className="bg-burgundy hover:bg-burgundy/90 text-white rounded-xl shadow-lg shadow-burgundy/20 print:hidden">
                  <Download className="mr-2 h-4 w-4" /> Download PDF
                </Button>
              </div>
              <div className="overflow-x-auto">
                <table className="w-full text-left text-sm">
                  <thead className="text-muted-foreground border-b border-border/50">
                    <tr>
                      <th className="pb-4 font-medium w-16">No</th>
                      <th className="pb-4 font-medium">Menu</th>
                      <th className="pb-4 font-medium text-center">Jumlah</th>
                      <th className="pb-4 font-medium">Tanggal Pembelian</th>
                      <th className="pb-4 font-medium text-right">Pendapatan</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border/50">
                    {data.history.map((item: any, idx: number) => (
                      <tr key={idx} className="hover:bg-muted/30 transition-colors">
                        <td className="py-4 text-muted-foreground">{idx + 1}</td>
                        <td className="py-4 font-medium">{item.menu_name}</td>
                        <td className="py-4 text-center">{item.quantity_sold}</td>
                        <td className="py-4 text-muted-foreground">{item.purchase_date}</td>
                        <td className="py-4 text-right font-medium text-burgundy">{formatRupiah(item.revenue)}</td>
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
