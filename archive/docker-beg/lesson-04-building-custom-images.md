# Lesson 04: Building Custom Images

## Overview

So far, we've been using pre-built images from Docker Hub. But what about your own applications? In this lesson, you'll learn how to create custom Docker images using Dockerfiles, understand layer architecture, optimize builds, and apply best practices for production-ready images.

## Core Concepts

### What is a Dockerfile?

A **Dockerfile** is a text file containing instructions to build a Docker image. Think of it as a recipe that tells Docker:

- What base image to start from
- What files to copy into the image
- What commands to run
- What ports to expose
- What command to execute when the container starts

### Dockerfile Structure

```dockerfile
# Comments start with #
FROM ubuntu:22.04                    # Base image
RUN apt-get update                   # Execute commands
COPY app.py /app/                    # Copy files
WORKDIR /app                         # Set working directory
EXPOSE 8080                          # Document exposed ports
CMD ["python", "app.py"]             # Default command
```

### Common Dockerfile Instructions

| Instruction  | Purpose                                  | Example                           |
| ------------ | ---------------------------------------- | --------------------------------- |
| `FROM`       | Set base image (required, must be first) | `FROM node:18`                    |
| `RUN`        | Execute commands during build            | `RUN npm install`                 |
| `COPY`       | Copy files from host to image            | `COPY . /app`                     |
| `ADD`        | Like COPY but can extract archives       | `ADD app.tar.gz /app`             |
| `WORKDIR`    | Set working directory                    | `WORKDIR /app`                    |
| `ENV`        | Set environment variables                | `ENV NODE_ENV=production`         |
| `EXPOSE`     | Document which ports the app listens on  | `EXPOSE 3000`                     |
| `CMD`        | Default command when container starts    | `CMD ["node", "server.js"]`       |
| `ENTRYPOINT` | Configure container as executable        | `ENTRYPOINT ["python"]`           |
| `ARG`        | Build-time variables                     | `ARG VERSION=1.0`                 |
| `LABEL`      | Add metadata to image                    | `LABEL maintainer="you@mail.com"` |
| `USER`       | Set user for RUN, CMD, ENTRYPOINT        | `USER appuser`                    |

### How Docker Builds Images

Docker builds images **layer by layer**:

```
┌─────────────────────────────────┐
│  CMD ["python", "app.py"]       │  ← Layer 5 (your CMD)
├─────────────────────────────────┤
│  COPY app.py /app/              │  ← Layer 4 (your code)
├─────────────────────────────────┤
│  RUN pip install flask          │  ← Layer 3 (dependencies)
├─────────────────────────────────┤
│  WORKDIR /app                   │  ← Layer 2 (setup)
├─────────────────────────────────┤
│  FROM python:3.11-slim          │  ← Layer 1 (base image)
└─────────────────────────────────┘
```

Each instruction creates a **new layer**. Docker caches layers, so unchanged layers are reused in subsequent builds.

## Hands-On Exercises

### Exercise 1: Your First Dockerfile

Let's build a simple Python web application!

**Step 1:** Create a project directory

```bash
mkdir -p ~/docker-lesson4/hello-app
cd ~/docker-lesson4/hello-app
```

**Step 2:** Create a simple Python app

```bash
cat > app.py << 'EOF'
from http.server import HTTPServer, BaseHTTPRequestHandler

class SimpleHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()
        self.wfile.write(b'<h1>Hello from Docker!</h1>')

if __name__ == '__main__':
    server = HTTPServer(('0.0.0.0', 8080), SimpleHandler)
    print('Server running on port 8080...')
    server.serve_forever()
EOF
```

**Step 3:** Create a Dockerfile

```bash
cat > Dockerfile << 'EOF'
# Use official Python runtime as base image
FROM python:3.11-slim

# Set working directory
WORKDIR /app

# Copy application code
COPY app.py .

# Expose port
EXPOSE 8080

# Run the application
CMD ["python", "app.py"]
EOF
```

**Step 4:** Build the image

```bash
docker build -t hello-app:1.0 .
```

**What just happened?**

- `-t hello-app:1.0`: Tag the image with name and version
- `.`: Build context (current directory)

Watch the output - you'll see each layer being built!

**Step 5:** Verify the image

```bash
docker images | grep hello-app
```

**Step 6:** Run your custom image

```bash
docker run -d -p 8080:8080 --name my-hello-app hello-app:1.0
```

**Step 7:** Test it

```bash
curl localhost:8080
```

You should see: `<h1>Hello from Docker!</h1>`

**Step 8:** Cleanup

```bash
docker stop my-hello-app
docker rm my-hello-app
```

### Exercise 2: Understanding Build Cache

**Step 1:** Time your first build

```bash
time docker build -t hello-app:1.0 .
```

**Step 2:** Build again without changes

```bash
time docker build -t hello-app:1.0 .
```

Notice how much **faster** it is? That's layer caching in action!

You'll see messages like: `CACHED [2/4] WORKDIR /app`

**Step 3:** Modify the application

```bash
echo "print('Build time:', $(date))" >> app.py
```

**Step 4:** Build again

```bash
docker build -t hello-app:1.1 .
```

Notice that layers **before** `COPY app.py` are cached, but layers **after** are rebuilt!

**Step 5:** Force rebuild without cache

```bash
docker build --no-cache -t hello-app:1.2 .
```

Everything is rebuilt from scratch.

### Exercise 3: Optimizing Layer Order

Bad Dockerfile (rebuilds dependencies every time code changes):

```dockerfile
FROM node:18
WORKDIR /app
COPY . .                    # ❌ Copies everything first
RUN npm install             # ❌ Reinstalls on ANY file change
CMD ["node", "server.js"]
```

Good Dockerfile (leverages cache):

```dockerfile
FROM node:18
WORKDIR /app
COPY package*.json ./       # ✅ Copy only dependency files first
RUN npm install             # ✅ Cached unless package.json changes
COPY . .                    # ✅ Copy code last
CMD ["node", "server.js"]
```

**Step 1:** Create a Node.js project

```bash
mkdir -p ~/docker-lesson4/node-app
cd ~/docker-lesson4/node-app
```

**Step 2:** Create package.json

```bash
cat > package.json << 'EOF'
{
  "name": "node-app",
  "version": "1.0.0",
  "dependencies": {
    "express": "^4.18.0"
  }
}
EOF
```

**Step 3:** Create app

```bash
cat > server.js << 'EOF'
const express = require('express');
const app = express();

app.get('/', (req, res) => {
  res.send('<h1>Node.js in Docker!</h1>');
});

app.listen(3000, () => {
  console.log('Server running on port 3000');
});
EOF
```

**Step 4:** Create optimized Dockerfile

```bash
cat > Dockerfile << 'EOF'
FROM node:18-slim
WORKDIR /app

# Copy dependency files first
COPY package*.json ./

# Install dependencies (this layer is cached!)
RUN npm install

# Copy application code last
COPY server.js .

EXPOSE 3000
CMD ["node", "server.js"]
EOF
```

**Step 5:** Build and test

```bash
docker build -t node-app:1.0 .
docker run -d -p 3000:3000 --name node-app node-app:1.0
curl localhost:3000
docker stop node-app && docker rm node-app
```

**Step 6:** Modify code and rebuild

```bash
# Change the message
sed -i "s/Node.js/Node.js v2/g" server.js

# Rebuild
docker build -t node-app:1.1 .
```

Notice: `npm install` was **cached**! Only the code copy layer was rebuilt.

### Exercise 4: Multi-Stage Builds

Multi-stage builds create **smaller production images** by separating build and runtime environments.

**Step 1:** Create a Go application

```bash
mkdir -p ~/docker-lesson4/go-app
cd ~/docker-lesson4/go-app
```

**Step 2:** Create a simple Go app

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<h1>Go app in Docker!</h1>")
}

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("Server running on :8080")
    http.ListenAndServe(":8080", nil)
}
EOF
```

**Step 3:** Create single-stage Dockerfile (large image)

```bash
cat > Dockerfile.single << 'EOF'
FROM golang:1.21
WORKDIR /app
COPY main.go .
RUN go build -o server main.go
EXPOSE 8080
CMD ["./server"]
EOF
```

**Step 4:** Build and check size

```bash
docker build -f Dockerfile.single -t go-app:single .
docker images go-app:single
```

Notice the size: **~800MB+** (includes entire Go toolchain!)

**Step 5:** Create multi-stage Dockerfile (small image)

```bash
cat > Dockerfile << 'EOF'
# Build stage
FROM golang:1.21 AS builder
WORKDIR /app
COPY main.go .
RUN go build -o server main.go

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
EOF
```

**Step 6:** Build and compare

```bash
docker build -t go-app:multi .
docker images | grep go-app
```

Multi-stage image: **~10-15MB** - that's a **98% reduction**!

**Step 7:** Test both work the same

```bash
docker run -d -p 8081:8080 --name go-single go-app:single
docker run -d -p 8082:8080 --name go-multi go-app:multi

curl localhost:8081
curl localhost:8082

docker stop go-single go-multi
docker rm go-single go-multi
```

Both work identically, but multi-stage is **much smaller**!

### Exercise 5: Working with .dockerignore

Just like `.gitignore`, `.dockerignore` prevents unnecessary files from being copied into images.

**Step 1:** Create a project with extra files

```bash
mkdir -p ~/docker-lesson4/ignore-demo
cd ~/docker-lesson4/ignore-demo

# Create application
echo "print('Hello')" > app.py

# Create files we don't want in the image
mkdir node_modules logs .git
touch .env secrets.txt README.md
dd if=/dev/zero of=largefile.bin bs=1M count=100
```

**Step 2:** Create Dockerfile that copies everything

```bash
cat > Dockerfile << 'EOF'
FROM python:3.11-slim
WORKDIR /app
COPY . .
CMD ["python", "app.py"]
EOF
```

**Step 3:** Build and check size

```bash
docker build -t ignore-demo:without .
docker images ignore-demo:without
```

Image is **huge** because of `largefile.bin`!

**Step 4:** Create .dockerignore

```bash
cat > .dockerignore << 'EOF'
# Development files
.git
.gitignore
.env

# Dependencies
node_modules
__pycache__

# Logs and temp files
*.log
logs/
*.tmp

# Documentation
README.md
*.md

# Large files
*.bin
EOF
```

**Step 5:** Build again and compare

```bash
docker build -t ignore-demo:with .
docker images | grep ignore-demo
```

Much smaller now!

**Step 6:** Verify ignored files aren't in the image

```bash
docker run --rm ignore-demo:with ls -la
# largefile.bin is not there!
```

### Exercise 6: Environment Variables and ARG

**Step 1:** Create an app that uses environment variables

```bash
mkdir -p ~/docker-lesson4/env-app
cd ~/docker-lesson4/env-app

cat > app.py << 'EOF'
import os

env = os.getenv('ENVIRONMENT', 'development')
port = os.getenv('PORT', '8080')

print(f"Running in {env} mode on port {port}")

from http.server import HTTPServer, BaseHTTPRequestHandler

class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()
        msg = f'<h1>Environment: {env}</h1><p>Port: {port}</p>'
        self.wfile.write(msg.encode())

server = HTTPServer(('0.0.0.0', int(port)), Handler)
server.serve_forever()
EOF
```

**Step 2:** Create Dockerfile with ENV and ARG

```bash
cat > Dockerfile << 'EOF'
FROM python:3.11-slim

# Build-time variable (only during build)
ARG APP_VERSION=1.0.0

# Runtime environment variable (available in container)
ENV ENVIRONMENT=production
ENV PORT=8080

WORKDIR /app
COPY app.py .

# Use ARG in a label
LABEL version=${APP_VERSION}

EXPOSE ${PORT}

CMD ["python", "app.py"]
EOF
```

**Step 3:** Build with default ARG

```bash
docker build -t env-app:1.0 .
```

**Step 4:** Build with custom ARG

```bash
docker build --build-arg APP_VERSION=2.0.0 -t env-app:2.0 .
```

**Step 5:** Check the label

```bash
docker image inspect env-app:1.0 -f '{{.Config.Labels.version}}'
docker image inspect env-app:2.0 -f '{{.Config.Labels.version}}'
```

**Step 6:** Run with default ENV

```bash
docker run -d -p 8080:8080 --name env-default env-app:1.0
curl localhost:8080
docker logs env-default
```

**Step 7:** Override ENV at runtime

```bash
docker run -d -p 8081:9000 --name env-custom \
  -e ENVIRONMENT=staging \
  -e PORT=9000 \
  env-app:1.0

curl localhost:8081
docker logs env-custom
```

**Step 8:** Cleanup

```bash
docker stop env-default env-custom
docker rm env-default env-custom
```

## Common Commands Reference

### Building Images

```bash
docker build -t <name>:<tag> .              # Build image
docker build -t <name>:<tag> -f <file> .    # Use specific Dockerfile
docker build --no-cache -t <name> .         # Build without cache
docker build --build-arg VAR=value .        # Pass build argument
docker build --target <stage> .             # Build specific stage (multi-stage)
```

### Managing Images

```bash
docker images                               # List images
docker image ls                             # Same as above
docker image inspect <image>                # View image details
docker image history <image>                # Show layers
docker image rm <image>                     # Remove image
docker image prune                          # Remove unused images
docker image prune -a                       # Remove all unused images
docker tag <source> <target>                # Tag image
```

### Image Information

```bash
docker image inspect <image> --format='{{.Size}}'           # Image size
docker image inspect <image> --format='{{.Config.Env}}'     # Environment vars
docker image inspect <image> --format='{{.Config.Cmd}}'     # Default command
docker history <image> --no-trunc                           # Full layer history
```

## Best Practices

### 1. Use Official Base Images

```dockerfile
# ✅ Good - official, maintained, secure
FROM node:18-slim
FROM python:3.11-alpine
FROM nginx:stable-alpine

# ❌ Bad - random user image, no security updates
FROM johndoe/custom-node:latest
```

### 2. Use Specific Tags, Not `latest`

```dockerfile
# ✅ Good - reproducible builds
FROM node:18.17.0-alpine3.18

# ❌ Bad - breaks when "latest" updates
FROM node:latest
```

### 3. Minimize Layers

```dockerfile
# ❌ Bad - creates 3 layers
RUN apt-get update
RUN apt-get install -y curl
RUN apt-get clean

# ✅ Good - creates 1 layer
RUN apt-get update && \
    apt-get install -y curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
```

### 4. Order Instructions by Change Frequency

```dockerfile
# ✅ Good - dependencies cached, code changes don't reinstall
COPY package.json .
RUN npm install
COPY . .

# ❌ Bad - code changes trigger npm install
COPY . .
RUN npm install
```

### 5. Use .dockerignore

```
.git
.gitignore
node_modules
*.log
.env
.vscode
.idea
*.md
```

### 6. Don't Run as Root

```dockerfile
# ✅ Good - run as non-root user
RUN useradd -m appuser
USER appuser

# ❌ Bad - running as root is a security risk
# (no USER instruction means root)
```

### 7. Use Multi-Stage Builds

```dockerfile
# ✅ Good - small final image
FROM node:18 AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM node:18-alpine
COPY --from=builder /app/dist ./dist
CMD ["node", "dist/server.js"]
```

### 8. Combine Commands to Reduce Layers

```dockerfile
# ✅ Good - single layer
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        curl \
        git && \
    rm -rf /var/lib/apt/lists/*

# ❌ Bad - multiple layers
RUN apt-get update
RUN apt-get install -y curl
RUN apt-get install -y git
```

### 9. Use COPY Instead of ADD

```dockerfile
# ✅ Good - explicit and clear
COPY app.py /app/

# ❌ Bad - ADD has magic behavior (auto-extract, URLs)
ADD app.py /app/
```

Use `ADD` only when you need its special features (extracting tar files, fetching URLs).

### 10. Clean Up in the Same Layer

```dockerfile
# ✅ Good - cache cleaned in same layer
RUN apt-get update && \
    apt-get install -y package && \
    rm -rf /var/lib/apt/lists/*

# ❌ Bad - cache stays in the layer
RUN apt-get update
RUN apt-get install -y package
RUN rm -rf /var/lib/apt/lists/*  # Too late, layer already created
```

## Challenges

### Challenge 1: Build a Flask API

Create a custom image for a Flask REST API:

1. Create a simple Flask app that responds to `/api/health` with `{"status": "ok"}`
2. Write a Dockerfile that:
   - Uses Python 3.11 slim
   - Installs Flask using `requirements.txt`
   - Copies the application code
   - Runs on port 5000
   - Uses a non-root user
3. Build and run the image
4. Test the endpoint with curl

### Challenge 2: Optimize a Real Application

Given this inefficient Dockerfile:

```dockerfile
FROM ubuntu:22.04
RUN apt-get update
RUN apt-get install -y python3
RUN apt-get install -y python3-pip
COPY . /app
WORKDIR /app
RUN pip3 install -r requirements.txt
CMD ["python3", "app.py"]
```

Optimize it for:

- Smaller image size
- Better caching
- Faster builds
- Security (non-root user)

### Challenge 3: Multi-Stage Rust Build

Create a multi-stage build for a Rust application:

1. Build stage: Use `rust:1.73` to compile the application
2. Runtime stage: Use `debian:bookworm-slim` (Rust binaries need glibc)
3. The final image should be as small as possible
4. Compare sizes with a single-stage build

### Challenge 4: Dynamic Configuration

Create an image that:

1. Accepts build arguments: `VERSION`, `BUILD_DATE`
2. Sets environment variables: `APP_ENV`, `LOG_LEVEL`
3. Adds labels with metadata
4. Can be configured differently for dev/staging/production
5. Build three variants: dev, staging, production

## Solutions

<details>
<summary>Challenge 1 Solution</summary>

```bash
# Create project
mkdir -p ~/docker-lesson4/flask-api
cd ~/docker-lesson4/flask-api

# Create requirements.txt
cat > requirements.txt << 'EOF'
flask==3.0.0
EOF

# Create app
cat > app.py << 'EOF'
from flask import Flask, jsonify

app = Flask(__name__)

@app.route('/api/health')
def health():
    return jsonify({"status": "ok"})

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
EOF

# Create Dockerfile
cat > Dockerfile << 'EOF'
FROM python:3.11-slim

# Create non-root user
RUN useradd -m -u 1000 appuser

WORKDIR /app

# Install dependencies first (caching)
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application
COPY app.py .

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

EXPOSE 5000

CMD ["python", "app.py"]
EOF

# Build
docker build -t flask-api:1.0 .

# Run
docker run -d -p 5000:5000 --name flask-api flask-api:1.0

# Test
curl localhost:5000/api/health

# Cleanup
docker stop flask-api && docker rm flask-api
```

</details>

<details>
<summary>Challenge 2 Solution</summary>

```dockerfile
# Optimized Dockerfile
FROM python:3.11-slim

# Install system dependencies in one layer
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        build-essential && \
    rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd -m -u 1000 appuser

WORKDIR /app

# Copy and install dependencies first (better caching)
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code last
COPY app.py .

# Set ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

CMD ["python", "app.py"]
```

**Improvements:**

- Used `python:3.11-slim` instead of Ubuntu (smaller base)
- Combined RUN commands (fewer layers)
- Ordered COPY for better caching
- Added non-root user
- Cleaned apt cache in same layer
- Used `--no-cache-dir` for pip

</details>

<details>
<summary>Challenge 3 Solution</summary>

```bash
# Create project
mkdir -p ~/docker-lesson4/rust-app
cd ~/docker-lesson4/rust-app

# Create simple Rust app
cat > main.rs << 'EOF'
use std::net::TcpListener;
use std::io::Write;

fn main() {
    let listener = TcpListener::bind("0.0.0.0:8080").unwrap();
    println!("Server running on port 8080");

    for stream in listener.incoming() {
        let mut stream = stream.unwrap();
        let response = "HTTP/1.1 200 OK\r\n\r\n<h1>Rust in Docker!</h1>";
        stream.write_all(response.as_bytes()).unwrap();
    }
}
EOF

# Single-stage Dockerfile
cat > Dockerfile.single << 'EOF'
FROM rust:1.73
WORKDIR /app
COPY main.rs .
RUN rustc main.rs -o server
CMD ["./server"]
EOF

# Multi-stage Dockerfile
cat > Dockerfile << 'EOF'
# Build stage
FROM rust:1.73 AS builder
WORKDIR /app
COPY main.rs .
RUN rustc -C opt-level=3 main.rs -o server

# Runtime stage
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
EOF

# Build both
docker build -f Dockerfile.single -t rust-app:single .
docker build -t rust-app:multi .

# Compare sizes
docker images | grep rust-app

# Test multi-stage
docker run -d -p 8080:8080 --name rust-app rust-app:multi
curl localhost:8080
docker stop rust-app && docker rm rust-app
```

**Result:** Multi-stage is ~100MB vs 1.5GB+ for single-stage!

</details>

<details>
<summary>Challenge 4 Solution</summary>

```dockerfile
# Dockerfile with build args and env vars
ARG PYTHON_VERSION=3.11
FROM python:${PYTHON_VERSION}-slim

# Build arguments
ARG VERSION=1.0.0
ARG BUILD_DATE
ARG APP_ENV=production

# Labels (metadata)
LABEL version="${VERSION}" \
      build-date="${BUILD_DATE}" \
      environment="${APP_ENV}"

# Runtime environment variables
ENV APP_ENV=${APP_ENV} \
    LOG_LEVEL=info \
    PORT=8080

WORKDIR /app
COPY app.py .

EXPOSE ${PORT}
CMD ["python", "app.py"]
```

```bash
# Build different variants
docker build \
  --build-arg VERSION=1.0.0 \
  --build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --build-arg APP_ENV=development \
  -t myapp:dev .

docker build \
  --build-arg VERSION=1.0.0 \
  --build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --build-arg APP_ENV=staging \
  -t myapp:staging .

docker build \
  --build-arg VERSION=1.0.0 \
  --build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --build-arg APP_ENV=production \
  -t myapp:prod .

# Inspect labels
docker image inspect myapp:dev -f '{{.Config.Labels}}'
docker image inspect myapp:prod -f '{{.Config.Labels}}'

# Run with different configs
docker run -d -p 8080:8080 --name app-dev myapp:dev
docker run -d -p 8081:8080 --name app-prod -e LOG_LEVEL=error myapp:prod

docker stop app-dev app-prod
docker rm app-dev app-prod
```

</details>

## Key Takeaways

1. **Dockerfiles** define how to build images from instructions
2. **Layer caching** makes rebuilds fast - order matters!
3. **Multi-stage builds** create smaller production images
4. **Copy dependencies first**, code last for optimal caching
5. **Use .dockerignore** to keep images small and builds fast
6. **Don't run as root** - create a user with `USER`
7. **Use specific tags**, not `latest`
8. **Combine RUN commands** to minimize layers
9. **Clean up in the same layer** where you create temporary files
10. **ARG** for build-time, **ENV** for runtime configuration

## Common Dockerfile Pitfalls

### Pitfall 1: Cache Busting

```dockerfile
# ❌ This busts cache on every build
RUN apt-get update && apt-get install -y curl && \
    echo "Build date: $(date)"
```

### Pitfall 2: Secrets in Images

```dockerfile
# ❌ NEVER do this - secrets end up in image history
COPY .env /app/.env
RUN echo "API_KEY=secret123" > /app/secrets.txt
```

Use environment variables at runtime or Docker secrets instead!

### Pitfall 3: Large Images

```dockerfile
# ❌ Installs unnecessary files
FROM ubuntu:22.04
RUN apt-get install -y python3 python3-pip build-essential

# ✅ Use slim variants and clean up
FROM python:3.11-slim
RUN apt-get update && \
    apt-get install -y --no-install-recommends build-essential && \
    rm -rf /var/lib/apt/lists/*
```

### Pitfall 4: Running Multiple Processes

```dockerfile
# ❌ Docker containers should run ONE process
CMD service nginx start && service mysql start && python app.py
```

Use Docker Compose or Kubernetes for multi-process applications!

## What's Next?

In **Lesson 05: Docker Volumes and Data Persistence**, you'll learn:

- How to persist data beyond container lifecycle
- Volume types (named volumes, bind mounts, tmpfs)
- Sharing data between containers
- Backup and restore strategies
- Database container best practices

---

**Status**: ✓ Lesson 04 Complete
**Next**: Lesson 05 - Docker Volumes and Data Persistence
