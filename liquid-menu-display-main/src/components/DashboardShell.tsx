import { Link, useNavigate, useLocation } from "@tanstack/react-router";
import { ChefHat, LogOut } from "lucide-react";
import { useAuth } from "@/stores/auth";
import { useEffect, type ReactNode } from "react";
import { Button } from "@/components/ui/button";

interface NavItem { to: string; label: string; icon: ReactNode }

export function DashboardShell({
  title,
  role,
  nav,
  children,
}: {
  title: string;
  role: "admin" | "tenant";
  nav: NavItem[];
  children: ReactNode;
}) {
  const auth = useAuth();
  const navigate = useNavigate();
  const loc = useLocation();

  useEffect(() => {
    if (typeof window === "undefined" || !auth._hasHydrated) return;
    
    if (!auth.token || auth.role !== role) {
      navigate({ to: "/login" });
    }
  }, [auth.token, auth.role, role, navigate, auth._hasHydrated]);

  return (
    <div className="flex min-h-screen">
      <aside className="glass-strong sticky top-0 hidden h-screen w-64 flex-col gap-2 border-r p-6 md:flex print:hidden">
        <Link to="/" className="mb-8 flex items-center gap-2 font-display text-xl font-bold text-burgundy">
          <ChefHat className="h-6 w-6" /> Saji
        </Link>
        <div className="mb-2 px-2 text-xs uppercase tracking-wider text-muted-foreground">{title}</div>
        <nav className="flex flex-col gap-1">
          {nav.map((n) => {
            const active = loc.pathname === n.to;
            return (
              <Link
                key={n.to}
                to={n.to}
                className={`flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-medium transition ${
                  active ? "bg-burgundy text-burgundy-foreground" : "hover:bg-muted"
                }`}
              >
                {n.icon} {n.label}
              </Link>
            );
          })}
        </nav>
        <div className="mt-auto space-y-2">
          <div className="rounded-xl border bg-white/50 p-3 text-xs">
            <div className="text-muted-foreground">Signed in as</div>
            <div className="truncate font-medium">{auth.email}</div>
          </div>
          <Button variant="outline" className="w-full" onClick={() => { auth.logout(); navigate({ to: "/login" }); }}>
            <LogOut className="h-4 w-4" /> Logout
          </Button>
        </div>
      </aside>
      <main className="flex-1 px-6 py-8 md:px-10">{children}</main>
    </div>
  );
}
