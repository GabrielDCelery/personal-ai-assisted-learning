from http.server import BaseHTTPRequestHandler, HTTPServer


class SimpleHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header("Content-tpe", "text/html")
        self.end_headers()
        self.wfile.write(b"<h1>Hello form Docker!</h1>")


def main():
    server = HTTPServer(("0.0.0.0", 8080), SimpleHandler)
    print("Server running on port 8080..")
    server.serve_forever()


if __name__ == "__main__":
    main()
