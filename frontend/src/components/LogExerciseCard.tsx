"use client";

import { useState } from "react";
import { useFetch } from "@/lib/useFetch";
import { api } from "@/lib/api";
import { useSettings } from "@/lib/settings";
import { prescription } from "@/lib/format";
import type { Exercise, ExerciseLog, SetLog, Side } from "@/lib/types";

interface Props {
  exercise: Exercise;
  cycleId: string;
}

type Cell = { weight: string; reps: string };
type SetEntry = Record<string, Cell>; // keyed by side ("both" | "left" | "right")

const fieldKey = (i: number, side: string, field: keyof Cell) => `${i}-${side}-${field}`;

// One exercise within the log-session flow: shows the target, pre-fills any
// existing sets for the selected cycle, and saves the whole log in one call.
export function LogExerciseCard({ exercise, cycleId }: Props) {
  // Re-mounted per cycle (keyed by parent), so fetching once on mount is fine.
  const { data: logs, loading } = useFetch<ExerciseLog[]>(`/exercises/${exercise.id}/logs`);
  if (loading) return <div className="card text-muted text-sm">Loading {exercise.name}…</div>;
  const existing = logs?.find((l) => l.cycleId === cycleId);
  return <Inner exercise={exercise} cycleId={cycleId} existing={existing} />;
}

function Inner({ exercise, cycleId, existing }: Props & { existing?: ExerciseLog }) {
  const { autofillSets } = useSettings();
  const sides: Side[] = exercise.isUnilateral ? ["left", "right"] : ["both"];
  const emptyEntry = (): SetEntry =>
    Object.fromEntries(sides.map((s) => [s, { weight: "", reps: "" }]));

  const buildInitial = (): SetEntry[] => {
    if (existing?.sets.length) {
      const bySet = new Map<number, SetEntry>();
      existing.sets.forEach((s) => {
        if (!bySet.has(s.setIndex)) bySet.set(s.setIndex, emptyEntry());
        const side = sides.includes(s.side) ? s.side : sides[0];
        bySet.get(s.setIndex)![side] = {
          weight: s.weightKg?.toString() ?? "",
          reps: s.reps?.toString() ?? "",
        };
      });
      return [...bySet.keys()].sort((a, b) => a - b).map((k) => bySet.get(k)!);
    }
    return Array.from({ length: Math.max(1, exercise.targetSets ?? 2) }, emptyEntry);
  };

  const [rows, setRows] = useState<SetEntry[]>(buildInitial);
  const [note, setNote] = useState(existing?.note ?? "");
  const [status, setStatus] = useState<"idle" | "saving" | "saved" | "error">(
    existing ? "saved" : "idle"
  );
  const [badFields, setBadFields] = useState<Record<string, boolean>>({});
  const [errors, setErrors] = useState<string[]>([]);

  const setCell = (i: number, side: string, field: keyof Cell, v: string) => {
    setRows((prev) => {
      const oldVal = prev[i][side][field];
      const next = prev.map((row, j) =>
        j === i ? { ...row, [side]: { ...row[side], [field]: v } } : row
      );
      // Auto-fill: typing a weight propagates to later sets (same side) that are
      // still empty or were tracking the old value.
      if (autofillSets && field === "weight") {
        for (let j = i + 1; j < next.length; j++) {
          const cur = next[j][side].weight;
          if (cur === "" || cur === oldVal) {
            next[j] = { ...next[j], [side]: { ...next[j][side], weight: v } };
          }
        }
      }
      return next;
    });
    setBadFields((p) => (p[fieldKey(i, side, field)] ? { ...p, [fieldKey(i, side, field)]: false } : p));
  };
  const addSet = () =>
    setRows((p) => {
      const last = p[p.length - 1];
      const entry: SetEntry = Object.fromEntries(
        sides.map((s) => [
          s,
          // Copy the previous set's values when auto-fill is on; otherwise blank.
          autofillSets && last
            ? { weight: last[s]?.weight ?? "", reps: last[s]?.reps ?? "" }
            : { weight: "", reps: "" },
        ])
      );
      return [...p, entry];
    });
  const removeSet = (i: number) => setRows((p) => (p.length > 1 ? p.filter((_, j) => j !== i) : p));

  const save = async () => {
    const sets: SetLog[] = [];
    const bad: Record<string, boolean> = {};
    const messages: string[] = [];

    rows.forEach((entry, i) => {
      const n = i + 1;
      sides.forEach((side) => {
        const cell = entry[side];
        const w = cell.weight.trim();
        const reps = cell.reps.trim();
        if (w === "" && reps === "") return; // empty — ignored

        const tag = side === "both" ? `Set ${n}` : `Set ${n} (${side === "left" ? "L" : "R"})`;
        let weightKg: number | null = null;
        let repsVal: number | null = null;

        if (w !== "") {
          const v = Number(w);
          if (Number.isNaN(v) || v < 0) {
            bad[fieldKey(i, side, "weight")] = true;
            messages.push(`${tag}: weight must be a number.`);
          } else {
            weightKg = v;
          }
        }
        if (reps === "") {
          bad[fieldKey(i, side, "reps")] = true;
          messages.push(`${tag}: enter the reps you did.`);
        } else {
          const v = Number(reps);
          if (Number.isNaN(v) || v <= 0) {
            bad[fieldKey(i, side, "reps")] = true;
            messages.push(`${tag}: reps must be a positive number (decimals like 8.5 are fine).`);
          } else {
            repsVal = v;
          }
        }
        sets.push({ setIndex: n, side, weightKg, reps: repsVal, rpe: null, isDropSet: false });
      });
    });

    if (messages.length > 0) {
      setBadFields(bad);
      setErrors(messages);
      setStatus("idle");
      return;
    }
    if (sets.length === 0 && note.trim() === "") {
      setErrors(["Enter at least one set (weight and reps) before saving."]);
      setStatus("idle");
      return;
    }

    setBadFields({});
    setErrors([]);
    setStatus("saving");
    try {
      await api.put(`/exercises/${exercise.id}/logs/${cycleId}`, { note, sets });
      setStatus("saved");
    } catch (e) {
      setStatus("error");
      setErrors([e instanceof Error ? e.message : "Could not save. Please try again."]);
    }
  };

  const cell = (i: number, side: Side, field: keyof Cell, placeholder: string) => (
    <input
      className={`input py-2 text-center text-lg ${badFields[fieldKey(i, side, field)] ? "border-red-500" : ""}`}
      inputMode="decimal"
      value={rows[i][side][field]}
      onChange={(e) => setCell(i, side, field, e.target.value)}
      placeholder={placeholder}
    />
  );

  return (
    <div className="card space-y-3">
      <div className="flex items-center justify-between">
        <div className="min-w-0">
          <p className="font-medium truncate">{exercise.name}</p>
          <p className="text-xs text-muted">
            {[prescription(exercise), exercise.isUnilateral ? "per side" : ""]
              .filter(Boolean)
              .join(" · ")}
          </p>
        </div>
        {status === "saved" && <span className="text-xs text-green-400 shrink-0">✓ saved</span>}
      </div>

      <div className="space-y-3">
        {!exercise.isUnilateral && (
          <div className="grid grid-cols-[2rem_1fr_1fr_1.5rem] gap-2 text-xs text-muted px-1">
            <span>Set</span>
            <span>Weight (kg)</span>
            <span>Reps</span>
            <span />
          </div>
        )}

        {rows.map((_, i) =>
          exercise.isUnilateral ? (
            <div key={i} className="rounded-xl border border-border p-2 space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Set {i + 1}</span>
                <button onClick={() => removeSet(i)} className="text-muted text-xl leading-none" aria-label="Remove set">
                  ×
                </button>
              </div>
              <div className="grid grid-cols-[1.5rem_1fr_1fr] gap-2 text-xs text-muted px-1">
                <span />
                <span>Weight (kg)</span>
                <span>Reps</span>
              </div>
              {sides.map((side) => (
                <div key={side} className="grid grid-cols-[1.5rem_1fr_1fr] gap-2 items-center">
                  <span className="text-muted text-sm font-medium text-center">
                    {side === "left" ? "L" : "R"}
                  </span>
                  {cell(i, side, "weight", "–")}
                  {cell(i, side, "reps", "–")}
                </div>
              ))}
            </div>
          ) : (
            <div key={i} className="grid grid-cols-[2rem_1fr_1fr_1.5rem] gap-2 items-center">
              <span className="text-muted text-center">{i + 1}</span>
              {cell(i, "both", "weight", "–")}
              {cell(i, "both", "reps", "–")}
              <button onClick={() => removeSet(i)} className="text-muted text-xl" aria-label="Remove set">
                ×
              </button>
            </div>
          )
        )}

        <div className="flex items-center justify-between">
          <button onClick={addSet} className="text-accent text-sm">
            + Add set
          </button>
          <span className="text-[11px] text-muted">Leave weight blank for bodyweight</span>
        </div>
      </div>

      <textarea
        className="input min-h-16 text-sm"
        value={note}
        onChange={(e) => setNote(e.target.value)}
        placeholder="Notes for this cycle (e.g. felt strong, go heavier next time)…"
      />

      {errors.length > 0 && (
        <ul className="text-red-400 text-sm space-y-1" role="alert">
          {errors.map((m, i) => (
            <li key={i}>• {m}</li>
          ))}
        </ul>
      )}

      <button onClick={save} disabled={status === "saving"} className="btn-secondary w-full">
        {status === "saving" ? "Saving…" : status === "error" ? "Retry save" : "Save"}
      </button>
    </div>
  );
}
