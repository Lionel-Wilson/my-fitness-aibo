"use client";

import { useState } from "react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import type { ProgressPoint } from "@/lib/types";

type Metric = { key: keyof ProgressPoint; label: string };
const METRICS: Metric[] = [
  { key: "bestE1rm", label: "Est. 1RM" },
  { key: "topWeightKg", label: "Top weight" },
  { key: "volumeKg", label: "Volume" },
];

// Line chart of an exercise's progress across cycles. Bilateral exercises show a
// single line; unilateral exercises show separate Left and Right lines.
export function ProgressChart({ data }: { data: ProgressPoint[] }) {
  const [metric, setMetric] = useState<Metric>(METRICS[0]);

  if (data.length === 0) {
    return <p className="text-muted text-sm card">No logged sets for this exercise yet.</p>;
  }

  const unilateral = data.some((d) => d.side === "left" || d.side === "right");

  // Merge points by cycle, with a column per side.
  const byCycle = new Map<string, { name: string; order: number; left?: number; right?: number; both?: number }>();
  data.forEach((p) => {
    const name = p.label || `C${p.cycleNumber}`;
    const entry = byCycle.get(name) ?? { name, order: p.cycleNumber };
    entry[p.side] = Math.round((p[metric.key] as number) * 10) / 10;
    byCycle.set(name, entry);
  });
  const chartData = [...byCycle.values()].sort((a, b) => a.order - b.order);

  return (
    <div className="card space-y-3">
      <div className="flex gap-2">
        {METRICS.map((m) => (
          <button
            key={m.key}
            onClick={() => setMetric(m)}
            className={`text-xs px-3 py-1.5 rounded-full border ${
              metric.key === m.key ? "bg-accent text-black border-accent" : "border-border text-muted"
            }`}
          >
            {m.label}
          </button>
        ))}
      </div>

      <div className="h-64 -ml-2">
        <ResponsiveContainer width="100%" height="100%">
          <LineChart data={chartData} margin={{ top: 8, right: 12, bottom: 0, left: 0 }}>
            <CartesianGrid stroke="#2a2a36" strokeDasharray="3 3" />
            <XAxis dataKey="name" stroke="#8a8a99" fontSize={12} />
            <YAxis stroke="#8a8a99" fontSize={12} width={40} domain={["auto", "auto"]} />
            <Tooltip
              contentStyle={{
                background: "#16161d",
                border: "1px solid #2a2a36",
                borderRadius: 12,
                color: "#ededed",
              }}
              labelStyle={{ color: "#8a8a99" }}
            />
            {/* Children must be direct descendants — Recharts does not look through fragments. */}
            {unilateral && <Legend wrapperStyle={{ fontSize: 12 }} />}
            {unilateral && (
              <Line type="monotone" dataKey="left" name="Left" stroke="#f5a623" strokeWidth={2.5} dot={{ r: 4, fill: "#f5a623" }} activeDot={{ r: 6 }} connectNulls />
            )}
            {unilateral && (
              <Line type="monotone" dataKey="right" name="Right" stroke="#3b82f6" strokeWidth={2.5} dot={{ r: 4, fill: "#3b82f6" }} activeDot={{ r: 6 }} connectNulls />
            )}
            {!unilateral && (
              <Line type="monotone" dataKey="both" name={metric.label} stroke="#f5a623" strokeWidth={2.5} dot={{ r: 4, fill: "#f5a623" }} activeDot={{ r: 6 }} />
            )}
          </LineChart>
        </ResponsiveContainer>
      </div>

      <p className="text-xs text-muted text-center">
        {metric.label} per cycle {metric.key !== "volumeKg" ? "(kg)" : "(kg total)"}
        {unilateral ? " · left vs right" : ""}
      </p>
    </div>
  );
}
