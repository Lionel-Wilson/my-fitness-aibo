"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { Protected } from "@/components/Protected";
import { PageHeader } from "@/components/PageHeader";
import { LogExerciseCard } from "@/components/LogExerciseCard";
import { useFetch } from "@/lib/useFetch";
import { api } from "@/lib/api";
import type { Workout, Exercise, Cycle } from "@/lib/types";

export default function LogPage() {
  return (
    <Protected>
      <LogInner />
    </Protected>
  );
}

function LogInner() {
  const { workoutId } = useParams<{ workoutId: string }>();
  const { data: workout } = useFetch<Workout>(`/workouts/${workoutId}`);
  const { data: exercises } = useFetch<Exercise[]>(`/workouts/${workoutId}/exercises`);
  const { data: cycles, refetch: refetchCycles } = useFetch<Cycle[]>(
    workout ? `/plans/${workout.planId}/cycles` : null
  );

  const [cycleId, setCycleId] = useState<string>("");

  // Default to the most recent cycle once cycles load.
  useEffect(() => {
    if (cycles && cycles.length > 0 && !cycleId) setCycleId(cycles[0].id);
  }, [cycles, cycleId]);

  const label = (c: Cycle) => c.label || `Cycle ${c.cycleNumber}`;

  const startCycle = async () => {
    const next = (cycles?.[0]?.cycleNumber ?? 0) + 1;
    const created = await api.post<Cycle>(`/plans/${workout!.planId}/cycles`, { label: `C${next}` });
    await refetchCycles();
    setCycleId(created.id);
  };

  return (
    <>
      <PageHeader title="Log session" subtitle={workout?.name} back />

      <main className="flex-1 px-4 py-4 space-y-4">
        {/* No cycles yet: block logging with a clear call to action. */}
        {cycles && cycles.length === 0 ? (
          <div className="card text-center py-8 space-y-3">
            <p className="text-3xl">🔄</p>
            <p className="font-medium">Start a cycle to begin logging</p>
            <p className="text-sm text-muted">
              A cycle is one round through your whole plan. Your sets are recorded against
              it, so you need to open one before you can log.
            </p>
            <button onClick={startCycle} className="btn-primary w-full">
              Start cycle 1
            </button>
          </div>
        ) : (
          <>
            <div className="card flex items-center gap-3">
              <div className="flex-1">
                <label className="label">Logging into</label>
                <select
                  className="input"
                  value={cycleId}
                  onChange={(e) => setCycleId(e.target.value)}
                  disabled={!cycles || cycles.length === 0}
                >
                  {cycles?.map((c) => (
                    <option key={c.id} value={c.id}>
                      {label(c)}
                    </option>
                  ))}
                </select>
              </div>
              <button onClick={startCycle} className="btn-secondary text-sm shrink-0 self-end">
                + New
              </button>
            </div>

            {cycleId &&
              exercises?.map((e) => (
                // Key by cycle so inputs reset/prefill when the cycle changes.
                <LogExerciseCard key={`${e.id}-${cycleId}`} exercise={e} cycleId={cycleId} />
              ))}

            {exercises && exercises.length === 0 && (
              <p className="text-muted text-sm card">This workout has no exercises yet.</p>
            )}
          </>
        )}
      </main>
    </>
  );
}
