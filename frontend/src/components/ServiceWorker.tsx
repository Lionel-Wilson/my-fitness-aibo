"use client";

import { useEffect } from "react";

// Registers the service worker so the app is installable and works offline-ish.
export function ServiceWorker() {
  useEffect(() => {
    if ("serviceWorker" in navigator) {
      navigator.serviceWorker.register("/sw.js").catch(() => {
        /* registration is best-effort */
      });
    }
  }, []);
  return null;
}
