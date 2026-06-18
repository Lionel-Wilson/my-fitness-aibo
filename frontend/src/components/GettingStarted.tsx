"use client";

const STEPS = [
  { title: "Create a plan", body: "Your training block, named after its goal (e.g. “Building muscle”)." },
  { title: "Add workouts", body: "The training days in the plan (e.g. “Monday — Legs”)." },
  { title: "Add exercises", body: "Put movements in each workout with target sets, reps and RPE." },
  { title: "Start a cycle", body: "One full round through the whole plan (C1, C2…). You log against it." },
  { title: "Log a session", body: "Open a workout → “Log this workout” → enter your weight & reps." },
  { title: "Track progress", body: "The Progress tab charts each exercise across your cycles." },
];

// Ordered onboarding guide. Shown expanded for first-time users (no plans) and
// available as a collapsible refresher otherwise.
export function GettingStarted({ defaultOpen = false }: { defaultOpen?: boolean }) {
  return (
    <details open={defaultOpen} className="card group">
      <summary className="flex items-center justify-between cursor-pointer list-none font-medium">
        <span>How it works</span>
        <span className="text-muted text-sm group-open:rotate-180 transition-transform">⌄</span>
      </summary>
      <ol className="mt-4 space-y-3">
        {STEPS.map((s, i) => (
          <li key={i} className="flex gap-3">
            <span className="shrink-0 h-6 w-6 rounded-full bg-accent/15 text-accent text-sm font-semibold flex items-center justify-center">
              {i + 1}
            </span>
            <div className="min-w-0">
              <p className="font-medium leading-6">{s.title}</p>
              <p className="text-sm text-muted">{s.body}</p>
            </div>
          </li>
        ))}
      </ol>
    </details>
  );
}
