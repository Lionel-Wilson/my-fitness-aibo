"use client";

import { useCallback, useEffect, useState } from "react";
import { api } from "./api";

interface State<T> {
  data: T | null;
  loading: boolean;
  error: string | null;
  refetch: () => void;
}

// Simple GET hook with loading/error state and a refetch trigger.
export function useFetch<T>(path: string | null): State<T> {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [tick, setTick] = useState(0);

  const refetch = useCallback(() => setTick((t) => t + 1), []);

  useEffect(() => {
    if (!path) {
      setLoading(false);
      return;
    }
    let active = true;
    setLoading(true);
    api
      .get<T>(path)
      .then((d) => active && setData(d))
      .catch((e) => active && setError(e.message))
      .finally(() => active && setLoading(false));
    return () => {
      active = false;
    };
  }, [path, tick]);

  return { data, loading, error, refetch };
}
