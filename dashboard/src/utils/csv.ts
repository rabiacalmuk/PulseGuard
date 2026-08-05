import type { Event } from "../api/client";

export function eventsToCSV(events: Event[]): string {  //events listesini csv formatında bir metne çeviriyoruz
  const header = "timestamp,host_id,check_type,level,message,metric_name,metric_value,metric_unit";
  const rows = events.map((e) =>
    [
      e.timestamp,
      e.host_id,
      e.check_type,
      e.level,
      escapeCSVField(e.message),
      e.metric?.name ?? "",
      e.metric?.value ?? "",
      e.metric?.unit ?? "",
    ].join(",")
  );
  return [header, ...rows].join("\n");
}

function escapeCSVField(field: string): string {
  if (field.includes(",") || field.includes('"') || field.includes("\n")) {
    return `"${field.replace(/"/g, '""')}"`;
  }
  return field;
}

export function downloadCSV(filename: string, csvContent: string) {   //üretilen csv matnini gerçek bir dosya olarak indirmemizi sağlayan fonksiyon
  const blob = new Blob([csvContent], { type: "text/csv;charset=utf-8;" });
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = filename;
  link.click();
  URL.revokeObjectURL(url);
}