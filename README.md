# strom — Stromdaten als Cloud-Lernprojekt

**Status: In Entwicklung · TypeScript-Frontend und Python-Backend laufen lokal.**

Dieses Lernprojekt wächst Schritt für Schritt: von einer kleinen lokalen Anwendung zu einem reproduzierbaren Cloud-Deployment. Die Anwendungslogik bleibt bewusst klein. Der Schwerpunkt liegt auf verständlichen Schnittstellen, SQL, Docker, Linux, Terraform und dem Entwicklungs- und Deployment-Prozess.

Das Frontend besteht aus **HTML, CSS und TypeScript**, ohne UI-Framework. Es lädt Beispieldaten über HTTP vom **Python-Backend** und zeichnet ein Tagesdiagramm. Der zweite Tab erklärt den geplanten Systemaufbau mit D2-Diagrammen.

Wir entwickeln und betreiben vorerst nur dieses eine Backend. Der dokumentierte [API-Vertrag](api/README.md) ermöglicht einen späteren Wechsel, etwa zu Go oder TypeScript. Eine weitere Implementierung ergänzen wir erst bei konkretem Lernbedarf. Frühere Versuche mit anderen Backends bleiben in der Git-Historie erhalten.

Die Architekturdiagramme zeigen das Zielbild, einschließlich noch nicht implementierter Komponenten. Es gibt noch keine öffentliche Live-Anwendung, keine echten Stromdaten und keine Datenbank.

## Lokal starten

Benötigt Python 3.11+ sowie Node.js 24+ mit npm zum Bauen des TypeScript-Frontends. Node.js wird nicht als Backend verwendet.

Nach dem Klonen aus dem Projektordner starten:

```sh
npm ci
npm run build
npm start
```

http://127.0.0.1:4176 öffnen. Mit Strg+C beenden. Alternativ startet `python3 backend/server.py` direkt denselben Server.

Ein anderer Port ist möglich:

```sh
PORT=8080 npm start
```

Der Server lauscht nur auf dem eigenen Rechner. Nach Änderungen an TypeScript erneut `npm run build` ausführen und die Browserseite neu laden. Nach Änderungen am Python-Code den Server neu starten.

## Aufbau

```text
strom-dashboard/
├── index.html                 # Beide Ansichten
├── styles.css                 # Layout und Farben
├── frontend/main.ts           # Daten laden und Diagramm zeichnen
├── backend/server.py          # Python-Webserver und Demo-API
├── api/
│   ├── README.md              # HTTP-Vertrag
│   └── demo-prices.json        # 24 erfundene Stundenwerte
├── tests/contract_test.py      # HTTP-Test gegen den Python-Server
├── diagrams/                  # D2-Quellen und SVGs
├── package.json
├── package-lock.json
└── tsconfig.json
```

`dist/` und `node_modules/` werden lokal erzeugt und nicht eingecheckt. Die Struktur wächst mit dem tatsächlich benötigten Code.

- `GET /api/prices` liefert die Beispieldaten. Das Frontend prüft die Antwort zur Laufzeit und zeigt Ladefehler an.
- `/healthz` liefert `{"status":"ok"}` und prüft nur den Webprozess.
- Der Server gibt ausschließlich freigegebene Dateien aus, nicht den ganzen Projektordner.
- Der API-Vertrag legt URLs, Felder, Einheiten und HTTP-Verhalten fest. Ein späterer Sprachwechsel erfordert eine neue Implementierung der Backendlogik; er ist keine automatische Übersetzung.

Python verwendet zunächst die Standardbibliothek. Der lokale Lernserver auf Basis von `http.server` ist kein Produktionsserver; den passenden Produktionsbetrieb richten wir vor dem Cloud-Deployment ein.

`"private": true` in `package.json` verhindert eine versehentliche Veröffentlichung als npm-Paket. Das GitHub-Repository bleibt öffentlich.

## Prüfen

```sh
npm run check
npm test
```

`npm test` baut TypeScript und startet einen Python-Testserver auf einem freien lokalen Port. Der Test prüft JSON-Daten, Frontenddateien, Health-Endpunkt, HEAD-Anfragen, nicht öffentliche Pfade und abgewiesene Schreibanfragen. Er lädt keine externen Stromdaten und beendet den Testserver wieder.

## Lernpfad

- [x] Kleines Dashboard und den Tab „So funktioniert’s“ mit D2-Diagrammen gestalten.
- [x] TypeScript-Frontend an ein Python-Backend mit dokumentierter HTTP-API anbinden.
- [ ] Echte Stromdaten aus einer öffentlichen API abrufen und Fehlerfälle verstehen.
- [ ] Daten in PostgreSQL speichern und im Dashboard anzeigen.
- [ ] Anwendung und Datenbank lokal mit Docker Compose betreiben.
- [ ] Cloud-Infrastruktur mit Terraform aufsetzen und die Anwendung auf Linux deployen.
- [ ] Tests und Deployment mit GitHub Actions automatisieren; Betrieb und Wiederherstellung dokumentieren.

Kubernetes gehört vorerst nicht zum Umfang. Bei einem späteren Backendwechsel müssen neben der HTTP-Schnittstelle auch Datenbankschema und Fachlogik zusammenpassen.

## D2-Diagramme aktualisieren

Die vier SVGs werden mit D2 0.9.0 erzeugt. Jede Quelldatei legt `theme-id: 101` (Orange Creamsicle), `sketch: true` und den Randabstand unter `vars.d2-config` fest; der Renderbefehl übernimmt diese Einstellungen automatisch. D2 wird nur zum Bearbeiten und Rendern benötigt, nicht beim Besuch der Website:

```sh
d2 diagrams/system.d2 diagrams/system.svg
```

Entsprechend für `dataflow`, `database` und `deployment` wiederholen. Diagramme und Text werden bei Architekturänderungen gemeinsam gepflegt; es gibt noch keinen automatischen D2-Build.
