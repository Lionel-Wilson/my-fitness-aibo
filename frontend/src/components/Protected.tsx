"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth";

// Wraps authenticated pages: shows a loader while resolving the session and
// redirects to /login when there is no user.
export function Protected({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && !user) router.replace("/login");
  }, [loading, user, router]);

  if (loading || !user) {
    return (
      <div className="flex-1 flex items-center justify-center text-muted">Loading…</div>
    );
  }
  return <>{children}</>;
}
