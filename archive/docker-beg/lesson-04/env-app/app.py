import os

env = os.getenv("ENVIRONMENT", "development")
port = os.getenv("PORT", "8080")

print(f"Running in {env} mode on port {port}")

from http.server import BaseHTTPRequestHandler, HTTPServer


class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header("Content-type", "text/html")
        self.end_headers()
        msg = f"<h1>Environment: {env}</h1><p>Port: {port}</p>"
        self.wfile.write(msg.encode())


server = HTTPServer(("0.0.0.0", int(port)), Handler)
server.serve_forever()
