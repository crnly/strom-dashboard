"""Kleiner lokaler Lernserver. Gleiche API wie Go und TypeScript."""

import os
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from urllib.parse import urlsplit

ROOT = Path(__file__).resolve().parents[2]


def public_file(path):
    # Explizite Freigabe statt eines Dateiservers für das ganze Repository.
    files = {
        "/": ("index.html", "text/html; charset=utf-8"),
        "/styles.css": ("styles.css", "text/css; charset=utf-8"),
        "/assets/main.js": ("dist/frontend/main.js", "text/javascript; charset=utf-8"),
        "/api/prices": ("api/demo-prices.json", "application/json"),
    }
    for name in ("system", "dataflow", "database", "deployment"):
        files[f"/diagrams/{name}.svg"] = (f"diagrams/{name}.svg", "image/svg+xml")
        files[f"/diagrams/{name}.d2"] = (f"diagrams/{name}.d2", "text/plain; charset=utf-8")
    return files.get(path)


class Handler(BaseHTTPRequestHandler):
    def respond(self, status, body, content_type):
        self.send_response(status)
        self.send_header("Content-Type", content_type)
        self.send_header("Content-Length", str(len(body)))
        if status == 405:
            self.send_header("Allow", "GET, HEAD")
        self.end_headers()
        if self.command != "HEAD":
            self.wfile.write(body)

    def do_GET(self):
        path = urlsplit(self.path).path
        if path == "/healthz":
            self.respond(200, b'{"status":"ok"}\n', "application/json")
            return
        entry = public_file(path)
        if entry:
            filename, content_type = entry
            try:
                body = (ROOT / filename).read_bytes()
            except FileNotFoundError:
                self.respond(500, b"Frontend fehlt: zuerst npm run build ausfuehren.\n", "text/plain")
                return
            self.respond(200, body, content_type)
            return
        self.respond(404, b"Not Found\n", "text/plain")

    do_HEAD = do_GET

    def reject_write(self):
        self.respond(405, b"Method Not Allowed\n", "text/plain")

    do_POST = do_PUT = do_PATCH = do_DELETE = do_OPTIONS = reject_write


if __name__ == "__main__":
    if not (ROOT / "dist/frontend/main.js").is_file():
        raise SystemExit("Frontend fehlt: zuerst npm run build ausführen.")
    server = ThreadingHTTPServer(("127.0.0.1", int(os.environ.get("PORT", "4176"))), Handler)
    print(f"Listening http://127.0.0.1:{server.server_port} (Python)", flush=True)
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        pass
    finally:
        server.server_close()
