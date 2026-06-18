import type { Exercise } from "./types";

// "6–8", "8", or "" depending on which rep targets are set.
export function repRange(e: Pick<Exercise, "repLow" | "repHigh">): string {
  if (e.repLow != null && e.repHigh != null) return `${e.repLow}–${e.repHigh}`;
  if (e.repLow != null) return `${e.repLow}`;
  if (e.repHigh != null) return `${e.repHigh}`;
  return "";
}

export function rpeRange(e: Pick<Exercise, "rpeLow" | "rpeHigh">): string {
  const f = (n: number) => (Number.isInteger(n) ? `${n}` : n.toFixed(1));
  if (e.rpeLow != null && e.rpeHigh != null && e.rpeLow !== e.rpeHigh)
    return `${f(e.rpeLow)}–${f(e.rpeHigh)}`;
  if (e.rpeLow != null) return f(e.rpeLow);
  if (e.rpeHigh != null) return f(e.rpeHigh);
  return "";
}

// "2×6–8 @ RPE 7" style prescription summary.
export function prescription(e: Exercise): string {
  const parts: string[] = [];
  const reps = repRange(e);
  if (e.targetSets != null && reps) parts.push(`${e.targetSets}×${reps}`);
  else if (e.targetSets != null) parts.push(`${e.targetSets} sets`);
  else if (reps) parts.push(`${reps} reps`);
  const rpe = rpeRange(e);
  if (rpe) parts.push(`@ RPE ${rpe}`);
  return parts.join(" ");
}

export function restLabel(seconds: number | null): string {
  if (seconds == null) return "";
  if (seconds < 60) return `${seconds}s rest`;
  const m = seconds / 60;
  return `${Number.isInteger(m) ? m : m.toFixed(1)} min rest`;
}

export function num(n: number): string {
  return Number.isInteger(n) ? `${n}` : n.toFixed(1);
}
