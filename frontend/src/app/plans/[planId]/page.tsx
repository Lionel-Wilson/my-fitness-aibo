"use client";

import { useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { Protected } from "@/components/Protected";
import { PageHeader } from "@/components/PageHeader";
import { Sheet } from "@/components/Sheet";
import { useFetch } from "@/lib/useFetch";
import { api } from "@/lib/api";
import type { Plan, Workout, Cycle } from "@/lib/types";

export default function PlanPage() {
  return (
    <Protected>
      <PlanInner />
    </Protected>
  );
}

function PlanInner() {
  const { planId } = useParams<{ planId: string }>();
  const { data: plan } = useFetch<Plan>(`/plans/${planId}`);
  const { data: workouts, refetch: refetchWorkouts } = useFetch<Workout[]>(`/plans/${planId}/workouts`);
  const { data: cycles, refetch: refetchCycles } = useFetch<Cycle[]>(`/plans/${planId}/cycles`);
  const [open, setOpen] = useState(false);

  const label = plan?.cycleLabel || "Cycle";
  const current = cycles?.[0];

  const startCycle = async () => {
    const next = (cycles?.[0]?.cycleNumber ?? 0) + 1;
    await api.post(`/plans/${planId}/cycles`, { label: `${label[0]?.toUpperCase()}${next}` });
    refetchCycles();
  };

  return (
    <>
      <PageHeader title={plan?.name || "Plan"} subtitle={plan?.quality} back />

      <main className="flex-1 px-4 py-4 space-y-4">
        <div className="card flex items-center justify-between">
          <div>
            <p className="text-xs text-muted">Current {label.toLowerCase()}</p>
            <p className="text-xl font-semibold">
              {current ? current.label || `${label} ${current.cycleNumber}` : "—"}
            </p>
            <p className="text-xs text-muted">{cycles?.length ?? 0} total</p>
          </div>
          <button onClick={startCycle} className="btn-secondary text-sm">
            + New {label.toLowerCase()}
          </button>
        </div>

        <div>
          <div className="flex items-center justify-between mb-2">
            <h2 className="text-sm uppercase tracking-wide text-muted">Workouts</h2>
            <button onClick={() => setOpen(true)} className="text-accent text-sm">
              + Add
            </button>
          </div>
          <div className="space-y-2">
            {workouts?.length === 0 && (
              <div className="card text-sm text-muted">
                <p className="text-white font-medium mb-1">Next: add your workouts</p>
                <p>Add the training days in this plan (e.g. “Monday — Legs”). Then open a
                  workout to add its exercises, start a cycle, and log your sessions.</p>
              </div>
            )}
            {workouts?.map((w) => (
              <Link key={w.id} href={`/workouts/${w.id}`} className="card block active:bg-surface2">
                <p className="font-medium">{w.name}</p>
                {w.dayLabel && <p className="text-sm text-muted">{w.dayLabel}</p>}
              </Link>
            ))}
          </div>
        </div>
      </main>

      <WorkoutForm
        planId={planId}
        open={open}
        onClose={() => setOpen(false)}
        onSaved={() => {
          setOpen(false);
          refetchWorkouts();
        }}
      />
    </>
  );
}

function WorkoutForm({ planId, open, onClose, onSaved }: { planId: string; open: boolean; onClose: () => void; onSaved: () => void }) {
  const [name, setName] = useState("");
  const [dayLabel, setDayLabel] = useState("");
  const [busy, setBusy] = useState(false);

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    setBusy(true);
    try {
      await api.post<Workout>(`/plans/${planId}/workouts`, { name, dayLabel });
      setName("");
      setDayLabel("");
      onSaved();
    } finally {
      setBusy(false);
    }
  };

  return (
    <Sheet open={open} onClose={onClose} title="New workout">
      <form onSubmit={submit} className="space-y-4">
        <div>
          <label className="label">Workout name</label>
          <input className="input" value={name} onChange={(e) => setName(e.target.value)} placeholder="Monday — LEGS" required />
        </div>
        <div>
          <label className="label">Day label (optional)</label>
          <input className="input" value={dayLabel} onChange={(e) => setDayLabel(e.target.value)} placeholder="Monday / Day 1" />
        </div>
        <button className="btn-primary w-full" disabled={busy}>
          {busy ? "…" : "Add workout"}
        </button>
      </form>
    </Sheet>
  );
}
