# strom — Go-Webserver und UI-Prototyp

Ein kleiner Go-Webserver liefert zwei Ansichten aus: ein Tagesdiagramm mit erfundenen Preisen und eine visuelle Tour durch den geplanten Systemaufbau. Die Oberfläche besteht aus HTML und CSS. Die Navigation funktioniert über URL-Fragmente und CSS `:has()` in aktuellen Browsern.

## Lokal starten

Benötigt Go 1.26 oder neuer. Aus diesem Ordner starten:

```sh
go run .
```

Danach http://127.0.0.1:4176 öffnen. `go run . serve` startet denselben Webserver. Mit Strg+C beenden.

Die Adresse lässt sich über eine Umgebungsvariable ändern:

```sh
ADDR=127.0.0.1:8080 go run .
```

Der Server lauscht standardmäßig nur auf dem eigenen Rechner. Es gibt keine zusätzlichen Go-Abhängigkeiten. Python und Node.js werden für den Start nicht benötigt.

## Was der erste Go-Schritt macht

- `main.go` startet den HTTP-Server mit der Go-Standardbibliothek.
- `go:embed` nimmt HTML, CSS und die D2-Dateien beim Kompilieren ins Programm auf. Ein gebautes Programm kann auch außerhalb dieses Ordners gestartet werden.
- `/` liefert die Seite mit beiden Ansichten aus.
- `/styles.css` und `/diagrams/<dateiname>` liefern die benötigten öffentlichen Dateien aus.
- `/healthz` antwortet mit `{"status":"ok"}`. Das prüft nur den Webprozess, noch keine Datenbank oder Datenaktualität.
- Andere Projektdateien und Verzeichnislisten werden nicht ausgeliefert.

Nach Änderungen an HTML, CSS oder Diagrammen muss der Go-Server neu gestartet werden, weil die Dateien beim Kompilieren eingebettet werden.

## Prüfen und bauen

```sh
go test ./...
go vet ./...
go build -o /tmp/strom-dashboard .
/tmp/strom-dashboard serve
```

Die HTTP-Tests prüfen beide Ansichten, öffentliche Dateien, nicht öffentliche Projektpfade, den Health-Endpunkt und abgewiesene Schreibanfragen. Sie benötigen keinen laufenden Server und keine Internetverbindung.

## Aktueller Stand

- `index.html`: beide Ansichten, Beispieldaten und geplante Architektur.
- `styles.css`: Darstellung, Diagrammbalken, Tab-Navigation und mobile Layouts.
- `main.go`, `main_test.go`, `go.mod`: Go-Webserver und seine Tests.
- `diagrams/`: vier D2-Quellen und die daraus erzeugten SVGs für Systemarchitektur, Datenfluss, Datenmodell und Deployment. Die Grafiken und Erklärungen zeigen einen Planungsstand.
- Sämtliche Preise sind erfunden und ausdrücklich als Beispieldaten gekennzeichnet.
- PostgreSQL, API-Abruf, Docker, Terraform und Deployment sind noch nicht implementiert.
- Ein lokales Git-Repository sichert die Projektgeschichte. Auf GitHub soll es zunächst privat angelegt werden; ein Remote und eine Cloud-Veröffentlichung sind noch nicht eingerichtet.

Die Struktur wächst mit den tatsächlich umgesetzten Schritten.

## D2-Diagramme aktualisieren

Die SVGs sind mit D2 0.9.0 erzeugt. D2 wird nur zum Bearbeiten und Rendern benötigt, nicht beim Besuch der Website. Beispiel:

```sh
d2 --pad 24 diagrams/system.d2 diagrams/system.svg
```

Entsprechend für `dataflow`, `database` und `deployment` wiederholen. Diagramme und Text werden bei Architekturänderungen gemeinsam gepflegt; es gibt noch keinen automatischen D2-Build.
