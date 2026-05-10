import { createFileRoute, Navigate } from "@tanstack/react-router";

export const Route = createFileRoute("/tenant/")({
  component: () => <Navigate to="/tenant/menus" />,
});
