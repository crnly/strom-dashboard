// Dieser Vertrag gilt für den aktuellen Demo-Endpunkt, noch nicht für echte Preise.
interface DemoPrices {
  demo: true;
  date: string;
  market: "DE-LU";
  timezone: "Europe/Berlin";
  unit: "EUR/MWh";
  interval: "1 hour";
  prices: number[];
}

function isDemoPrices(value: unknown): value is DemoPrices {
  if (typeof value !== "object" || value === null) return false;
  const data = value as Record<string, unknown>;
  return data.demo === true && data.date === "2026-10-31" &&
    data.market === "DE-LU" && data.timezone === "Europe/Berlin" &&
    data.unit === "EUR/MWh" && data.interval === "1 hour" &&
    Array.isArray(data.prices) && data.prices.length === 24 &&
    data.prices.every((price: unknown) =>
      typeof price === "number" && Number.isFinite(price) && price >= 0 && price <= 120);
}

async function loadPrices(): Promise<void> {
  const bars = document.querySelector<HTMLElement>(".bars");
  const chart = document.querySelector<HTMLElement>(".chart");
  const status = document.querySelector<HTMLElement>("#data-status");
  const json = document.querySelector<HTMLElement>("#prices-json");
  if (!bars || !chart || !status || !json) return;

  try {
    const response = await fetch("/api/prices", { signal: AbortSignal.timeout(5000) });
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    const data: unknown = await response.json();
    // TypeScript-Typen prüfen keine Netzwerkantworten zur Laufzeit.
    if (!isDemoPrices(data)) throw new Error("Unerwartetes Demoformat");

    const fragment = document.createDocumentFragment();
    data.prices.forEach((price, hour) => {
      const bar = document.createElement("span");
      bar.className = "bar";
      bar.style.setProperty("--price", String(price));
      const start = String(hour).padStart(2, "0");
      const end = String(hour + 1).padStart(2, "0");
      bar.title = `${start}:00–${end}:00 · ${price} €/MWh`;
      fragment.append(bar);
    });
    bars.replaceChildren(fragment);
    json.textContent = JSON.stringify(data, null, 2);
    status.textContent = "24 erfundene Stundenwerte · per API geladen";
    chart.setAttribute("aria-label", "Erfundene stündliche Strompreise in Euro pro Megawattstunde. Alle 24 Werte stehen unter dem Diagramm als JSON bereit.");
  } catch {
    bars.replaceChildren();
    status.textContent = "Beispieldaten konnten nicht geladen werden. Bitte die Seite neu laden.";
    json.textContent = "Keine Daten geladen.";
    chart.setAttribute("aria-label", "Keine Strompreise geladen.");
  } finally {
    chart.setAttribute("aria-busy", "false");
  }
}

void loadPrices();
