import { useEffect, useState } from "react";
import { createFileRoute } from "@tanstack/react-router";
import { DashboardShell } from "@/components/DashboardShell";
import { UtensilsCrossed, ClipboardList, Wallet, TrendingUp, Receipt } from "lucide-react";
import { api } from "@/lib/api";
import { formatRupiah } from "@/stores/cart";
import { Skeleton } from "@/components/ui/skeleton";
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
  useEffect(() => {
    api<any>("/tenant/earnings").then((d) => setData(d.data || d)).catch((e) => toast.error(e.message));
  }, []);

  return (
    <DashboardShell title="Tenant" role="tenant" nav={NAV}>
      <h1 className="mb-8 font-display text-3xl font-bold">Earnings</h1>
      {!data ? (
        <div className="grid gap-5 sm:grid-cols-2"><Skeleton className="h-40 rounded-2xl" /><Skeleton className="h-40 rounded-2xl" /></div>
      ) : (
        <div className="grid gap-5 sm:grid-cols-2">
          <Stat icon={<TrendingUp />} label="Total Revenue" value={formatRupiah(data.total_revenue || data.total_earnings || 0)} />
          <Stat icon={<Receipt />} label="Total Orders" value={String(data.total_orders || 0)} />
        </div>
      )}
    </DashboardShell>
  );
}

function Stat({ icon, label, value }: { icon: React.ReactNode; label: string; value: string }) {
  return (
    <div className="glass-strong rounded-3xl p-8">
      <div className="flex items-center gap-3 text-burgundy">{icon}<span className="text-sm font-medium">{label}</span></div>
      <div className="mt-3 font-display text-4xl font-bold">{value}</div>
    </div>
  );
}
