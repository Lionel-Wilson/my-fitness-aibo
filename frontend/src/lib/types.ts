// Shapes mirror the Go API JSON responses.

export interface User {
  id: string;
  email: string;
  createdAt: string;
}

export interface Plan {
  id: string;
  name: string;
  quality: string;
  description: string;
  cycleLabel: string;
  periodStart: string | null;
  periodEnd: string | null;
  isActive: boolean;
  createdAt: string;
}

export interface Workout {
  id: string;
  planId: string;
  name: string;
  dayLabel: string;
  orderIndex: number;
  durationMin: number | null;
  notes: string;
  createdAt: string;
}

export interface Exercise {
  id: string;
  workoutId: string;
  name: string;
  orderIndex: number;
  targetSets: number | null;
  repLow: number | null;
  repHigh: number | null;
  rpeLow: number | null;
  rpeHigh: number | null;
  restSeconds: number | null;
  instructions: string;
  orGroup: string;
  isOptional: boolean;
  isUnilateral: boolean;
  createdAt: string;
}

export type Side = "both" | "left" | "right";

export interface Cycle {
  id: string;
  planId: string;
  cycleNumber: number;
  label: string;
  startedAt: string;
  completedAt: string | null;
  notes: string;
}

export interface SetLog {
  id?: string;
  setIndex: number;
  side: Side;
  weightKg: number | null;
  reps: number | null;
  rpe: number | null;
  isDropSet: boolean;
}

export interface ExerciseLog {
  id: string;
  exerciseId: string;
  cycleId: string;
  note: string;
  workingWeightKg: number | null;
  createdAt: string;
  sets: SetLog[];
}

export interface ProgressPoint {
  cycleId: string;
  cycleNumber: number;
  label: string;
  side: Side;
  topWeightKg: number;
  volumeKg: number;
  bestE1rm: number;
  totalReps: number;
  startedAt: string;
}
