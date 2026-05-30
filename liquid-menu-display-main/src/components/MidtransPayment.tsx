import { useEffect } from "react";
import { api } from "@/lib/api";
import { toast } from "sonner";
import { MidtransResult } from "@/types/midtrans";

interface MidtransPaymentProps {
  token: string;
  orderId: string;
  onSuccess: (result: MidtransResult) => void;
  onPending?: (result: MidtransResult) => void;
  onError?: (error: any) => void;
  onClose?: () => void;
}

export default function MidtransPayment({
  token,
  orderId,
  onSuccess,
  onPending,
  onError,
  onClose,
}: MidtransPaymentProps) {
  useEffect(() => {
    // Load Snap.js
    const snapUrl = import.meta.env.VITE_MIDTRANS_SNAP_URL || "https://app.sandbox.midtrans.com/snap/snap.js";
    const clientKey = import.meta.env.VITE_MIDTRANS_CLIENT_KEY;

    // Check if script already exists
    let script = document.getElementById("midtrans-snap") as HTMLScriptElement;
    
    if (!script) {
      script = document.createElement("script");
      script.id = "midtrans-snap";
      script.src = snapUrl;
      script.setAttribute("data-client-key", clientKey);
      script.async = true;
      document.body.appendChild(script);
    }

    const triggerSnap = () => {
      if (window.snap) {
        window.snap.pay(token, {
          onSuccess: async (result) => {
            console.log("Midtrans Success:", result);
            // Verify locally just in case webhook is slow
            try {
              await api("/midtrans/verify", {
                method: "POST",
                body: JSON.stringify({ order_id: orderId }),
              });
            } catch (err) {
              console.error("Local verification failed:", err);
            }
            onSuccess(result);
          },
          onPending: (result) => {
            console.log("Midtrans Pending:", result);
            if (onPending) onPending(result);
          },
          onError: (result) => {
            console.error("Midtrans Error:", result);
            toast.error("Pembayaran gagal. Silakan coba lagi.");
            if (onError) onError(result);
          },
          onClose: () => {
            console.log("Midtrans Popup Closed");
            if (onClose) onClose();
          },
        });
      } else {
        console.error("Snap is not loaded yet");
        setTimeout(triggerSnap, 500);
      }
    };

    if (window.snap) {
      triggerSnap();
    } else {
      script.addEventListener("load", triggerSnap);
    }

    return () => {
      script.removeEventListener("load", triggerSnap);
    };
  }, [token, orderId, onSuccess, onPending, onError, onClose]);

  return null;
}
