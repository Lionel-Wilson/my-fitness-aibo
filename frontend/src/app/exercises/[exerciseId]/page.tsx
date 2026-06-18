"use client";

import { useState } from "react";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { Protected } from "@/components/Protected";
import { PageHeader } from "@/components/PageHeader";
import { ExerciseForm } from "@/components/ExerciseForm";
import { useFetch } from "@/lib/useFetch";
import { api } from "@/lib/api";
import { prescription, restLabel, num } from "@/lib/format";
import type { Exercise, Workout, Cycle, ExerciseLog } from "@/lib/types";

export default function ExercisePage() {
  return (
    <Protected>
      <ExerciseInner />
    </Protected>
  );
}

function ExerciseInner() {
  const { exerciseId } = useParams<{ exerciseId: string }>();
  const router = useRouter();
  const { data: exercise, refetch } = useFetch<Exercise>(`/exercises/${exerciseId}`);
  const { data: workout } = useFetch<Workout>(exercise ? `/workouts/${exercise.workoutId}` : null);
  const { data: cycles } = useFetch<Cycle[]>(workout ? `/plans/${workout.planId}/cycles` : null);
  const { data: logs } = useFetch<ExerciseLog[]>(`/exercises/${exerciseId}/logs`);
  const [editing, setEditing] = useState(false);

  const remove = async () => {
    if (!confirm("Delete this exercise and its logs?")) return;
    await api.delete(`/exercises/${exerciseId}`);
    router.back();
  };

  const cycleById = new Map((cycles ?? []).map((c) => [c.id, c]));
  // Columns: cycles that have a log, ordered by cycle number.
  const cols = (logs ?? [])
    .map((l) => ({ log: l, cycle: cycleById.get(l.cycleId) }))
    .sort((a, b) => (a.cycle?.cycleNumber ?? 0) - (b.cycle?.cycleNumber ?? 0));
  const maxSets = Math.max(0, ...cols.flatMap((c) => c.log.sets.map((s) => s.setIndex)));
  const unilateral = exercise?.isUnilateral ?? false;

  return (
    <>
      <PageHeader
        title={exercise?.name || "Exercise"}
        subtitle={workout?.name}
        back
        right={
          <button onClick={() => setEditing(true)} className="text-sm text-accent">
            Edit
          </button>
        }
      />

      <main className="flex-1 px-4 py-4 space-y-4">
        {exercise && (
          <div className="card space-y-2">
            <p className="font-medium">{prescription(exercise) || "No target set"}</p>
            {restLabel(exercise.restSeconds) && (
              <p className="text-sm text-muted">{restLabel(exercise.restSeconds)}</p>
            )}
            {exercise.instructions && (
              <p className="text-sm whitespace-pre-wrap text-muted/90 border-t border-border pt-2">
                {exercise.instructions}
              </p>
            )}
          </div>
        )}

        <Link href={`/log/${exercise?.workoutId}`} className="btn-primary w-full">
          ▶ Log a session
        </Link>

        <div>
          <div className="flex items-center justify-between mb-2">
            <h2 className="text-sm uppercase tracking-wide text-muted">History</h2>
            <Link href={`/dashboard?exercise=${exerciseId}`} className="text-accent text-sm">
              View chart →
            </Link>
          </div>

          {cols.length === 0 ? (
            <p className="text-muted text-sm card">No logged sets yet.</p>
          ) : (
            <div className="card overflow-x-auto p-0">
              <table className="w-full text-sm border-collapse">
                <thead>
                  <tr className="text-muted">
                    <th className="text-left px-3 py-2 sticky left-0 bg-surface">Set</th>
                    {cols.map(({ log, cycle }) => (
                      <th key={log.id} className="px-3 py-2 text-center whitespace-nowrap">
                        {cycle?.label || `C${cycle?.cycleNumber ?? "?"}`}
                        {log.workingWeightKg != null && (
                          <span className="block text-[10px] text-accent">
                            {num(log.workingWeightKg)}kg
                          </span>
                        )}
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {Array.from({ length: maxSets }).map((_, i) => {
                    const setNo = i + 1;
                    return (
                      <tr key={i} className="border-t border-border">
                        <td className="px-3 py-2 text-muted sticky left-0 bg-surface align-top">{setNo}</td>
                        {cols.map(({ log }) => {
                          const atSet = log.sets.filter((s) => s.setIndex === setNo);
                          const cellFor = (side: string) => atSet.find((s) => s.side === side);
                          const render = (s?: typeof atSet[number]) =>
                            s ? (
                              <>
                                <span className="font-medium">{s.reps ?? "–"}</span>
                                {s.weightKg != null && (
                                  <span className="text-[10px] text-muted"> · {num(s.weightKg)}kg</span>
                                )}
                              </>
                            ) : (
                              <span className="text-muted">–</span>
                            );
                          return (
                            <td key={log.id} className="px-3 py-2 text-center whitespace-nowrap align-top">
                              {unilateral ? (
                                <div className="space-y-0.5">
                                  <div className="text-xs">
                                    <span className="text-muted mr-1">L</span>
                                    {render(cellFor("left"))}
                                  </div>
                                  <div className="text-xs">
                                    <span className="text-muted mr-1">R</span>
                                    {render(cellFor("right"))}
                                  </div>
                                </div>
                              ) : (
                                render(cellFor("both"))
                              )}
                            </td>
                          );
                        })}
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          )}

          {/* Per-cycle notes */}
          <div className="space-y-2 mt-3">
            {cols
              .filter((c) => c.log.note.trim())
              .map(({ log, cycle }) => (
                <div key={log.id} className="card">
                  <p className="text-xs text-accent mb-1">
                    {cycle?.label || `C${cycle?.cycleNumber}`}
                  </p>
                  <p className="text-sm whitespace-pre-wrap">{log.note}</p>
                </div>
              ))}
          </div>
        </div>

        <button onClick={remove} className="text-sm text-red-400 w-full py-3">
          Delete exercise
        </button>
      </main>

      {exercise && (
        <ExerciseForm
          key={exercise.id}
          exercise={exercise}
          open={editing}
          onClose={() => setEditing(false)}
          onSaved={() => {
            setEditing(false);
            refetch();
          }}
        />
      )}
    </>
  );
}
