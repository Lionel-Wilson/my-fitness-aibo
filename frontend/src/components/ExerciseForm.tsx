"use client";

import { useState } from "react";
import { Sheet } from "./Sheet";
import { api } from "@/lib/api";
import type { Exercise } from "@/lib/types";

type Draft = {
  name: string;
  targetSets: string;
  repLow: string;
  repHigh: string;
  rpeLow: string;
  rpeHigh: string;
  restSeconds: string;
  instructions: string;
  orGroup: string;
  isOptional: boolean;
  isUnilateral: boolean;
};

function toDraft(e?: Exercise): Draft {
  const s = (n: number | null | undefined) => (n == null ? "" : String(n));
  return {
    name: e?.name ?? "",
    targetSets: s(e?.targetSets),
    repLow: s(e?.repLow),
    repHigh: s(e?.repHigh),
    rpeLow: s(e?.rpeLow),
    rpeHigh: s(e?.rpeHigh),
    restSeconds: s(e?.restSeconds),
    instructions: e?.instructions ?? "",
    orGroup: e?.orGroup ?? "",
    isOptional: e?.isOptional ?? false,
    isUnilateral: e?.isUnilateral ?? false,
  };
}

function numOrNull(s: string): number | null {
  if (s.trim() === "") return null;
  const n = Number(s);
  return Number.isNaN(n) ? null : n;
}

interface Props {
  open: boolean;
  onClose: () => void;
  onSaved: () => void;
  workoutId?: string; // for create
  exercise?: Exercise; // for edit
}

// Bottom-sheet form to create or edit an exercise, including the per-exercise
// prescription and the static how-to instructions.
export function ExerciseForm({ open, onClose, onSaved, workoutId, exercise }: Props) {
  const [d, setD] = useState<Draft>(toDraft(exercise));
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const set = (k: keyof Draft, v: string | boolean) => setD((p) => ({ ...p, [k]: v }));

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    setBusy(true);
    setError(null);
    const body = {
      name: d.name,
      targetSets: numOrNull(d.targetSets),
      repLow: numOrNull(d.repLow),
      repHigh: numOrNull(d.repHigh),
      rpeLow: numOrNull(d.rpeLow),
      rpeHigh: numOrNull(d.rpeHigh),
      restSeconds: numOrNull(d.restSeconds),
      instructions: d.instructions,
      orGroup: d.orGroup,
      isOptional: d.isOptional,
      isUnilateral: d.isUnilateral,
    };
    try {
      if (exercise) await api.patch(`/exercises/${exercise.id}`, body);
      else await api.post(`/workouts/${workoutId}/exercises`, body);
      onSaved();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed");
    } finally {
      setBusy(false);
    }
  };

  return (
    <Sheet open={open} onClose={onClose} title={exercise ? "Edit exercise" : "New exercise"}>
      <form onSubmit={submit} className="space-y-4">
        <div>
          <label className="label">Exercise name</label>
          <input className="input" value={d.name} onChange={(e) => set("name", e.target.value)} placeholder="Front Squat" required />
        </div>

        <div className="grid grid-cols-3 gap-3">
          <div>
            <label className="label">Sets</label>
            <input className="input" inputMode="numeric" value={d.targetSets} onChange={(e) => set("targetSets", e.target.value)} placeholder="2" />
          </div>
          <div>
            <label className="label">Reps low</label>
            <input className="input" inputMode="numeric" value={d.repLow} onChange={(e) => set("repLow", e.target.value)} placeholder="6" />
          </div>
          <div>
            <label className="label">Reps high</label>
            <input className="input" inputMode="numeric" value={d.repHigh} onChange={(e) => set("repHigh", e.target.value)} placeholder="8" />
          </div>
        </div>

        <div className="grid grid-cols-3 gap-3">
          <div>
            <label className="label">RPE low</label>
            <input className="input" inputMode="decimal" value={d.rpeLow} onChange={(e) => set("rpeLow", e.target.value)} placeholder="7" />
          </div>
          <div>
            <label className="label">RPE high</label>
            <input className="input" inputMode="decimal" value={d.rpeHigh} onChange={(e) => set("rpeHigh", e.target.value)} placeholder="8" />
          </div>
          <div>
            <label className="label">Rest (s)</label>
            <input className="input" inputMode="numeric" value={d.restSeconds} onChange={(e) => set("restSeconds", e.target.value)} placeholder="150" />
          </div>
        </div>

        <div>
          <label className="label">Notes / tips (how to perform)</label>
          <textarea className="input min-h-24" value={d.instructions} onChange={(e) => set("instructions", e.target.value)} placeholder="Form cues, progression rules…" />
        </div>

        <div>
          <label className="label">Alternative group (optional)</label>
          <input className="input" value={d.orGroup} onChange={(e) => set("orGroup", e.target.value)} placeholder="e.g. squat — links A OR B options" />
        </div>

        <label className="flex items-center gap-3 text-sm">
          <input type="checkbox" checked={d.isOptional} onChange={(e) => set("isOptional", e.target.checked)} className="h-5 w-5 accent-accent" />
          Optional exercise
        </label>

        <label className="flex items-start gap-3 text-sm">
          <input type="checkbox" checked={d.isUnilateral} onChange={(e) => set("isUnilateral", e.target.checked)} className="h-5 w-5 accent-accent mt-0.5" />
          <span>
            Unilateral (per side)
            <span className="block text-xs text-muted">
              Log left & right independently — e.g. lateral raises, split squats.
            </span>
          </span>
        </label>

        {error && <p className="text-red-400 text-sm">{error}</p>}
        <button className="btn-primary w-full" disabled={busy}>
          {busy ? "…" : exercise ? "Save changes" : "Add exercise"}
        </button>
      </form>
    </Sheet>
  );
}
