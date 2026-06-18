"use client";

import { Suspense, useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import { Protected } from "@/components/Protected";
import { PageHeader } from "@/components/PageHeader";
import { BottomNav } from "@/components/BottomNav";
import { ProgressChart } from "@/components/ProgressChart";
import { useFetch } from "@/lib/useFetch";
import { api } from "@/lib/api";
import type { Plan, Workout, Exercise, ProgressPoint } from "@/lib/types";

export default function DashboardPage() {
  return (
    <Protected>
      <Suspense fallback={<div className="flex-1" />}>
        <DashboardInner />
      </Suspense>
    </Protected>
  );
}

function DashboardInner() {
  const search = useSearchParams();
  const presetExercise = search.get("exercise");

  const { data: plans } = useFetch<Plan[]>("/plans");
  const [planId, setPlanId] = useState("");
  const [workoutId, setWorkoutId] = useState("");
  const [exerciseId, setExerciseId] = useState("");
  const [presetApplied, setPresetApplied] = useState(false);

  const { data: workouts } = useFetch<Workout[]>(planId ? `/plans/${planId}/workouts` : null);
  const { data: exercises } = useFetch<Exercise[]>(workoutId ? `/workouts/${workoutId}/exercises` : null);
  const { data: progress } = useFetch<ProgressPoint[]>(
    exerciseId ? `/exercises/${exerciseId}/progress` : null
  );

  // Resolve a preset ?exercise= into the plan/workout/exercise selectors.
  useEffect(() => {
    if (!presetExercise || presetApplied) return;
    (async () => {
      try {
        const ex = await api.get<Exercise>(`/exercises/${presetExercise}`);
        const wk = await api.get<Workout>(`/workouts/${ex.workoutId}`);
        setPlanId(wk.planId);
        setWorkoutId(wk.id);
        setExerciseId(ex.id);
      } catch {
        /* fall back to manual selection */
      } finally {
        setPresetApplied(true);
      }
    })();
  }, [presetExercise, presetApplied]);

  // Default selections when no preset.
  useEffect(() => {
    if (!presetExercise && plans && plans.length && !planId) setPlanId(plans[0].id);
  }, [plans, planId, presetExercise]);

  const selectClass = "input";

  return (
    <>
      <PageHeader title="Progress" />

      <main className="flex-1 px-4 py-4 space-y-3">
        <div className="card space-y-3">
          <div>
            <label className="label">Plan</label>
            <select
              className={selectClass}
              value={planId}
              onChange={(e) => {
                setPlanId(e.target.value);
                setWorkoutId("");
                setExerciseId("");
              }}
            >
              <option value="">Select a plan…</option>
              {plans?.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name}
                </option>
              ))}
            </select>
          </div>

          {planId && (
            <div>
              <label className="label">Workout</label>
              <select
                className={selectClass}
                value={workoutId}
                onChange={(e) => {
                  setWorkoutId(e.target.value);
                  setExerciseId("");
                }}
              >
                <option value="">Select a workout…</option>
                {workouts?.map((w) => (
                  <option key={w.id} value={w.id}>
                    {w.name}
                  </option>
                ))}
              </select>
            </div>
          )}

          {workoutId && (
            <div>
              <label className="label">Exercise</label>
              <select
                className={selectClass}
                value={exerciseId}
                onChange={(e) => setExerciseId(e.target.value)}
              >
                <option value="">Select an exercise…</option>
                {exercises?.map((ex) => (
                  <option key={ex.id} value={ex.id}>
                    {ex.name}
                  </option>
                ))}
              </select>
            </div>
          )}
        </div>

        {exerciseId && progress && <ProgressChart data={progress} />}
        {!exerciseId && (
          <p className="text-muted text-sm card text-center py-8">
            Pick an exercise to see how your strength has changed over time.
          </p>
        )}
      </main>

      <BottomNav />
    </>
  );
}
