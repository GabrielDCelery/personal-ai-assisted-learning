# Lesson 02: Images vs Containers

## Overview

Understanding the distinction between Docker images and containers is fundamental to working with Docker effectively. This lesson will clarify these core concepts and show you how they relate to each other.

## Core Concepts

### What is a Docker Image?

A Docker image is a **read-only template** that contains:

- Application code
- Runtime environment
- System tools and libraries
- Environment variables and configuration
- Metadata about the image

Think of an image as a **class** in object-oriented programming - it's a blueprint that defines what something should be.

**Key characteristics:**

- **Immutable**: Once built, images don't change
- **Layered**: Built from multiple layers stacked on top of each other
- **Shareable**: Can be distributed via Docker registries (Docker Hub, etc.)
- **Reusable**: One image can create many containers

### What is a Docker Container?

A container is a **running instance** of an image. When you run an image, Docker creates a container from it.

Think of a container as an **object** instantiated from a class - it's a living, running version of the blueprint.

**Key characteristics:**

- **Writable**: Has a thin writable layer on top of the image
- **Isolated**: Runs in its own process space
- **Ephemeral**: Can be stopped, started, deleted
- **Stateful**: Changes made during runtime (unless stored in volumes)

### The Relationship

```
Image (Blueprint)          Container (Running Instance)
     │                              │
     │                              │
     ├──────────────────────────────┤
     │                              │
     │    docker run                │
     │  ──────────────────>         │
     │                              │
     │                              │
     │                              │
   nginx:latest              my-nginx (running)
   (stored on disk)          (process in memory)
```

**Analogy:**

- **Image** = Recipe for a cake
- **Container** = The actual cake you baked from that recipe
- You can bake many cakes (containers) from one recipe (image)

## Hands-On Exercises

### Exercise 1: Exploring Images

**Step 1:** List all images on your system

```bash
docker images
# or
docker image ls
```

**What you'll see:**

- `REPOSITORY`: Name of the image
- `TAG`: Version identifier (e.g., latest, 1.0, alpine)
- `IMAGE ID`: Unique identifier (first 12 chars of SHA256 hash)
- `CREATED`: When the image was built
- `SIZE`: Disk space used

**Step 2:** Pull an image without running it

```bash
docker pull ubuntu:22.04
```

This downloads the image but doesn't create a container.

**Step 3:** Inspect an image

```bash
docker image inspect ubuntu:22.04
```

Look at the output - you'll see:

- Layers that make up the image
- Environment variables
- Default command (CMD)
- Architecture and OS

**Step 4:** Check image history (layers)

```bash
docker image history ubuntu:22.04
```

This shows how the image was built, layer by layer.

### Exercise 2: Understanding Containers

**Step 1:** Create a container from an image

```bash
docker run -d --name my-ubuntu ubuntu:22.04 sleep 3600
```

Let's break this down:

- `run`: Create and start a container
- `-d`: Detached mode (run in background)
- `--name my-ubuntu`: Give it a friendly name
- `ubuntu:22.04`: The image to use
- `sleep 3600`: Command to run (keeps container alive for 1 hour)

**Step 2:** List running containers

```bash
docker ps
```

**Step 3:** List ALL containers (including stopped)

```bash
docker ps -a
```

**Step 4:** Inspect a container

```bash
docker container inspect my-ubuntu
```

Notice the difference from image inspect - you see:

- State (running, paused, stopped)
- Network settings (IP address, ports)
- Mounts (volumes)
- Logs location

**Step 5:** See the writable layer

```bash
# Create a file inside the container
docker exec my-ubuntu touch /tmp/test-file.txt

# Check container's filesystem changes
docker diff my-ubuntu
```

The `A` prefix means "Added". This change exists only in this container's writable layer, not in the image.

### Exercise 3: One Image, Many Containers

**Step 1:** Create multiple containers from the same image

```bash
docker run -d --name nginx1 -p 8081:80 nginx:latest
docker run -d --name nginx2 -p 8082:80 nginx:latest
docker run -d --name nginx3 -p 8083:80 nginx:latest
```

**Step 2:** Verify they're all running

```bash
docker ps --filter ancestor=nginx:latest
```

**Step 3:** Check they're using the same image

```bash
docker ps -a --format "table {{.Names}}\t{{.Image}}\t{{.ID}}"
```

All three containers share the same base image layers! Docker is efficient - it doesn't duplicate the image data.

**Step 4:** Make changes to one container

```bash
# Modify nginx1's homepage
docker exec nginx1 bash -c 'echo "<h1>Container 1</h1>" > /usr/share/nginx/html/index.html'

# Check the websites
curl localhost:8081  # Shows "Container 1"
curl localhost:8082  # Shows default nginx page
```

Each container has its own writable layer!

**Step 5:** Cleanup

```bash
docker stop nginx1 nginx2 nginx3
docker rm nginx1 nginx2 nginx3
```

### Exercise 4: Container Lifecycle

**Step 1:** Create a container without starting it

```bash
docker create --name lifecycle-demo alpine:latest echo "Hello Docker"
```

**Step 2:** Check its status

```bash
docker ps -a --filter name=lifecycle-demo
```

Status: `Created` (not running)

**Step 3:** Start the container

```bash
docker start lifecycle-demo
```

**Step 4:** View the output

```bash
docker logs lifecycle-demo
```

**Step 5:** Check status again

```bash
docker ps -a --filter name=lifecycle-demo
```

Status: `Exited` (it ran and completed)

**Step 6:** Restart it

```bash
docker restart lifecycle-demo
docker logs lifecycle-demo
```

You'll see the message twice!

**Step 7:** Cleanup

```bash
docker rm lifecycle-demo
```

## Key Differences Summary

| Aspect       | Image                     | Container                                                 |
| ------------ | ------------------------- | --------------------------------------------------------- |
| **Nature**   | Template/Blueprint        | Running instance                                          |
| **State**    | Static, immutable         | Dynamic, stateful                                         |
| **Layers**   | Read-only layers          | Read-only layers + writable layer                         |
| **Storage**  | Stored on disk            | Runs in memory with disk storage                          |
| **Lifespan** | Permanent (until deleted) | Ephemeral (can be stopped/started)                        |
| **Purpose**  | Define what to run        | Actually run it                                           |
| **Sharing**  | Can be pushed/pulled      | Cannot be directly shared (but can be committed to image) |

## Common Commands Reference

### Images

```bash
docker images                    # List images
docker pull <image>             # Download image
docker rmi <image>              # Remove image
docker image inspect <image>    # View image details
docker image history <image>    # Show image layers
docker image prune              # Remove unused images
```

### Containers

```bash
docker ps                       # List running containers
docker ps -a                    # List all containers
docker run <image>             # Create and start container
docker create <image>          # Create container (don't start)
docker start <container>       # Start stopped container
docker stop <container>        # Stop running container
docker restart <container>     # Restart container
docker rm <container>          # Remove container
docker container inspect <id>  # View container details
docker diff <container>        # Show filesystem changes
```

## Challenges

### Challenge 1: Image Investigation

Pull the `postgres:16` image and answer these questions:

1. How many layers does it have?
2. What is the default command (CMD)?
3. What environment variables are set by default?
4. How large is the image?

### Challenge 2: Container States

Create a container that:

1. Uses the `alpine` image
2. Is named `state-demo`
3. Runs the command `ping localhost -c 5`
4. Practice starting, stopping, and restarting it
5. View the logs each time

### Challenge 3: Writable Layer Experiment

1. Run an `nginx` container
2. Modify `/usr/share/nginx/html/index.html` inside the container
3. Use `docker diff` to see what changed
4. Stop and remove the container
5. Start a new nginx container - is your change still there? Why or why not?

### Challenge 4: Resource Efficiency

1. Pull the `redis:latest` image
2. Create 5 containers from it (name them redis1, redis2, etc.)
3. Use `docker images` to verify you only downloaded the image once
4. Use `docker system df` to see storage usage
5. Cleanup all containers

## Solutions

<details>
<summary>Challenge 1 Solution</summary>

```bash
# Pull the image
docker pull postgres:16

# 1. Count layers
docker image history postgres:16 --no-trunc | wc -l

# 2. Find default command
docker image inspect postgres:16 | grep -A 5 "Cmd"
# Or more readable:
docker image inspect postgres:16 --format='{{.Config.Cmd}}'

# 3. Environment variables
docker image inspect postgres:16 --format='{{.Config.Env}}'

# 4. Image size
docker images postgres:16 --format "{{.Size}}"
```

</details>

<details>
<summary>Challenge 2 Solution</summary>

```bash
# Create the container
docker run --name state-demo alpine ping localhost -c 5

# View logs
docker logs state-demo

# Restart and check logs again
docker restart state-demo
docker logs state-demo
# You'll see 10 pings total (5 from first run, 5 from restart)

# Start/stop practice
docker stop state-demo
docker ps -a | grep state-demo  # Status: Exited
docker start state-demo
docker logs state-demo  # Now 15 pings!

# Cleanup
docker rm state-demo
```

</details>

<details>
<summary>Challenge 3 Solution</summary>

```bash
# 1. Run nginx
docker run -d --name nginx-test nginx

# 2. Modify the file
docker exec nginx-test bash -c 'echo "<h1>Modified</h1>" > /usr/share/nginx/html/index.html'

# 3. Check changes
docker diff nginx-test
# You'll see changes to the HTML file

# 4. Remove container
docker stop nginx-test
docker rm nginx-test

# 5. Start new container
docker run -d --name nginx-new -p 8080:80 nginx
curl localhost:8080
# Shows default nginx page!

# Why? The writable layer is deleted with the container.
# The original image remains unchanged.

# Cleanup
docker stop nginx-new
docker rm nginx-new
```

</details>

<details>
<summary>Challenge 4 Solution</summary>

```bash
# 1. Pull image
docker pull redis:latest

# 2. Create 5 containers
for i in {1..5}; do
  docker run -d --name redis$i redis:latest
done

# 3. Verify single image
docker images | grep redis
# Only one image listed

# 4. Check storage
docker system df
# Shows image is used by 5 containers but stored once

# Also try:
docker ps --format "{{.Names}}: {{.Image}}"

# 5. Cleanup
docker stop redis{1..5}
docker rm redis{1..5}
```

</details>

## Key Takeaways

1. **Images are blueprints**, containers are running instances
2. **One image** can create **many containers**
3. **Images are immutable** - changes happen in container's writable layer
4. **Containers are ephemeral** - when removed, their writable layer is lost
5. Docker is **efficient** - multiple containers share the same image layers
6. Use `docker images` for images, `docker ps` for containers
7. Understanding this relationship is crucial for managing data persistence (covered in Lesson 09)

## What's Next?

In **Lesson 03: Running Your First Container**, you'll dive deeper into:

- Container runtime options
- Interactive vs detached mode
- Port mapping in detail
- Environment variables
- Container naming strategies

Ready to continue? Let me know when you want to start Lesson 03!
