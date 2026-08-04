import { useState } from "react";
import { FleetView } from "./pages/FleetView";
import { HostDetail } from "./pages/HostDetail";

function App() {
  const [selectedHost, setSelectedHost] = useState<string | null>(null);

  return (
    <div style={{ padding: 32, fontFamily: "sans-serif" }}>
      <h1>PulseGuard — Filo Görünümü</h1>
      {selectedHost ? (
        <HostDetail hostId={selectedHost} onBack={() => setSelectedHost(null)} />
      ) : (
        <FleetView onSelectHost={setSelectedHost} />
      )}
    </div>
  );
}

export default App;