import os

from flask import Flask, jsonify

port = os.getenv("PORT", "8080")

app = Flask(__name__)


@app.route("/api/health")
def health():
    return jsonify({"status": "ok"})


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=port)
