import { useEffect, useState } from "react";
import { createFileRoute } from "@tanstack/react-router";
import { DashboardShell } from "@/components/DashboardShell";
import { UtensilsCrossed, ClipboardList, Wallet, Pencil, Trash2, Plus, Upload, Loader2, Image as ImageIcon } from "lucide-react";
import { api } from "@/lib/api";
import { formatRupiah, type Menu } from "@/stores/cart";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger, DialogFooter } from "@/components/ui/dialog";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";
import { supabase } from "@/integrations/supabase/client";

const NAV = [
  { to: "/tenant/menus", label: "Menus", icon: <UtensilsCrossed className="h-4 w-4" /> },
  { to: "/tenant/orders", label: "Orders", icon: <ClipboardList className="h-4 w-4" /> },
  { to: "/tenant/earnings", label: "Earnings", icon: <Wallet className="h-4 w-4" /> },
];

export const Route = createFileRoute("/tenant/menus")({
  component: TenantMenus,
});

interface FormData {
  name: string;
  description: string;
  price: number;
  stock: number;
  image_url: string;
  is_available: boolean;
}

const empty: FormData = {
  name: "",
  description: "",
  price: 0,
  stock: 0,
  image_url: "",
  is_available: true,
};

function TenantMenus() {
  const [menus, setMenus] = useState<Menu[] | null>(null);
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<Menu | null>(null);
  const [form, setForm] = useState<FormData>(empty);
  const [imageFile, setImageFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [isUploading, setIsUploading] = useState(false);

  const load = async () => {
    try {
      const res = await api<any>("/tenant/menus");
      setMenus(res.data || res.menus || (Array.isArray(res) ? res : []));
    } catch (e: any) {
      toast.error(e.message);
      setMenus([]);
    }
  };

  useEffect(() => {
    load();
  }, []);

  // Handle image preview
  useEffect(() => {
    if (!imageFile) {
      setPreviewUrl(null);
      return;
    }
    const url = URL.createObjectURL(imageFile);
    setPreviewUrl(url);
    return () => URL.revokeObjectURL(url);
  }, [imageFile]);

  const openCreate = () => {
    setEditing(null);
    setForm(empty);
    setImageFile(null);
    setOpen(true);
  };

  const openEdit = (m: Menu) => {
    setEditing(m);
    setForm({
      name: m.name,
      description: m.description || "",
      price: m.price,
      stock: m.stock,
      image_url: m.image_url || "",
      is_available: m.is_available ?? true,
    });
    setImageFile(null);
    setOpen(true);
  };

  const save = async () => {
    if (!form.name || form.price <= 0) {
      toast.error("Please fill name and valid price");
      return;
    }

    setIsUploading(true);
    try {
      let finalImageUrl = form.image_url;

      // 1. Upload to Supabase Storage if new image selected
      if (imageFile) {
        const fileExt = imageFile.name.split(".").pop();
        const fileName = `${Date.now()}-${Math.random().toString(36).substring(7)}.${fileExt}`;
        const filePath = `menus/${fileName}`;

        const { error: uploadError } = await supabase.storage
          .from("menu-images")
          .upload(filePath, imageFile);

        if (uploadError) throw uploadError;

        const { data: { publicUrl } } = supabase.storage
          .from("menu-images")
          .getPublicUrl(filePath);

        finalImageUrl = publicUrl;
      }

      // 2. Save to Golang API
      const payload = { ...form, image_url: finalImageUrl };

      if (editing) {
        await api(`/menus/${editing.id}`, {
          method: "PUT",
          body: JSON.stringify(payload),
        });
        toast.success("Menu updated successfully");
      } else {
        const { is_available, ...createPayload } = payload;
        await api("/menus", {
          method: "POST",
          body: JSON.stringify(createPayload),
        });
        toast.success("Menu created successfully");
      }

      setOpen(false);
      setImageFile(null);
      load();
    } catch (e: any) {
      console.error("Save error:", e);
      toast.error(e.message || "Failed to save menu");
    } finally {
      setIsUploading(false);
    }
  };

  const del = async (id: string) => {
    if (!confirm("Are you sure you want to delete this menu item?")) return;
    try {
      await api(`/menus/${id}`, { method: "DELETE" });
      toast.success("Menu deleted");
      load();
    } catch (e: any) {
      toast.error(e.message);
    }
  };

  return (
    <DashboardShell title="Menu Management" role="tenant" nav={NAV}>
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="font-display text-3xl font-bold">My Menus</h1>
          <p className="text-sm text-muted-foreground">Add and manage your food items.</p>
        </div>
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogTrigger asChild>
            <Button onClick={openCreate} className="rounded-xl">
              <Plus className="mr-2 h-4 w-4" /> New Menu
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[425px]">
            <DialogHeader>
              <DialogTitle>{editing ? "Edit Menu Item" : "Create New Menu Item"}</DialogTitle>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="space-y-2">
                <Label htmlFor="name">Item Name</Label>
                <Input
                  id="name"
                  placeholder="e.g. Nasi Goreng Spesial"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="description">Description</Label>
                <Textarea
                  id="description"
                  placeholder="Tell customers about this dish..."
                  value={form.description}
                  onChange={(e) => setForm({ ...form, description: e.target.value })}
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="price">Price (Rp)</Label>
                  <Input
                    id="price"
                    type="number"
                    min={0}
                    value={form.price}
                    onChange={(e) => setForm({ ...form, price: +e.target.value })}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="stock">Stock</Label>
                  <Input
                    id="stock"
                    type="number"
                    min={0}
                    value={form.stock}
                    onChange={(e) => setForm({ ...form, stock: +e.target.value })}
                  />
                </div>
              </div>
              <div className="space-y-2">
                <Label>Item Image</Label>
                <div className="flex flex-col items-center gap-4">
                  {(previewUrl || form.image_url) ? (
                    <div className="relative aspect-video w-full overflow-hidden rounded-xl border">
                      <img
                        src={previewUrl || form.image_url}
                        alt="Preview"
                        className="h-full w-full object-cover"
                      />
                      <Button
                        type="button"
                        variant="secondary"
                        size="sm"
                        className="absolute bottom-2 right-2"
                        onClick={() => document.getElementById("image-upload")?.click()}
                      >
                        Change Image
                      </Button>
                    </div>
                  ) : (
                    <div
                      className="flex aspect-video w-full cursor-pointer flex-col items-center justify-center rounded-xl border-2 border-dashed border-muted-foreground/25 transition-colors hover:border-muted-foreground/50 hover:bg-muted/50"
                      onClick={() => document.getElementById("image-upload")?.click()}
                    >
                      <Upload className="mb-2 h-8 w-8 text-muted-foreground" />
                      <p className="text-xs text-muted-foreground">Click to upload image</p>
                    </div>
                  )}
                  <Input
                    id="image-upload"
                    type="file"
                    accept="image/*"
                    className="hidden"
                    onChange={(e) => setImageFile(e.target.files?.[0] || null)}
                  />
                </div>
              </div>
              {editing && (
                <div className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    id="available"
                    checked={form.is_available}
                    onChange={(e) => setForm({ ...form, is_available: e.target.checked })}
                    className="h-4 w-4 rounded border-gray-300 text-burgundy focus:ring-burgundy"
                  />
                  <Label htmlFor="available">Item is currently available</Label>
                </div>
              )}
            </div>
            <DialogFooter>
              <Button onClick={save} disabled={isUploading} className="w-full rounded-xl bg-burgundy hover:bg-burgundy/90">
                {isUploading ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Processing...
                  </>
                ) : (
                  "Save Menu Item"
                )}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      {!menus ? (
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <Skeleton key={i} className="h-72 w-full rounded-3xl" />
          ))}
        </div>
      ) : menus.length === 0 ? (
        <div className="glass flex flex-col items-center justify-center rounded-3xl p-20 text-center">
          <div className="mb-4 rounded-full bg-muted p-4 text-muted-foreground">
            <UtensilsCrossed className="h-10 w-10" />
          </div>
          <h2 className="text-xl font-bold">No menus found</h2>
          <p className="max-w-xs text-muted-foreground">
            Start by adding your first menu item to show it to customers.
          </p>
          <Button onClick={openCreate} variant="outline" className="mt-6 rounded-xl">
            Add Your First Item
          </Button>
        </div>
      ) : (
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {menus.map((m) => (
            <div key={m.id} className="glass group overflow-hidden rounded-3xl transition-all hover:shadow-xl">
              <div className="relative aspect-[16/10] overflow-hidden">
                {m.image_url ? (
                  <img
                    src={m.image_url}
                    alt={m.name}
                    className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-110"
                  />
                ) : (
                  <div className="flex h-full w-full items-center justify-center bg-muted">
                    <ImageIcon className="h-10 w-10 text-muted-foreground/20" />
                  </div>
                )}
                {!m.is_available && (
                  <div className="absolute inset-0 flex items-center justify-center bg-black/60 backdrop-blur-[2px]">
                    <span className="rounded-full bg-white/10 px-4 py-2 text-xs font-bold text-white ring-1 ring-white/20">
                      UNAVAILABLE
                    </span>
                  </div>
                )}
              </div>
              <div className="p-5">
                <div className="mb-2 flex items-start justify-between">
                  <h3 className="font-display text-lg font-bold leading-tight">{m.name}</h3>
                  <span className="text-base font-bold text-burgundy">{formatRupiah(m.price)}</span>
                </div>
                <p className="line-clamp-2 text-sm text-muted-foreground">{m.description || "No description provided."}</p>
                <div className="mt-6 flex items-center justify-between">
                  <div className="flex flex-col">
                    <span className="text-[10px] uppercase tracking-wider text-muted-foreground">Availability</span>
                    <span className="text-xs font-medium">
                      {m.stock > 0 ? `${m.stock} in stock` : "Out of stock"}
                    </span>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      size="icon"
                      variant="secondary"
                      className="h-9 w-9 rounded-xl"
                      onClick={() => openEdit(m)}
                    >
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                      size="icon"
                      variant="secondary"
                      className="h-9 w-9 rounded-xl hover:bg-destructive hover:text-destructive-foreground"
                      onClick={() => del(m.id)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
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
