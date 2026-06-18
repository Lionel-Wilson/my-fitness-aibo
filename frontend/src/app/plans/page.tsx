"use client";

import { useState } from "react";
import Link from "next/link";
import { Protected } from "@/components/Protected";
import { PageHeader } from "@/components/PageHeader";
import { BottomNav } from "@/components/BottomNav";
import { Sheet } from "@/components/Sheet";
import { GettingStarted } from "@/components/GettingStarted";
import { useFetch } from "@/lib/useFetch";
import { api } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import type { Plan } from "@/lib/types";

export default function PlansPage() {
  return (
    <Protected>
      <PlansInner />
    </Protected>
  );
}

function PlansInner() {
  const { logout } = useAuth();
  const { data: plans, loading, refetch } = useFetch<Plan[]>("/plans");
  const [open, setOpen] = useState(false);

  return (
    <>
      <PageHeader
        title="My Plans"
        right={
          <button onClick={logout} className="text-sm text-muted active:text-white">
            Log out
          </button>
        }
      />

      <main className="flex-1 px-4 py-4 space-y-3">
        {loading && <p className="text-muted">Loading…</p>}

        {/* Onboarding: expanded for new users, a collapsible refresher otherwise. */}
        {plans && <GettingStarted defaultOpen={plans.length === 0} />}

        {plans && plans.length === 0 && (
          <div className="card text-center text-muted py-6">
            <p className="mb-1 text-2xl">🗂️</p>
            <p>No plans yet.</p>
            <p className="text-sm">Start with step 1 below.</p>
          </div>
        )}
        {plans?.map((p) => (
          <Link key={p.id} href={`/plans/${p.id}`} className="card block active:bg-surface2">
            <div className="flex items-center justify-between">
              <div className="min-w-0">
                <p className="font-semibold truncate">{p.name}</p>
                {p.quality && <p className="text-sm text-muted capitalize">{p.quality}</p>}
              </div>
              {p.isActive && (
                <span className="text-xs bg-accent/15 text-accent px-2 py-1 rounded-full shrink-0">
                  Active
                </span>
              )}
            </div>
          </Link>
        ))}

        <button onClick={() => setOpen(true)} className="btn-primary w-full mt-2">
          + New plan
        </button>
      </main>

      <BottomNav />

      <PlanForm
        open={open}
        onClose={() => setOpen(false)}
        onSaved={() => {
          setOpen(false);
          refetch();
        }}
      />
    </>
  );
}

function PlanForm({ open, onClose, onSaved }: { open: boolean; onClose: () => void; onSaved: () => void }) {
  const [name, setName] = useState("");
  const [quality, setQuality] = useState("");
  const [cycleLabel, setCycleLabel] = useState("Cycle");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    setBusy(true);
    setError(null);
    try {
      await api.post<Plan>("/plans", { name, quality, cycleLabel });
      setName("");
      setQuality("");
      setCycleLabel("Cycle");
      onSaved();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed");
    } finally {
      setBusy(false);
    }
  };

  return (
    <Sheet open={open} onClose={onClose} title="New plan">
      <form onSubmit={submit} className="space-y-4">
        <div>
          <label className="label">Plan name</label>
          <input className="input" value={name} onChange={(e) => setName(e.target.value)} placeholder="Building muscle" required />
        </div>
        <div>
          <label className="label">Quality (optional)</label>
          <input className="input" value={quality} onChange={(e) => setQuality(e.target.value)} placeholder="hypertrophy / strength / power" />
        </div>
        <div>
          <label className="label">What do you call a round of this plan?</label>
          <input className="input" value={cycleLabel} onChange={(e) => setCycleLabel(e.target.value)} placeholder="Cycle" />
          <p className="text-xs text-muted mt-1">Used as column labels (e.g. Cycle, Week, Block).</p>
        </div>
        {error && <p className="text-red-400 text-sm">{error}</p>}
        <button className="btn-primary w-full" disabled={busy}>
          {busy ? "…" : "Create plan"}
        </button>
      </form>
    </Sheet>
  );
}
