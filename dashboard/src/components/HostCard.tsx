import type { Host } from "../api/client";

const statusColors: Record<string, string> = {
  ok: "#22c55e",
  warning: "#eab308",
  error: "#ef4444",
  unknown: "#9ca3af",
};

interface Props {
  host: Host;
  onClick?: () => void;

}

export function HostCard({ host, onClick }: Props) {
  const color = statusColors[host.status] ?? statusColors.unknown;

  return (
    
    <div
      onClick={onClick} 
      style={{ border: `2px solid ${color}`, borderRadius: 8, padding: 16, minWidth: 200 , cursor: onClick ? "pointer" : "default" }}>
       <div style={{ display: "flex", alignItems: "center", gap: 8 }}>
         <span style={{ width: 12, height: 12, borderRadius: "50%", backgroundColor: color, display: "inline-block" }} />
         <strong>{host.host_id}</strong>
     </div>
      <p style={{ margin: "8px 0 0", color: "#6b7280", fontSize: 14 }}>
        Son gorulme: {new Date(host.last_seen).toLocaleString()}
      </p>
    </div>
  );
}