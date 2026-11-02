# HTTP-Schnittstelle · Demo v1

Das TypeScript-Frontend lädt Daten vom Python-Backend. Beide kommen von derselben Adresse; das Frontend verwendet relative URLs. Ein späteres Backend in einer anderen Sprache müsste diesen Vertrag erfüllen, damit das Frontend unverändert weiterarbeiten kann. Aktuell implementieren und pflegen wir nur Python.

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

## HTTP-Verhalten

- `GET /`, `/styles.css`, `/assets/main.js` und die vier Diagramme als SVG/D2 liefern die Frontenddateien aus.
- `HEAD` liefert für diese Routen die gleichen Statuscodes und Metadaten ohne Antwortkörper.
- Unbekannte Pfade antworten mit 404. Quellcode, API-Fixture und Verzeichnislisten sind nicht direkt öffentlich erreichbar.
- `POST`, `PUT`, `PATCH`, `DELETE` und `OPTIONS` auf vorhandene Routen antworten mit 405. Es gibt keine Schreiboperationen.
- Das Frontend zeigt Ladefehler oder unerwartete Daten an, statt sie als Preise zu zeichnen.

`tests/contract_test.py` prüft diesen Vertrag gegen einen gestarteten Python-Server. TypeScript-Typen allein prüfen keine HTTP-Antworten von Python; deshalb validiert das Frontend die empfangenen Daten zusätzlich zur Laufzeit.
