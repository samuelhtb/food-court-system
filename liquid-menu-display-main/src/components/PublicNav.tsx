import { Link } from "@tanstack/react-router";
import { ShoppingBag, Search, ChefHat } from "lucide-react";
import { useCart } from "@/stores/cart";
import { Button } from "@/components/ui/button";

export function PublicNav({ onCart }: { onCart: () => void }) {
  const count = useCart((s) => s.count());
  return (
    <header className="sticky top-0 z-40 px-4 pt-4">
      <nav className="glass-strong mx-auto flex max-w-6xl items-center justify-between rounded-full px-5 py-3">
        <Link to="/" className="flex items-center gap-2 font-display text-xl font-bold text-burgundy">
          <ChefHat className="h-6 w-6" />
          <span>Saji</span>
        </Link>
        <div className="hidden items-center gap-6 text-sm font-medium md:flex">
          <Link to="/" className="hover:text-burgundy transition">Menu</Link>
          <Link to="/track-order" className="hover:text-burgundy transition">Track Order</Link>
          <Link to="/login" className="hover:text-burgundy transition">Sign in</Link>
        </div>
        <div className="flex items-center gap-2">
          <Button onClick={onCart} variant="default" size="sm" className="relative rounded-full">
            <ShoppingBag className="h-4 w-4" />
            <span className="ml-2 hidden sm:inline">Cart</span>
            {count > 0 && (
              <span className="absolute -right-1 -top-1 flex h-5 w-5 items-center justify-center rounded-full bg-foreground text-[10px] font-bold text-background">
                {count}
              </span>
            )}
          </Button>
        </div>
      </nav>
    </header>
  );
}
