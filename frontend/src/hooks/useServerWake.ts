import { useCallback, useEffect, useRef, useState } from "react";

export default function useServerWake() {
  const [isWaking, setIsWaking] = useState(false);
  const timerRef = useRef<number | null>(null);

  const runWithWake = useCallback(async <T,>(request: Promise<T>) => {
    setIsWaking(false);
    timerRef.current = window.setTimeout(() => setIsWaking(true), 3000);
    try {
      return await request;
    } finally {
      if (timerRef.current !== null) {
        window.clearTimeout(timerRef.current);
        timerRef.current = null;
      }
      setIsWaking(false);
    }
  }, []);

  useEffect(
    () => () => {
      if (timerRef.current !== null) window.clearTimeout(timerRef.current);
    },
    [],
  );

  return { isWaking, runWithWake };
}
