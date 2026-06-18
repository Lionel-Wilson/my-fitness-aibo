"use client";

import { createContext, useContext, useEffect, useState } from "react";

interface Settings {
  // Auto-fill a set's weight into the other sets, and copy the previous set
  // when adding a new one.
  autofillSets: boolean;
}

const DEFAULTS: Settings = {
  autofillSets: true,
};

const STORAGE_KEY = "aibo_settings";

interface SettingsContextValue extends Settings {
  setSetting: <K extends keyof Settings>(key: K, value: Settings[K]) => void;
}

const SettingsContext = createContext<SettingsContextValue | undefined>(undefined);

// Lightweight per-device settings persisted to localStorage.
export function SettingsProvider({ children }: { children: React.ReactNode }) {
  const [settings, setSettings] = useState<Settings>(DEFAULTS);

  useEffect(() => {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      if (raw) setSettings({ ...DEFAULTS, ...JSON.parse(raw) });
    } catch {
      /* ignore malformed settings */
    }
  }, []);

  const setSetting: SettingsContextValue["setSetting"] = (key, value) => {
    setSettings((prev) => {
      const next = { ...prev, [key]: value };
      try {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(next));
      } catch {
        /* ignore quota errors */
      }
      return next;
    });
  };

  return (
    <SettingsContext.Provider value={{ ...settings, setSetting }}>
      {children}
    </SettingsContext.Provider>
  );
}

export function useSettings() {
  const ctx = useContext(SettingsContext);
  if (!ctx) throw new Error("useSettings must be used within SettingsProvider");
  return ctx;
}
