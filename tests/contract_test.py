"""Derselbe HTTP-Vertrag gegen alle drei echten Serverprozesse."""

import http.client
import json
import os
from pathlib import Path
import re
import selectors
import subprocess
import time
import unittest

ROOT = Path(__file__).resolve().parents[1]
COMMANDS = {
    "python": ["python3", "backends/python/server.py"],
    "go": [str(ROOT / "bin/strom-dashboard")],
    "typescript": ["node", "dist/backends/typescript/server.js"],
}


class BackendContract(unittest.TestCase):
    def test_all_backends(self):
        for backend, command in COMMANDS.items():
            with self.subTest(backend=backend):
                env = {**os.environ, "PORT": "0"}
                env.pop("ADDR", None)
                with subprocess.Popen(command, cwd=ROOT, env=env, stdout=subprocess.PIPE,
                                      stderr=subprocess.STDOUT, text=True) as process:
                    try:
                        port = self.wait_for_start(process)
                        self.check_contract(port)
                        print(f"  {backend}: gemeinsamer HTTP-Vertrag erfüllt", flush=True)
                    finally:
                        process.terminate()
                        try:
                            process.wait(timeout=5)
                        except subprocess.TimeoutExpired:
                            process.kill()
                            process.wait()

    def wait_for_start(self, process):
        lines = []
        with selectors.DefaultSelector() as selector:
            selector.register(process.stdout, selectors.EVENT_READ)
            deadline = time.monotonic() + 10
            while time.monotonic() < deadline:
                for _key, _event in selector.select(timeout=0.2):
                    line = process.stdout.readline()
                    lines.append(line)
                    match = re.search(r"Listening http://127\.0\.0\.1:(\d+)", line)
                    if match:
                        return int(match[1])
                if process.poll() is not None:
                    break
        self.fail("Server startet nicht: " + "".join(lines))

    def request(self, port, path, method="GET"):
        connection = http.client.HTTPConnection("127.0.0.1", port, timeout=3)
        try:
            connection.request(method, path)
            response = connection.getresponse()
            return response.status, dict(response.getheaders()), response.read()
        finally:
            connection.close()

    def check_contract(self, port):
        status, headers, body = self.request(port, "/api/prices")
        self.assertEqual(status, 200)
        self.assertEqual(headers["Content-Type"], "application/json")
        self.assertEqual(json.loads(body), json.loads((ROOT / "api/demo-prices.json").read_text()))

        status, headers, body = self.request(port, "/healthz")
        self.assertEqual(status, 200)
        self.assertEqual(headers["Content-Type"], "application/json")
        self.assertEqual(json.loads(body), {"status": "ok"})

        assets = {
            "/": ("index.html", "text/html"),
            "/styles.css": ("styles.css", "text/css"),
            "/assets/main.js": ("dist/frontend/main.js", "text/javascript"),
        }
        for name in ("system", "dataflow", "database", "deployment"):
            assets[f"/diagrams/{name}.svg"] = (f"diagrams/{name}.svg", "image/svg+xml")
            assets[f"/diagrams/{name}.d2"] = (f"diagrams/{name}.d2", "text/plain")
        for path, (filename, content_type) in assets.items():
            status, headers, body = self.request(port, path)
            self.assertEqual(status, 200, path)
            self.assertTrue(headers["Content-Type"].startswith(content_type), path)
            self.assertEqual(body, (ROOT / filename).read_bytes(), path)
            head_status, head_headers, head_body = self.request(port, path, "HEAD")
            self.assertEqual(head_status, 200, path)
            self.assertEqual(head_body, b"", path)
            self.assertEqual(int(head_headers["Content-Length"]), len(body), path)

        for path in ("/api/prices", "/healthz"):
            status, _headers, body = self.request(port, path, "HEAD")
            self.assertEqual(status, 200)
            self.assertEqual(body, b"")
        for method in ("POST", "PUT", "PATCH", "DELETE", "OPTIONS"):
            for path in ("/", "/api/prices", "/healthz"):
                status, headers, _body = self.request(port, path, method)
                self.assertEqual(status, 405, (method, path))
                self.assertEqual(set(headers["Allow"].replace(" ", "").split(",")), {"GET", "HEAD"})
        for path in ("/missing", "/main.go", "/go.mod", "/README.md", "/.env",
                     "/frontend/main.ts", "/backends/python/server.py",
                     "/dist/frontend/main.js", "/api/demo-prices.json", "/diagrams/",
                     "/diagrams/missing.svg"):
            status, _headers, _body = self.request(port, path)
            self.assertEqual(status, 404, path)


if __name__ == "__main__":
    unittest.main()
