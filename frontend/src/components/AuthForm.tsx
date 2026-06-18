"use client";

import { useState } from "react";
import Link from "next/link";
import { useAuth } from "@/lib/auth";

// Shared email/password form for both login and signup.
export function AuthForm({ mode }: { mode: "login" | "signup" }) {
  const { login, signup } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setBusy(true);
    try {
      if (mode === "login") await login(email, password);
      else await signup(email, password);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Something went wrong");
      setBusy(false);
    }
  };

  return (
    <div className="flex-1 flex flex-col justify-center px-6 py-10">
      <div className="mb-8 text-center">
        <div className="text-4xl mb-2">🏋️</div>
        <h1 className="text-2xl font-bold">Fitness Aibo</h1>
        <p className="text-muted text-sm mt-1">
          {mode === "login" ? "Welcome back" : "Create your account"}
        </p>
      </div>

      <form onSubmit={submit} className="space-y-4">
        <div>
          <label className="label">Email</label>
          <input
            className="input"
            type="email"
            autoComplete="email"
            inputMode="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div>
          <label className="label">Password</label>
          <input
            className="input"
            type="password"
            autoComplete={mode === "login" ? "current-password" : "new-password"}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            minLength={8}
            required
          />
        </div>

        {error && <p className="text-red-400 text-sm">{error}</p>}

        <button className="btn-primary w-full" disabled={busy}>
          {busy ? "…" : mode === "login" ? "Log in" : "Sign up"}
        </button>
      </form>

      <p className="text-center text-sm text-muted mt-6">
        {mode === "login" ? (
          <>
            No account?{" "}
            <Link href="/signup" className="text-accent">
              Sign up
            </Link>
          </>
        ) : (
          <>
            Already have an account?{" "}
            <Link href="/login" className="text-accent">
              Log in
            </Link>
          </>
        )}
      </p>
    </div>
  );
}
