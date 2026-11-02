import { createServer } from "node:http";
import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";

// Nach tsc liegt diese Datei unter dist/backends/typescript/server.js.
const root = new URL("../../../", import.meta.url);
const files = new Map<string, [string, string]>([
  ["/", ["index.html", "text/html; charset=utf-8"]],
  ["/styles.css", ["styles.css", "text/css; charset=utf-8"]],
  ["/assets/main.js", ["dist/frontend/main.js", "text/javascript; charset=utf-8"]],
  ["/api/prices", ["api/demo-prices.json", "application/json"]],
]);
for (const name of ["system", "dataflow", "database", "deployment"]) {
  files.set(`/diagrams/${name}.svg`, [`diagrams/${name}.svg`, "image/svg+xml"]);
  files.set(`/diagrams/${name}.d2`, [`diagrams/${name}.d2`, "text/plain; charset=utf-8"]);
}

const port = Number(process.env.PORT ?? "4176");
if (!Number.isInteger(port) || port < 0 || port > 65535) throw new Error("Ungültiger PORT");

const server = createServer((request, response) => {
  if (request.method !== "GET" && request.method !== "HEAD") {
    response.writeHead(405, { Allow: "GET, HEAD", "Content-Type": "text/plain" });
    response.end("Method Not Allowed\n");
    return;
  }
  const path = (request.url ?? "/").split("?")[0] ?? "/";
  let body: Buffer;
  let contentType: string;
  if (path === "/healthz") {
    body = Buffer.from('{"status":"ok"}\n');
    contentType = "application/json";
  } else {
    const file = files.get(path);
    if (!file) {
      response.writeHead(404, { "Content-Type": "text/plain" });
      response.end(request.method === "HEAD" ? undefined : "Not Found\n");
      return;
    }
    try {
      body = readFileSync(fileURLToPath(new URL(file[0], root)));
    } catch {
      response.writeHead(500, { "Content-Type": "text/plain" });
      response.end(request.method === "HEAD" ? undefined : "Datei nicht verfuegbar.\n");
      return;
    }
    contentType = file[1];
  }
  response.writeHead(200, { "Content-Type": contentType, "Content-Length": body.length });
  response.end(request.method === "HEAD" ? undefined : body);
});

server.on("error", (error) => {
  console.error(error.message);
  process.exitCode = 1;
});
server.listen(port, "127.0.0.1", () => {
  const address = server.address();
  if (address && typeof address === "object") {
    console.log(`Listening http://127.0.0.1:${address.port} (TypeScript)`);
  }
});
