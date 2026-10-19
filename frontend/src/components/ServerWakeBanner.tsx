interface ServerWakeBannerProps {
  visible: boolean;
}

export default function ServerWakeBanner({
  visible,
}: ServerWakeBannerProps) {
  if (!visible) return null;

  return (
    <p
      className="border-l-4 border-primary bg-[#f3efff] px-4 py-3 text-sm font-medium text-black"
      role="status"
      aria-live="polite"
      data-testid="status-server-wake"
    >
      Waking up the server, this can take about a minute.
    </p>
  );
}
