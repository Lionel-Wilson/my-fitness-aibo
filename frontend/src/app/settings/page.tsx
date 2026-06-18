"use client";

import { Protected } from "@/components/Protected";
import { PageHeader } from "@/components/PageHeader";
import { BottomNav } from "@/components/BottomNav";
import { useSettings } from "@/lib/settings";
import { useAuth } from "@/lib/auth";

export default function SettingsPage() {
  return (
    <Protected>
      <SettingsInner />
    </Protected>
  );
}

function SettingsInner() {
  const { autofillSets, setSetting } = useSettings();
  const { user, logout } = useAuth();

  return (
    <>
      <PageHeader title="Settings" />

      <main className="flex-1 px-4 py-4 space-y-4">
        <section>
          <h2 className="text-sm uppercase tracking-wide text-muted mb-2">Logging</h2>
          <div className="card">
            <Toggle
              label="Auto-fill sets"
              description="When you enter a set's weight, copy it to the other sets. Adding a set copies the previous set's values."
              checked={autofillSets}
              onChange={(v) => setSetting("autofillSets", v)}
            />
          </div>
        </section>

        <section>
          <h2 className="text-sm uppercase tracking-wide text-muted mb-2">Account</h2>
          <div className="card space-y-3">
            <p className="text-sm text-muted">
              Signed in as <span className="text-white">{user?.email}</span>
            </p>
            <button onClick={logout} className="btn-secondary w-full">
              Log out
            </button>
          </div>
        </section>
      </main>

      <BottomNav />
    </>
  );
}

function Toggle({
  label,
  description,
  checked,
  onChange,
}: {
  label: string;
  description?: string;
  checked: boolean;
  onChange: (v: boolean) => void;
}) {
  return (
    <label className="flex items-start justify-between gap-4 cursor-pointer">
      <span className="min-w-0">
        <span className="font-medium">{label}</span>
        {description && <span className="block text-sm text-muted mt-0.5">{description}</span>}
      </span>
      <button
        type="button"
        role="switch"
        aria-checked={checked}
        onClick={() => onChange(!checked)}
        className={`relative h-7 w-12 shrink-0 rounded-full transition-colors ${
          checked ? "bg-accent" : "bg-surface2 border border-border"
        }`}
      >
        <span
          className={`absolute top-0.5 h-6 w-6 rounded-full bg-white transition-transform ${
            checked ? "translate-x-5" : "translate-x-0.5"
          }`}
        />
      </button>
    </label>
  );
}
