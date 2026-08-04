import { useEffect, useState } from "react";
import { getHosts, type Host } from "../api/client";
import { HostCard } from "../components/HostCard";

export function FleetView() {
  const [hosts, setHosts] = useState<Host[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function load() {
      try {
        const data = await getHosts();
        setHosts(data);
        setError(null);
      } catch (err) {
        setError(err instanceof Error ? err.message : "bilinmeyen hata");
      }
    }

    load();
    const interval = setInterval(load, 5000);
    return () => clearInterval(interval);
  }, []);

  if (error) {
    return <p style={{ color: "#ef4444" }}>Hata: {error}</p>;
  }

  if (hosts.length === 0) {
    return <p>Henuz hic host verisi yok.</p>;
  }

  return (
    <div style={{ display: "flex", gap: 16, flexWrap: "wrap" }}>
      {hosts.map((h) => (
        <HostCard key={h.host_id} host={h} />
      ))}
    </div>
  );
}