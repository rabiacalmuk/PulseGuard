const BASE_URL = "http://localhost:8080";

export interface Host {
  host_id: string;
  status: "ok" | "warning" | "error" | "unknown";
  last_seen: string;
}

export interface Metric {
  name: string;
  value: number;
  unit: string;
}

export interface Event {
  event_id: string;
  host_id: string;
  check_type: string;
  level: "INFO" | "WARNING" | "ERROR";
  message: string;
  metric: Metric;
  timestamp: string;
  schema_version: number;
}

export async function getHosts(): Promise<Host[]> {
  const res = await fetch(`${BASE_URL}/api/v1/hosts`);
  if (!res.ok) {
    throw new Error(`hosts alinamadi: ${res.status}`);
  }
  const data = await res.json();
  return data.hosts;
}

export async function getEvents(hostId?: string): Promise<Event[]> {
  const url = hostId
    ? `${BASE_URL}/api/v1/events?host_id=${encodeURIComponent(hostId)}`
    : `${BASE_URL}/api/v1/events`;
  const res = await fetch(url);
  if (!res.ok) {
    throw new Error(`events alinamadi: ${res.status}`);
  }
  const data = await res.json();
  return data.events;
}