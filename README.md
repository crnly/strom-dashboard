# strom — erster UI-Prototyp

Zwei Ansichten mit HTML und CSS: ein Tagesdiagramm mit erfundenen Preisen und eine ausführliche visuelle Tour durch den geplanten Systemaufbau. Die Navigation funktioniert über URL-Fragmente und CSS `:has()` in aktuellen Browsern. Zum Ansehen sind JavaScript, Build und Installation nicht nötig.

## Ansehen

`index.html` im Browser öffnen oder aus diesem Ordner starten:

```sh
python3 -m http.server 4175 --bind 127.0.0.1
```

Danach http://127.0.0.1:4175 öffnen. Python dient hier nur als lokaler Dateiserver. Das spätere Backend wird in Go entwickelt.

## Aktueller Stand

- `index.html`: beide Ansichten, Beispieldaten und geplante Architektur.
- `styles.css`: Darstellung, Diagrammbalken, Tab-Navigation und mobile Layouts.
- `diagrams/`: vier D2-Quellen und die daraus erzeugten SVGs für Systemarchitektur, Datenfluss, Datenmodell und Deployment. Die Grafiken und Erklärungen zeigen einen Planungsstand.
- Sämtliche Preise sind erfunden und ausdrücklich als Beispieldaten gekennzeichnet.
- Go, PostgreSQL, API-Abruf, Docker, Terraform und Deployment sind noch nicht implementiert.
- Das Projekt ist lokal; es ist noch kein GitHub-Repository und keine Cloud-Veröffentlichung eingerichtet.

Die Struktur wächst mit den tatsächlich umgesetzten Schritten.

## D2-Diagramme aktualisieren

Die SVGs sind mit D2 0.9.0 erzeugt. D2 wird nur zum Bearbeiten und Rendern benötigt, nicht beim Besuch der Website. Beispiel:

```sh
d2 --pad 24 diagrams/system.d2 diagrams/system.svg
```

Entsprechend für `dataflow`, `database` und `deployment` wiederholen. Diagramme und Text werden bei Architekturänderungen gemeinsam gepflegt; es gibt noch keinen automatischen D2-Build.
