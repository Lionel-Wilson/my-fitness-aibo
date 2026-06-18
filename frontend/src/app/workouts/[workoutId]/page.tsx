"use client";

import { useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { Protected } from "@/components/Protected";
import { PageHeader } from "@/components/PageHeader";
import { ExerciseForm } from "@/components/ExerciseForm";
import { useFetch } from "@/lib/useFetch";
import { prescription, restLabel } from "@/lib/format";
import type { Workout, Exercise } from "@/lib/types";

export default function WorkoutPage() {
  return (
    <Protected>
      <WorkoutInner />
    </Protected>
  );
}

function WorkoutInner() {
  const { workoutId } = useParams<{ workoutId: string }>();
  const { data: workout } = useFetch<Workout>(`/workouts/${workoutId}`);
  const { data: exercises, refetch } = useFetch<Exercise[]>(`/workouts/${workoutId}/exercises`);
  const [open, setOpen] = useState(false);

  return (
    <>
      <PageHeader title={workout?.name || "Workout"} subtitle={workout?.dayLabel} back />

      <main className="flex-1 px-4 py-4 space-y-4">
        {exercises && exercises.length > 0 && (
          <Link href={`/log/${workoutId}`} className="btn-primary w-full">
            ▶ Log this workout
          </Link>
        )}

        <div className="space-y-2">
          {exercises?.length === 0 && (
            <div className="card text-sm text-muted">
              <p className="text-white font-medium mb-1">Next: add exercises</p>
              <p>Add the movements for this workout with their target sets, reps and RPE.
                Once a plan has a cycle, you’ll log weights against them here.</p>
            </div>
          )}
          {exercises?.map((e) => (
            <Link key={e.id} href={`/exercises/${e.id}`} className="card block active:bg-surface2">
              <div className="flex items-center justify-between gap-2">
                <p className="font-medium">{e.name}</p>
                {e.isOptional && <span className="text-xs text-muted shrink-0">optional</span>}
              </div>
              <p className="text-sm text-muted">
                {[prescription(e), restLabel(e.restSeconds)].filter(Boolean).join(" · ")}
              </p>
            </Link>
          ))}
        </div>

        <button onClick={() => setOpen(true)} className="btn-secondary w-full">
          + Add exercise
        </button>
      </main>

      <ExerciseForm
        workoutId={workoutId}
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
