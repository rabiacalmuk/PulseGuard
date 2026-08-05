import { useEffect, useState } from "react";
import { getEvents, type Event } from "../api/client";
import { eventsToCSV, downloadCSV } from "../utils/csv";
interface Props {
  hostId: string;
  onBack: () => void;
}

const levelColors: Record<string, string> = {
  INFO: "#22c55e",
  WARNING: "#eab308",
  ERROR: "#ef4444",
};

export function HostDetail({ hostId, onBack }: Props) {
  const [events, setEvents] = useState<Event[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function load() {
      try {
        const data = await getEvents(hostId);
        setEvents(data);
        setError(null);
      } catch (err) {
        setError(err instanceof Error ? err.message : "bilinmeyen hata");
      }
    }

    load();
    const interval = setInterval(load, 5000);
    return () => clearInterval(interval);
  }, [hostId]);

  return (
    <div>
      <button onClick={onBack} style={{ marginBottom: 16 }}>
        ← Filo görünümüne dön
      </button>
      <h2>{hostId}</h2>
      <button
        onClick={() => downloadCSV(`${hostId}-rapor.csv`, eventsToCSV(events))}
        style={{ marginBottom: 16 }}
        disabled={events.length === 0}
    >
         CSV olarak indir
      </button>
      {error && <p style={{ color: "#ef4444" }}>Hata: {error}</p>}
      {!error && events.length === 0 && <p>Henuz hic event yok.</p>}

      <table style={{ borderCollapse: "collapse", width: "100%" }}>
        <thead>
          <tr style={{ textAlign: "left", borderBottom: "1px solid #e5e7eb" }}>
            <th style={{ padding: 8 }}>Zaman</th>
            <th style={{ padding: 8 }}>Check</th>
            <th style={{ padding: 8 }}>Seviye</th>
            <th style={{ padding: 8 }}>Mesaj</th>
          </tr>
        </thead>
        <tbody>
          {events.map((e) => (
            <tr key={e.event_id} style={{ borderBottom: "1px solid #f3f4f6" }}>
              <td style={{ padding: 8 }}>{new Date(e.timestamp).toLocaleString()}</td>
              <td style={{ padding: 8 }}>{e.check_type}</td>
              <td style={{ padding: 8, color: levelColors[e.level] }}>{e.level}</td>
              <td style={{ padding: 8 }}>{e.message}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}