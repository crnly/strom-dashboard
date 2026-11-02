# Gemeinsame HTTP-Schnittstelle · Demo v1

Python, Go und TypeScript sind alternative Implementierungen derselben Anwendung. Es läuft jeweils eines der Backends. Frontend und Backend kommen von derselben Adresse; das Frontend verwendet relative URLs und braucht beim Wechsel keine neue Konfiguration.

## `GET /api/prices`

Antwort: HTTP 200, `Content-Type: application/json`. Der Inhalt entspricht `demo-prices.json` in diesem Ordner.

| Feld | Bedeutung |
| --- | --- |
| `demo` | Immer `true`; keine echten Marktdaten |
| `date` | Fester Beispieltag `2026-10-31` |
| `market` | `DE-LU` |
| `timezone` | `Europe/Berlin` |
| `unit` | `EUR/MWh` |
| `interval` | `1 hour` |
| `prices` | 24 Zahlen, nach Stunde aufsteigend, jeweils im Bereich 0–120 |

Der Endpunkt akzeptiert noch keine Datums- oder Marktauswahl. Query-Parameter ändern die Antwort nicht. Das Zahlenarray und die feste Diagrammskala sind nur ein Demoformat. Für echte Daten ersetzen wir es gemeinsam durch explizite Zeitintervalle mit Unterstützung für negative und fehlende Preise.

## `GET /healthz`

Antwort: HTTP 200, `Content-Type: application/json`, `{"status":"ok"}`. Prüft nur, ob der Webprozess antwortet.

## Gemeinsames Verhalten

- `GET /`, `/styles.css`, `/assets/main.js` und die vier Diagramme als SVG/D2 liefern dieselben Frontenddateien aus.
- `HEAD` liefert für diese Routen die gleichen Statuscodes und Metadaten ohne Antwortkörper.
- Unbekannte Pfade antworten mit 404. Quellcode, API-Fixture und Verzeichnislisten sind nicht direkt öffentlich erreichbar.
- `POST`, `PUT`, `PATCH`, `DELETE` und `OPTIONS` auf vorhandene Routen antworten mit 405. Es gibt keine Schreiboperationen.
- Das Frontend zeigt Ladefehler oder unerwartete Daten an, statt sie als Preise zu zeichnen.

`tests/contract_test.py` prüft denselben Vertrag gegen alle drei gestarteten Backends. TypeScript-Typen allein sichern die Schnittstelle zu Python und Go nicht ab.
