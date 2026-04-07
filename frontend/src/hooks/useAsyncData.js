import { useCallback, useEffect, useState } from "react";

export function useAsyncData(loader, dependencies = [], options = {}) {
  const { immediate = true } = options;
  const [data, setData] = useState(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(immediate);

  const execute = useCallback(async () => {
    setLoading(true);
    setError("");

    try {
      const result = await loader();
      setData(result);
      return result;
    } catch (err) {
      setError(err.message || "Something went wrong.");
      throw err;
    } finally {
      setLoading(false);
    }
  }, dependencies);

  useEffect(() => {
    if (!immediate) return;
    execute().catch(() => null);
  }, [execute, immediate]);

  return { data, error, loading, execute, setData };
}
