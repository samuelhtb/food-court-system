import { useState } from "react";
import { createFileRoute, useNavigate, Link } from "@tanstack/react-router";
import { api } from "@/lib/api";
import { useAuth } from "@/stores/auth";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ChefHat } from "lucide-react";
import { toast } from "sonner";

export const Route = createFileRoute("/login")({
  component: Login,
});

function Login() {
  const nav = useNavigate();
  const setAuth = useAuth((s) => s.setAuth);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      const res: any = await api("/login", {
        method: "POST",
        auth: false,
        body: JSON.stringify({ email, password }),
      });
      const token = res.token || res.access_token || res.data?.token;
      const role = res.role || res.user?.role || res.data?.role;
      if (!token) throw new Error("No token in response");
      setAuth(token, role, email);
      toast.success("Welcome back");
      nav({ to: role === "admin" ? "/admin" : "/tenant" });
    } catch (err: any) {
      toast.error(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center px-4">
      <div className="glass-strong w-full max-w-md rounded-3xl p-8 shadow-2xl">
        <Link to="/" className="mb-6 flex items-center justify-center gap-2 font-display text-2xl font-bold text-burgundy">
          <ChefHat className="h-7 w-7" /> Saji
        </Link>
        <h1 className="text-center font-display text-3xl font-bold">Sign in</h1>
        <p className="mt-1 text-center text-sm text-muted-foreground">Tenant & Admin access</p>
        <form onSubmit={submit} className="mt-8 space-y-4">
          <div className="space-y-2">
            <Label>Email</Label>
            <Input type="email" required value={email} onChange={(e) => setEmail(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>Password</Label>
            <Input type="password" required value={password} onChange={(e) => setPassword(e.target.value)} />
          </div>
          <Button type="submit" disabled={loading} className="w-full" size="lg">
            {loading ? "Signing in..." : "Sign in"}
          </Button>
        </form>
        <Link to="/" className="mt-6 block text-center text-xs text-muted-foreground hover:text-foreground">
          ← Back to menu
        </Link>
      </div>
    </div>
  );
}
