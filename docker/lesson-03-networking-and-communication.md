# Lesson 03: Networking and Container Communication

## Overview

In the real world, applications rarely run in isolation. Your web app needs to talk to a database, your backend needs to communicate with Redis, and services need to discover each other. This lesson covers how Docker containers communicate with each other and the outside world.

## Core Concepts

### Docker Networks

Docker provides several networking drivers to connect containers:

1. **bridge** (default): Isolated network on a single host
2. **host**: Remove network isolation, use host's network directly
3. **none**: Disable all networking
4. **overlay**: Multi-host networking (for Docker Swarm)
5. **macvlan**: Assign MAC addresses to containers

We'll focus on **bridge networks** as they're the most common for development.

### Default Bridge vs Custom Bridge Networks

**Default Bridge Network:**

- Created automatically when Docker is installed
- All containers connect to it by default
- ❌ **No DNS resolution** - containers can't reach each other by name
- ✅ Can communicate by IP address
- Legacy `--link` flag required for name resolution (deprecated)

**Custom Bridge Networks:**

- Created by you with `docker network create`
- ✅ **Automatic DNS resolution** - containers can reach each other by name
- ✅ Better isolation between applications
- ✅ Can connect/disconnect containers without restarting them
- **Recommended for all multi-container applications**

### Container Communication Patterns

```
┌─────────────────────────────────────────────────────────┐
│                        Host Machine                      │
│                                                          │
│  ┌────────────────────────────────────────────────────┐ │
│  │          Custom Bridge Network: myapp-net          │ │
│  │                                                     │ │
│  │  ┌─────────────┐          ┌──────────────┐        │ │
│  │  │   nginx     │          │    redis     │        │ │
│  │  │  (port 80)  │─────────▶│  (port 6379) │        │ │
│  │  │             │  by name │              │        │ │
│  │  └──────┬──────┘          └──────────────┘        │ │
│  │         │                                          │ │
│  └─────────┼──────────────────────────────────────────┘ │
│            │                                             │
│            │ Port mapping (-p 8080:80)                   │
│            ▼                                             │
│     localhost:8080                                       │
└──────────────────────────────────────────────────────────┘
```

## Hands-On Exercises

### Exercise 1: Understanding the Default Bridge Network

**Step 1:** Check existing networks

```bash
docker network ls
```

You'll see three default networks:

- `bridge` - default network
- `host` - host network mode
- `none` - isolated network

**Step 2:** Run containers on default bridge

```bash
docker run -d --name redis-default redis:latest
docker run -d --name ubuntu-default ubuntu:24.04 sleep 3600
```

**Step 3:** Try to ping by name (this will FAIL)

```bash
# Install ping first
docker exec ubuntu-default bash -c "apt-get update -qq && apt-get install -y iputils-ping >/dev/null 2>&1"

# Try to ping redis by name
docker exec ubuntu-default ping -c 2 redis-default
```

**Result:** `ping: redis-default: Name or service not known` ❌

**Step 4:** Get Redis IP address and ping by IP

```bash
# Get Redis IP
docker inspect redis-default -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}'

# Let's say it's 172.17.0.2, ping it:
docker exec ubuntu-default ping -c 2 172.17.0.2
```

**Result:** Success! ✅ But this is not practical - IPs can change!

**Step 5:** Cleanup

```bash
docker stop redis-default ubuntu-default
docker rm redis-default ubuntu-default
```

### Exercise 2: Custom Bridge Networks (The Right Way)

**Step 1:** Create a custom bridge network

```bash
docker network create myapp-network
```

**Step 2:** Run containers on the custom network

```bash
docker run -d --name redis --network myapp-network redis:latest
docker run -d --name ubuntu --network myapp-network ubuntu:24.04 sleep 3600
```

**Step 3:** Test DNS resolution by name

```bash
# Install ping
docker exec ubuntu bash -c "apt-get update -qq && apt-get install -y iputils-ping >/dev/null 2>&1"

# Ping Redis by name
docker exec ubuntu ping -c 2 redis
```

**Result:** Success! ✅ DNS resolution works!

**Step 4:** Install and test redis-cli

```bash
# Install redis tools
docker exec ubuntu bash -c "apt-get update -qq && apt-get install -y redis-tools >/dev/null 2>&1"

# Connect to Redis by name
docker exec ubuntu redis-cli -h redis ping
```

**Result:** `PONG` ✅

**Step 5:** Inspect the network

```bash
docker network inspect myapp-network
```

Look at the `Containers` section - you'll see both containers with their IP addresses.

### Exercise 3: Connecting Existing Containers to Networks

You can connect running containers to additional networks without restarting them!

**Step 1:** Create containers on default network

```bash
docker run -d --name web nginx:latest
docker run -d --name db postgres:16 -e POSTGRES_PASSWORD=secret
```

**Step 2:** Try to connect (will fail)

```bash
docker exec web bash -c "apt-get update -qq && apt-get install -y curl >/dev/null 2>&1"
docker exec web curl db:5432
# Error: Could not resolve host: db
```

**Step 3:** Create custom network and connect containers

```bash
# Create network
docker network create app-network

# Connect both containers
docker network connect app-network web
docker network connect app-network db
```

**Step 4:** Now they can communicate!

```bash
docker exec web curl -v db:5432
# Connection successful!
```

**Step 5:** Check container's networks

```bash
docker inspect web -f '{{range $key, $value := .NetworkSettings.Networks}}{{$key}} {{end}}'
```

You'll see both `bridge` and `app-network`!

**Step 6:** Disconnect from default bridge

```bash
docker network disconnect bridge web
docker network disconnect bridge db
```

Now they're only on the custom network.

### Exercise 4: Port Mapping Deep Dive

**Port mapping** allows external access to containers.

**Step 1:** Run nginx with port mapping

```bash
docker run -d --name web1 -p 8080:80 nginx:latest
```

Format: `-p HOST_PORT:CONTAINER_PORT`

**Step 2:** Access from host

```bash
curl localhost:8080
```

You'll see the nginx welcome page!

**Step 3:** Multiple port mappings

```bash
docker run -d --name web2 \
  -p 8081:80 \
  -p 8443:443 \
  nginx:latest
```

**Step 4:** Let Docker choose the host port

```bash
docker run -d --name web3 -p 80 nginx:latest
```

**Step 5:** Find the assigned port

```bash
docker port web3
# Output: 80/tcp -> 0.0.0.0:xxxxx
```

**Step 6:** Bind to specific interface

```bash
docker run -d --name web4 -p 127.0.0.1:8082:80 nginx:latest
```

Now only accessible from localhost, not external network.

### Exercise 5: Environment Variables

**Step 1:** Pass environment variables

```bash
docker run -d --name db1 \
  -e POSTGRES_PASSWORD=mypassword \
  -e POSTGRES_USER=myuser \
  -e POSTGRES_DB=mydb \
  postgres:16
```

**Step 2:** Verify environment variables

```bash
docker exec db1 printenv | grep POSTGRES
```

**Step 3:** Use env file

Create a file `db.env`:

```bash
cat > db.env << 'EOF'
POSTGRES_PASSWORD=secret123
POSTGRES_USER=appuser
POSTGRES_DB=appdb
EOF
```

**Step 4:** Run container with env file

```bash
docker run -d --name db2 --env-file db.env postgres:16
```

**Step 5:** Check it worked

```bash
docker exec db2 printenv | grep POSTGRES
```

### Exercise 6: Real-World Multi-Container App

Let's build a complete application with a web server and database!

**Step 1:** Create network

```bash
docker network create webapp-network
```

**Step 2:** Run PostgreSQL database

```bash
docker run -d \
  --name webapp-db \
  --network webapp-network \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_USER=webuser \
  -e POSTGRES_DB=webapp \
  postgres:16
```

**Step 3:** Wait for database to be ready

```bash
docker logs webapp-db --tail 10
# Wait until you see "database system is ready to accept connections"
```

**Step 4:** Run a simple web app container

```bash
docker run -d \
  --name webapp \
  --network webapp-network \
  -p 8080:80 \
  -e DATABASE_HOST=webapp-db \
  -e DATABASE_PORT=5432 \
  -e DATABASE_NAME=webapp \
  nginx:latest
```

**Step 5:** Verify connectivity from web to db

```bash
docker exec webapp bash -c "apt-get update -qq && apt-get install -y postgresql-client >/dev/null 2>&1"
docker exec webapp psql -h webapp-db -U webuser -d webapp -c "SELECT version();"
# Enter password: secret
```

Success! The web container can reach the database by name!

## Common Commands Reference

### Network Management

```bash
docker network ls                           # List networks
docker network create <name>               # Create custom network
docker network inspect <name>              # View network details
docker network connect <network> <container>   # Connect container to network
docker network disconnect <network> <container> # Disconnect container
docker network rm <name>                   # Remove network
docker network prune                       # Remove unused networks
```

### Container Networking

```bash
docker run --network <name> <image>        # Run on specific network
docker run -p <host>:<container> <image>   # Port mapping
docker run -e KEY=VALUE <image>            # Environment variable
docker run --env-file <file> <image>       # Env file
docker port <container>                    # Show port mappings
docker inspect <container> -f '{{.NetworkSettings.Networks}}'  # Network info
```

## Networking Best Practices

1. **Always use custom bridge networks** for multi-container apps
2. **Use container names** for service discovery, not IP addresses
3. **Don't expose ports** unless necessary for external access
4. **Use environment variables** for configuration
5. **Create separate networks** for different applications
6. **Name your containers** meaningfully (db, cache, api, etc.)

## Challenges

### Challenge 1: Three-Tier Application

Build a three-tier application:

1. Create a custom network called `three-tier-net`
2. Run a Redis container (cache tier) named `cache`
3. Run an nginx container (web tier) named `web`
4. Run an Alpine container (client tier) named `client`
5. From `client`, ping both `web` and `cache` by name
6. Install redis-cli on client and connect to Redis

### Challenge 2: Network Isolation

Demonstrate network isolation:

1. Create two networks: `frontend-net` and `backend-net`
2. Run container `app` on `frontend-net`
3. Run container `db` on `backend-net`
4. Try to ping `db` from `app` (should fail)
5. Connect `app` to `backend-net` as well
6. Now ping should work

### Challenge 3: Port Mapping Scenarios

Practice different port mapping scenarios:

1. Run nginx on port 8080
2. Run another nginx on port 9090
3. Run a third nginx letting Docker choose the port
4. Run a fourth nginx accessible only from localhost:7070
5. Use `docker ps` to see all port mappings

### Challenge 4: Environment Variables

Set up a PostgreSQL database with:

1. Custom username: `admin`
2. Custom password: `supersecret`
3. Custom database name: `production`
4. Use an env file instead of multiple `-e` flags
5. Verify all variables are set correctly

## Solutions

<details>
<summary>Challenge 1 Solution</summary>

```bash
# 1. Create network
docker network create three-tier-net

# 2. Run Redis
docker run -d --name cache --network three-tier-net redis:latest

# 3. Run nginx
docker run -d --name web --network three-tier-net nginx:latest

# 4. Run Alpine
docker run -d --name client --network three-tier-net alpine:latest sleep 3600

# 5. Ping from client
docker exec client apk add --no-cache bind-tools
docker exec client ping -c 2 web
docker exec client ping -c 2 cache

# 6. Install redis-cli and connect
docker exec client apk add --no-cache redis
docker exec client redis-cli -h cache ping
# Output: PONG

# Cleanup
docker stop cache web client
docker rm cache web client
docker network rm three-tier-net
```

</details>

<details>
<summary>Challenge 2 Solution</summary>

```bash
# 1. Create networks
docker network create frontend-net
docker network create backend-net

# 2. Run app on frontend
docker run -d --name app --network frontend-net alpine:latest sleep 3600

# 3. Run db on backend
docker run -d --name db --network backend-net alpine:latest sleep 3600

# 4. Try to ping (should fail)
docker exec app apk add --no-cache bind-tools
docker exec app ping -c 2 db
# Error: bad address 'db'

# 5. Connect app to backend network
docker network connect backend-net app

# 6. Now ping works!
docker exec app ping -c 2 db
# Success!

# Cleanup
docker stop app db
docker rm app db
docker network rm frontend-net backend-net
```

</details>

<details>
<summary>Challenge 3 Solution</summary>

```bash
# 1. Nginx on 8080
docker run -d --name nginx1 -p 8080:80 nginx:latest

# 2. Nginx on 9090
docker run -d --name nginx2 -p 9090:80 nginx:latest

# 3. Let Docker choose port
docker run -d --name nginx3 -p 80 nginx:latest

# 4. Localhost only on 7070
docker run -d --name nginx4 -p 127.0.0.1:7070:80 nginx:latest

# 5. View all mappings
docker ps --format "table {{.Names}}\t{{.Ports}}"

# Test them
curl localhost:8080
curl localhost:9090
curl localhost:7070

# Find nginx3's port
docker port nginx3
curl localhost:<port>

# Cleanup
docker stop nginx1 nginx2 nginx3 nginx4
docker rm nginx1 nginx2 nginx3 nginx4
```

</details>

<details>
<summary>Challenge 4 Solution</summary>

```bash
# 1-3. Create env file
cat > prod.env << 'EOF'
POSTGRES_USER=admin
POSTGRES_PASSWORD=supersecret
POSTGRES_DB=production
EOF

# 4. Run with env file
docker run -d --name prod-db --env-file prod.env postgres:16

# 5. Verify
docker exec prod-db printenv | grep POSTGRES

# Also test connection
docker exec prod-db psql -U admin -d production -c "\l"

# Cleanup
docker stop prod-db
docker rm prod-db
rm prod.env
```

</details>

## Common Networking Issues and Solutions

### Issue 1: "Could not resolve host"

**Problem:** Container can't reach another by name

**Solution:** Check if both containers are on the same custom bridge network

```bash
docker network inspect <network-name>
```

### Issue 2: "Port already in use"

**Problem:** `-p 8080:80` fails because port 8080 is already used

**Solution:** Either stop the other container or use a different host port

```bash
docker ps --format "{{.Names}}: {{.Ports}}" | grep 8080
docker stop <container-using-8080>
```

### Issue 3: Can ping but can't connect to service

**Problem:** Container reachable but service not responding

**Solution:** Check if the service is actually running and listening

```bash
docker exec <container> netstat -tlnp
docker logs <container>
```

### Issue 4: Container can't reach the internet

**Problem:** No outbound connectivity

**Solution:** Check DNS and gateway settings

```bash
docker exec <container> cat /etc/resolv.conf
docker exec <container> ping 8.8.8.8
```

## Advanced Topics (Preview)

Topics we'll cover in future lessons:

- **Docker Compose** - Define multi-container apps in YAML
- **Volumes** - Persistent data between containers
- **Health checks** - Ensure services are actually ready
- **Container linking** - Legacy but still seen in old projects
- **Overlay networks** - Multi-host networking
- **Network plugins** - Third-party networking solutions

## Key Takeaways

1. **Default bridge network** = no DNS resolution, use custom networks instead
2. **Custom bridge networks** = automatic DNS, better isolation
3. **Container names** work as hostnames on custom networks
4. **Port mapping** (`-p`) exposes container ports to the host
5. **Environment variables** (`-e`) configure containers
6. **Multiple networks** can be attached to one container
7. **Network isolation** provides security between applications
8. Use `docker network create` for every multi-container project

## Interactive vs Detached Mode

We've been using `-d` (detached) a lot. Let's clarify:

**Detached mode (`-d`):**

```bash
docker run -d nginx
# Returns immediately, container runs in background
```

**Interactive mode (`-it`):**

```bash
docker run -it ubuntu bash
# Attaches to container's terminal, you're "inside" it
```

**Foreground mode (no -d):**

```bash
docker run nginx
# Terminal is attached to container logs, Ctrl+C stops it
```

## Container Naming Strategies

**Good names:**

```bash
docker run --name webapp-db postgres
docker run --name webapp-cache redis
docker run --name webapp-api node
```

**Bad names:**

```bash
docker run --name container1 postgres
docker run --name test redis
docker run --name asdf node
```

Name containers based on their **role** in your application!

## What's Next?

In **Lesson 04: Building Custom Images**, you'll learn:

- Writing Dockerfiles
- Building your own images
- Layer caching and optimization
- Multi-stage builds
- Best practices for image creation

Ready to continue? Let me know when you want to start Lesson 04!

---

**Status**: ✓ Lesson 03 Complete
**Next**: Lesson 04 - Building Custom Images
