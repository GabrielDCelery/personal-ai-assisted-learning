# Lesson 01: Docker Fundamentals

## What is Docker?

Docker is a platform for developing, shipping, and running applications in containers. Think of a container as a lightweight, standalone package that includes everything needed to run a piece of software: code, runtime, system tools, libraries, and settings.

## The Problem Docker Solves

### "It works on my machine!"

Before Docker, developers faced common challenges:

- Code works locally but breaks in production
- Different team members have different environments
- Setting up development environments takes hours or days
- Dependency conflicts between projects
- Difficult to scale applications

### Docker's Solution

Docker containers ensure that your application runs the same way everywhere:

- **Consistency**: Same environment from dev to production
- **Isolation**: Each container runs independently
- **Portability**: Run anywhere Docker is installed
- **Efficiency**: Lightweight compared to virtual machines

## Containers vs Virtual Machines

### Virtual Machines

```
[App A] [App B] [App C]
[Guest OS] [Guest OS] [Guest OS]
[Hypervisor]
[Host Operating System]
[Physical Server]
```

### Containers

```
[App A] [App B] [App C]
[Docker Engine]
[Host Operating System]
[Physical Server]
```

**Key Differences:**

- Containers share the host OS kernel (VMs have their own OS)
- Containers start in seconds (VMs take minutes)
- Containers use less resources (MBs vs GBs)
- Containers are more portable and easier to scale

## Core Concepts

### Images

- Read-only templates used to create containers
- Think of them as "classes" in OOP
- Stored in registries (like Docker Hub)
- Example: `nginx:latest`, `node:18`, `postgres:15`

### Containers

- Running instances of images
- Think of them as "objects" in OOP
- Can be started, stopped, deleted
- Isolated from each other and the host

### Dockerfile

- Text file with instructions to build an image
- Like a recipe for creating an image

### Registry

- Storage and distribution system for images
- Docker Hub is the public registry
- Can host private registries

## Hands-on Exercise 1: Your First Container

Let's run a simple container to see Docker in action!

### Step 1: Check Docker is running

```bash
docker --version
docker info
```

### Step 2: Run your first container

```bash
docker run hello-world
```

**What just happened?**

1. Docker checked if `hello-world` image exists locally
2. It didn't, so Docker downloaded it from Docker Hub
3. Docker created a container from the image
4. The container ran, printed a message, and exited

### Step 3: List images

```bash
docker images
```

You should see the `hello-world` image.

### Step 4: List containers

```bash
docker ps -a
```

The `-a` flag shows all containers (including stopped ones). You'll see your `hello-world` container.

## Hands-on Exercise 2: Interactive Container

Let's run something more interesting - an Ubuntu container!

### Step 1: Run Ubuntu interactively

```bash
docker run -it ubuntu:22.04 bash
```

**Flags explained:**

- `-i`: Interactive mode (keep STDIN open)
- `-t`: Allocate a pseudo-TTY (terminal)
- `bash`: Command to run inside the container

You're now inside a Ubuntu container!

### Step 2: Explore the container

Try these commands inside the container:

```bash
whoami
hostname
ls /
cat /etc/os-release
apt update && apt install -y curl
curl --version
```

### Step 3: Exit the container

```bash
exit
```

The container stops when the main process (bash) exits.

## Hands-on Exercise 3: Running a Web Server

Let's run an actual web server!

### Step 1: Run nginx

```bash
docker run -d -p 8080:80 --name my-nginx nginx:latest
```

**Flags explained:**

- `-d`: Detached mode (run in background)
- `-p 8080:80`: Port mapping (host:container)
- `--name`: Give the container a friendly name

### Step 2: Check it's running

```bash
docker ps
```

You should see your nginx container running!

### Step 3: Access the web server

Open your browser and visit: `http://localhost:8080`

Or use curl:

```bash
curl http://localhost:8080
```

You should see the nginx welcome page!

### Step 4: View logs

```bash
docker logs my-nginx
```

### Step 5: Stop and remove the container

```bash
docker stop my-nginx
docker rm my-nginx
```

## Challenge: Your Turn!

Now it's your turn to practice! Complete these tasks:

1. **Run a Python container** interactively and check the Python version:

   ```bash
   docker run -it python:3.11 bash
   python --version
   ```

2. **Run a Redis server** in the background on port 6379:
   - Name it `my-redis`
   - Use the `redis:latest` image
   - Check it's running with `docker ps`
   - View its logs with `docker logs my-redis`
   - Stop and remove it when done

3. **Experiment with container lifecycle:**
   - Run a container
   - Stop it
   - Start it again
   - Restart it
   - Remove it

## Solution to Challenge

<details>
<summary>Click to reveal solution</summary>

### Task 1: Python container

```bash
docker run -it python:3.11 bash
# Inside container:
python --version
exit
```

### Task 2: Redis server

```bash
# Run Redis
docker run -d -p 6379:6379 --name my-redis redis:latest

# Check it's running
docker ps

# View logs
docker logs my-redis

# Stop and remove
docker stop my-redis
docker rm my-redis
```

### Task 3: Container lifecycle

```bash
# Run a container
docker run -d --name lifecycle-test nginx:latest

# Stop it
docker stop lifecycle-test

# Check stopped containers
docker ps -a

# Start it again
docker start lifecycle-test

# Restart it
docker restart lifecycle-test

# Stop and remove
docker stop lifecycle-test
docker rm lifecycle-test
```

</details>

## Key Commands Learned

| Command          | Description                  |
| ---------------- | ---------------------------- |
| `docker run`     | Create and start a container |
| `docker ps`      | List running containers      |
| `docker ps -a`   | List all containers          |
| `docker images`  | List images                  |
| `docker stop`    | Stop a container             |
| `docker start`   | Start a stopped container    |
| `docker restart` | Restart a container          |
| `docker rm`      | Remove a container           |
| `docker logs`    | View container logs          |
| `docker info`    | Display system information   |

## Key Takeaways

- Docker containers are lightweight, isolated environments
- Containers are created from images
- Images are like templates, containers are running instances
- Containers are isolated but can expose ports to the host
- Docker solves the "works on my machine" problem
- Containers start quickly and use fewer resources than VMs

## Next Steps

In **Lesson 02**, we'll dive deeper into images and containers:

- How images are structured
- Layer architecture
- Building custom images
- Image tags and versioning

## Questions to Ponder

1. How is a Docker container different from installing software directly on your machine?
2. What happens to data inside a container when it's removed?
3. Can multiple containers run from the same image?
4. Why do we need port mapping (`-p` flag)?

We'll explore these questions in upcoming lessons!

---

**Status**: ✓ Lesson 01 Complete
**Next**: Lesson 02 - Images vs Containers
