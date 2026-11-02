# strom — Stromdaten als Cloud-Lernprojekt

**Status: In Entwicklung · TypeScript-Frontend und drei austauschbare Demo-Backends laufen lokal.**

Dieses Lernprojekt wächst Schritt für Schritt: von einer kleinen lokalen Anwendung zu einem reproduzierbaren Cloud-Deployment. Die Anwendungslogik bleibt bewusst klein. Der Schwerpunkt liegt auf verständlichen Schnittstellen, SQL, Docker, Linux, Terraform und dem Entwicklungs- und Deployment-Prozess.

Das Frontend besteht aus **HTML, CSS und TypeScript**, ohne UI-Framework. Es lädt Beispieldaten über HTTP und zeichnet ein Tagesdiagramm. Der zweite Tab erklärt den geplanten Systemaufbau mit D2-Diagrammen.

Als Backend wählen wir beim Start **Python, Go oder TypeScript mit Node.js**. Es läuft jeweils eine Variante. Alle liefern dieselben Dateien und erfüllen denselben [API-Vertrag](api/README.md). Das Frontend braucht beim Wechsel keine Änderung.

Die Architekturdiagramme zeigen das Zielbild, einschließlich noch nicht implementierter Komponenten. Es gibt noch keine öffentliche Live-Anwendung, keine echten Stromdaten und keine Datenbank.

## Lokal starten

Für den TypeScript-Build wird Node.js 24 oder neuer mit npm benötigt. Für das gewählte Backend zusätzlich Python 3.11+ oder Go 1.26+; das TypeScript-Backend läuft mit Node.js.

Einmal nach dem Klonen die festgeschriebenen Entwicklungsabhängigkeiten installieren und bauen:

```sh
npm ci
npm run build
```

Danach **einen** Server aus dem Projektordner starten:

| Backend | Startbefehl |
| --- | --- |
| Python · vorerst Standard | `npm start` |
| Go | `npm run start:go` |
| TypeScript / Node.js | `npm run start:typescript` |

http://127.0.0.1:4176 öffnen. Mit Strg+C beenden, dann bei Bedarf ein anderes Backend starten. Es gibt keinen Backend-Umschalter im Dashboard.

Ein anderer Port ist möglich, zum Beispiel:

```sh
PORT=8080 npm start
```

Alle Varianten lauschen standardmäßig nur auf dem eigenen Rechner. Go unterstützt weiterhin `ADDR` als vollständige Bind-Adresse; sie hat Vorrang vor `PORT`.

Nach Änderungen an TypeScript `npm run build` ausführen. Beim Go-Backend außerdem den Server neu starten, weil dessen Binärdatei HTML, CSS, JavaScript, Beispieldaten und Diagramme über `go:embed` enthält. Python und Node.js lesen diese Dateien aus dem Projektordner.

## Was gemeinsam ist — und was nicht

- `frontend/main.ts` lädt `GET /api/prices`, prüft die Antwort zur Laufzeit und zeichnet die vorhandenen Balken. Ladefehler werden sichtbar angezeigt.
- `api/demo-prices.json` enthält die 24 erfundenen Stundenwerte. Alle Backends liefern diese Datei unter derselben API-Adresse aus.
- `api/README.md` beschreibt URLs, Felder, Einheiten und HTTP-Verhalten.
- `/healthz` liefert `{"status":"ok"}` und prüft nur den Webprozess.
- Die Server geben nur explizit freigegebene Dateien aus, nicht den ganzen Projektordner.
- Die Programmlogik wird je Sprache separat implementiert. Die gemeinsame HTTP-Schnittstelle ersetzt keine automatische Übersetzung zwischen Python, Go und TypeScript.

Die Server verwenden zunächst ihre Standardbibliotheken und sind kleine lokale Lernvarianten. Insbesondere Pythons `http.server` ist kein Produktionsserver; bevor wir in die Cloud gehen, wählen wir den passenden Produktionsbetrieb für das dann aktive Backend. Die drei Varianten müssen nicht dauerhaft gleichzeitig ausgebaut werden.

`"private": true` in `package.json` verhindert eine versehentliche Veröffentlichung als npm-Paket. Das GitHub-Repository bleibt öffentlich.

## Ordnerstruktur

```text
strom-dashboard/
├── index.html                 # Beide Ansichten
├── styles.css                 # Bestehendes Layout und Farben
├── frontend/main.ts           # Gemeinsames Frontend
├── api/
│   ├── README.md              # Gemeinsamer API-Vertrag
│   └── demo-prices.json        # Gemeinsame Beispieldaten
├── backends/
│   ├── python/server.py       # Python-Variante
│   └── typescript/server.ts   # TypeScript-/Node-Variante
├── main.go                    # Bestehende Go-Variante
├── main_test.go
├── tests/contract_test.py      # Derselbe HTTP-Test für alle drei
├── diagrams/                  # D2-Quellen und SVGs
├── package.json
├── package-lock.json
├── tsconfig.json
└── go.mod
```

Die Go-Dateien bleiben vorerst an ihrem bisherigen Platz. `dist/`, `bin/` und `node_modules/` werden lokal erzeugt und nicht eingecheckt. Neue Ordner kommen erst dazu, wenn sie gebraucht werden.

## Prüfen und bauen

Typprüfung:

```sh
npm run check
```

Gemeinsamer Testlauf, benötigt alle drei Laufzeiten:

```sh
npm test
```

Der Testlauf baut TypeScript, führt Go-Tests und `go vet` aus, baut das Go-Programm und startet die drei Server nacheinander auf freien lokalen Ports. Gegen jeden Server laufen dieselben Prüfungen: identische JSON-Daten und Frontenddateien, Health-Endpunkt, HEAD-Anfragen, nicht öffentliche Pfade und abgewiesene Schreibanfragen. Er lädt keine externen Stromdaten und beendet die Testprozesse wieder.

Go einzeln bauen, nach `npm run build`:

```sh
go build -o bin/strom-dashboard .
./bin/strom-dashboard serve
```

## Lernpfad

- [x] Kleines Dashboard und den Tab „So funktioniert’s“ mit D2-Diagrammen gestalten.
- [x] Oberfläche mit einem Go-Webserver ausliefern und HTTP-Routen testen.
- [x] TypeScript-Frontend über einen gemeinsamen API-Vertrag an Python, Go und TypeScript anbinden.
- [ ] Echte Stromdaten aus einer öffentlichen API abrufen und Fehlerfälle verstehen.
- [ ] Daten in PostgreSQL speichern und im Dashboard anzeigen.
- [ ] Anwendung und Datenbank lokal mit Docker Compose betreiben.
- [ ] Cloud-Infrastruktur mit Terraform aufsetzen und die Anwendung auf Linux deployen.
- [ ] Tests und Deployment mit GitHub Actions automatisieren; Betrieb und Wiederherstellung dokumentieren.

Wir lernen zunächst mit Python weiter. Die anderen Backends können wir später anhand derselben Schnittstelle vergleichen. Ein Sprachwechsel übernimmt bereits gespeicherte Daten erst dann zuverlässig, wenn auch Datenbankschema und Fachlogik zusammenpassen. Kubernetes gehört vorerst nicht zum Umfang.

## D2-Diagramme aktualisieren

Die SVGs werden mit D2 0.9.0 erzeugt. D2 wird nur zum Bearbeiten und Rendern benötigt, nicht beim Besuch der Website:

```sh
d2 --pad 24 diagrams/system.d2 diagrams/system.svg
```

Entsprechend für `dataflow`, `database` und `deployment` wiederholen. Diagramme und Text werden bei Architekturänderungen gemeinsam gepflegt; es gibt noch keinen automatischen D2-Build.
