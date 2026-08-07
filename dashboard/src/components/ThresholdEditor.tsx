import { useEffect, useState } from "react";
import { getThresholds, setThresholds, type ThresholdValue } from "../api/client";

interface Props {
  hostId: string;
}

const CHECK_TYPES = ["cpu", "ram", "disk"];

export function ThresholdEditor({ hostId }: Props) {
  const [values, setValues] = useState<Record<string, ThresholdValue>>({});
  const [status, setStatus] = useState<"idle" | "saving" | "saved" | "error">("idle");
  const [errorMsg, setErrorMsg] = useState("");

  useEffect(() => {
    getThresholds(hostId)
      .then((data) => {
        const merged: Record<string, ThresholdValue> = {};
        for (const type of CHECK_TYPES) {
          merged[type] = data.thresholds[type] ?? { warning: 0, error: 0 };
        }
        setValues(merged);
      })
      .catch(() => {
        const empty: Record<string, ThresholdValue> = {};
        for (const type of CHECK_TYPES) {
          empty[type] = { warning: 0, error: 0 };
        }
        setValues(empty);
      });
  }, [hostId]);

  function updateField(type: string, field: "warning" | "error", value: number) {
    setValues((prev) => ({
      ...prev,
      [type]: { ...prev[type], [field]: value },
    }));
  }

  async function handleSave() {
     for (const type of CHECK_TYPES) {
    const v = values[type];
    if (v.warning < 0 || v.warning > 100 || v.error < 0 || v.error > 100) {
      setErrorMsg(`${type}: değerler 0-100 arasında olmalı`);
      setStatus("error");
      return;
    }
    if (v.error <= v.warning) {
      setErrorMsg(`${type}: hata eşiği, uyarı eşiğinden büyük olmalı`);
      setStatus("error");
      return;
    }
  }
    setStatus("saving");
    try {
      await setThresholds(hostId, values);
      setStatus("saved");
      setTimeout(() => setStatus("idle"), 2000);
    } catch (err) {
      setErrorMsg(err instanceof Error ? err.message : "bilinmeyen hata");
      setStatus("error");
    }
  }

  return (
    <div className="threshold-editor">
      <h3>Eşik Ayarları</h3>
      <table>
        <thead>
          <tr>
            <th>Check</th>
            <th>Uyarı</th>
            <th>Hata</th>
          </tr>
        </thead>
        <tbody>
          {CHECK_TYPES.map((type) => (
            <tr key={type}>
              <td>{type}</td>
              <td>
                <input
                  type="number"
                  min={0}
                  max={100}
                  value={values[type]?.warning ?? 0}
                  onChange={(e) => updateField(type, "warning", Number(e.target.value))}
                />
              </td>
              <td>
                <input
                  type="number"
                  min={0}
                  max={100}
                  value={values[type]?.error ?? 0}
                  onChange={(e) => updateField(type, "error", Number(e.target.value))}
                />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <button onClick={handleSave} disabled={status === "saving"}>
        {status === "saving" ? "Kaydediliyor..." : "Kaydet"}
      </button>
      {status === "saved" && <span className="status-ok"> Kaydedildi ✓</span>}
      {status === "error" && <span className="status-error"> Hata: {errorMsg}</span>}
    </div>
  );
}