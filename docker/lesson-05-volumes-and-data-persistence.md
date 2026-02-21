# Lesson 05: Volumes and Data Persistence

## Overview

By default, containers are **ephemeral** - when you delete a container, all its data is lost. But what about databases, user uploads, logs, or configuration files you want to keep? This lesson teaches you how to persist data beyond a container's lifecycle using Docker volumes and bind mounts.

## Core Concepts

### The Container Filesystem Problem

When you write files inside a container, they're stored in the container's **writable layer**:

```
┌─────────────────────────────────┐
│  Writable Layer (Container)     │  ← Your data is here
├─────────────────────────────────┤
│  Read-only Layers (Image)       │
└─────────────────────────────────┘
```

**Problem:** When you remove the container, the writable layer is **deleted forever**!

```bash
docker run --name db postgres:16
# Write data to database
docker rm db  # ❌ All data is LOST!
```

### Docker's Solution: Volumes

Docker provides three ways to persist data:

1. **Volumes** (Recommended) - Managed by Docker
2. **Bind Mounts** - Map to host filesystem paths
3. **tmpfs mounts** - Stored in memory only (Linux only)

```
┌──────────────────────────────────────────────────┐
│                    Host System                   │
│                                                  │
│  ┌────────────────┐         ┌────────────────┐  │
│  │  Docker Volume │         │  Host Directory│  │
│  │  (managed)     │         │  /path/on/host │  │
│  └────────┬───────┘         └────────┬───────┘  │
│           │                          │          │
│           └──────────┬───────────────┘          │
│                      │                          │
│              ┌───────▼────────┐                 │
│              │   Container    │                 │
│              │  /data/volume  │                 │
│              └────────────────┘                 │
└──────────────────────────────────────────────────┘
```

### Volumes vs Bind Mounts vs tmpfs

| Feature         | Volumes                    | Bind Mounts                | tmpfs                 |
| --------------- | -------------------------- | -------------------------- | --------------------- |
| **Managed by**  | Docker                     | You (filesystem paths)     | Docker (memory)       |
| **Location**    | Docker storage area        | Anywhere on host           | RAM                   |
| **Portability** | ✅ High                    | ⚠️ Path-dependent          | N/A                   |
| **Sharing**     | ✅ Easy between containers | ⚠️ Possible                | ❌ Container-specific |
| **Performance** | ✅ Optimized               | ✅ Native                  | ✅ Fastest (memory)   |
| **Backup**      | ✅ Easy                    | ⚠️ Manual                  | ❌ Lost on stop       |
| **Use Case**    | Databases, production data | Development, config files  | Temporary, secrets    |
| **Syntax**      | `-v name:/path`            | `-v /host/path:/container` | `--tmpfs /path`       |
| **Recommended** | ✅ Production              | ✅ Development             | ⚠️ Specific use cases |

## Hands-On Exercises

### Exercise 1: The Problem - Data Loss

Let's demonstrate the problem first!

**Step 1:** Run PostgreSQL without a volume

```bash
docker run -d --name db-temp \
  -e POSTGRES_PASSWORD=secret \
  postgres:16
```

**Step 2:** Create some data

```bash
# Wait for database to be ready
sleep 5

# Create a table and insert data
docker exec db-temp psql -U postgres -c "
  CREATE TABLE users (id SERIAL, name TEXT);
  INSERT INTO users (name) VALUES ('Alice'), ('Bob');
"

# Verify data exists
docker exec db-temp psql -U postgres -c "SELECT * FROM users;"
```

You should see Alice and Bob!

**Step 3:** Remove the container

```bash
docker stop db-temp
docker rm db-temp
```

**Step 4:** Start a new container with the same name

```bash
docker run -d --name db-temp \
  -e POSTGRES_PASSWORD=secret \
  postgres:16

sleep 5

# Try to query the data
docker exec db-temp psql -U postgres -c "SELECT * FROM users;"
```

**Result:** Error! The table doesn't exist. All data is **gone**! 💀

**Step 5:** Cleanup

```bash
docker stop db-temp
docker rm db-temp
```

### Exercise 2: Named Volumes (The Solution!)

**Step 1:** Create a named volume

```bash
docker volume create postgres-data
```

**Step 2:** List volumes

```bash
docker volume ls
```

You'll see your `postgres-data` volume!

**Step 3:** Inspect the volume

```bash
docker volume inspect postgres-data
```

Look at the `Mountpoint` - this is where Docker stores the data on your host.

**Step 4:** Run PostgreSQL with the volume

```bash
docker run -d --name db-persistent \
  -e POSTGRES_PASSWORD=secret \
  -v postgres-data:/var/lib/postgresql/data \
  postgres:16
```

**Format:** `-v volume-name:/container/path`

**Step 5:** Create data

```bash
sleep 5

docker exec db-persistent psql -U postgres -c "
  CREATE TABLE users (id SERIAL, name TEXT);
  INSERT INTO users (name) VALUES ('Alice'), ('Bob'), ('Charlie');
"

# Verify
docker exec db-persistent psql -U postgres -c "SELECT * FROM users;"
```

**Step 6:** Remove the container

```bash
docker stop db-persistent
docker rm db-persistent
```

**Step 7:** Start a NEW container using the SAME volume

```bash
docker run -d --name db-persistent-new \
  -e POSTGRES_PASSWORD=secret \
  -v postgres-data:/var/lib/postgresql/data \
  postgres:16

sleep 5

# Query the data
docker exec db-persistent-new psql -U postgres -c "SELECT * FROM users;"
```

**Result:** Alice, Bob, and Charlie are still there! ✅ **Data persisted!**

**Step 8:** Cleanup (keeping volume)

```bash
docker stop db-persistent-new
docker rm db-persistent-new
# Don't remove the volume yet - we'll use it later
```

### Exercise 3: Anonymous Volumes

You can let Docker create volumes automatically without naming them.

**Step 1:** Run container without specifying volume name

```bash
docker run -d --name db-anon \
  -e POSTGRES_PASSWORD=secret \
  -v /var/lib/postgresql/data \
  postgres:16
```

Notice: No volume name before the `:` - Docker generates one automatically!

**Step 2:** Find the auto-generated volume

```bash
docker inspect db-anon -f '{{range .Mounts}}{{.Name}}{{end}}'
```

You'll see a long hash like `abc123def456...`

**Step 3:** List volumes to confirm

```bash
docker volume ls
```

You'll see the anonymous volume with the hash name.

**Step 4:** Cleanup

```bash
docker stop db-anon
docker rm db-anon
```

**Problem:** Anonymous volumes are harder to track and reuse!

**Step 5:** Remove unused volumes

```bash
docker volume prune
```

This removes anonymous volumes not used by any containers.

### Exercise 4: Bind Mounts

Bind mounts map a host directory directly into the container.

**Step 1:** Create a directory on the host

```bash
mkdir -p ~/docker-lesson5/nginx-content
cd ~/docker-lesson5/nginx-content
```

**Step 2:** Create an HTML file

```bash
cat > index.html << 'EOF'
<!DOCTYPE html>
<html>
<head><title>My Custom Page</title></head>
<body>
  <h1>Hello from bind mount!</h1>
  <p>This file is on my host machine.</p>
</body>
</html>
EOF
```

**Step 3:** Run nginx with a bind mount

```bash
docker run -d --name nginx-bind \
  -p 8080:80 \
  -v ~/docker-lesson5/nginx-content:/usr/share/nginx/html \
  nginx:latest
```

**Format:** `-v /host/path:/container/path`

**Step 4:** Test it

```bash
curl localhost:8080
```

You'll see your custom HTML!

**Step 5:** Modify the file on the host

```bash
echo "<h2>Updated live!</h2>" >> ~/docker-lesson5/nginx-content/index.html
```

**Step 6:** Check again

```bash
curl localhost:8080
```

The changes appear **immediately** - no container restart needed! ✨

**Step 7:** Verify files are on the host

```bash
ls -la ~/docker-lesson5/nginx-content/
cat ~/docker-lesson5/nginx-content/index.html
```

Files exist on both host AND in the container!

**Step 8:** Cleanup

```bash
docker stop nginx-bind
docker rm nginx-bind
```

### Exercise 5: Volume Mounting Options

**Read-only volumes** prevent the container from modifying data.

**Step 1:** Create a config file

```bash
mkdir -p ~/docker-lesson5/config
echo "SECRET_KEY=my-secret-key" > ~/docker-lesson5/config/app.conf
```

**Step 2:** Mount as read-only

```bash
docker run -d --name app-readonly \
  -v ~/docker-lesson5/config:/app/config:ro \
  alpine:latest sleep 3600
```

Notice the `:ro` suffix (read-only)!

**Step 3:** Try to write (should fail)

```bash
docker exec app-readonly sh -c 'echo "test" > /app/config/new.txt'
```

**Result:** Permission denied! Container can't write to read-only volumes.

**Step 4:** Read is allowed

```bash
docker exec app-readonly cat /app/config/app.conf
```

This works fine!

**Step 5:** Cleanup

```bash
docker stop app-readonly
docker rm app-readonly
```

### Exercise 6: Sharing Volumes Between Containers

Multiple containers can share the same volume!

**Step 1:** Create a shared volume

```bash
docker volume create shared-data
```

**Step 2:** Run a writer container

```bash
docker run -d --name writer \
  -v shared-data:/data \
  alpine:latest sh -c 'while true; do date >> /data/log.txt; sleep 2; done'
```

This writes timestamps to `/data/log.txt` every 2 seconds.

**Step 3:** Run a reader container

```bash
docker run -d --name reader \
  -v shared-data:/data \
  alpine:latest sh -c 'tail -f /data/log.txt'
```

**Step 4:** Check the reader's logs

```bash
docker logs reader -f
```

You'll see timestamps being written by the `writer` container in real-time! 🔄

Press `Ctrl+C` to stop watching logs.

**Step 5:** Add more readers or writers

```bash
docker run --rm -v shared-data:/data alpine:latest ls -la /data
docker run --rm -v shared-data:/data alpine:latest cat /data/log.txt
```

**Step 6:** Cleanup

```bash
docker stop writer reader
docker rm writer reader
docker volume rm shared-data
```

### Exercise 7: Backup and Restore Volumes

**Backing up a volume:**

**Step 1:** Create a volume with data

```bash
docker volume create backup-demo
docker run --rm -v backup-demo:/data alpine:latest sh -c 'echo "Important data" > /data/file.txt'
```

**Step 2:** Backup the volume to a tar file

```bash
docker run --rm \
  -v backup-demo:/data \
  -v $(pwd):/backup \
  alpine:latest tar czf /backup/backup-demo.tar.gz -C /data .
```

**What happened?**

- Mounted the volume to `/data`
- Mounted current directory to `/backup`
- Created a tar archive of `/data` contents

**Step 3:** Verify backup file

```bash
ls -lh backup-demo.tar.gz
```

**Restoring a volume:**

**Step 4:** Create a new empty volume

```bash
docker volume create backup-demo-restored
```

**Step 5:** Restore from the tar file

```bash
docker run --rm \
  -v backup-demo-restored:/data \
  -v $(pwd):/backup \
  alpine:latest tar xzf /backup/backup-demo.tar.gz -C /data
```

**Step 6:** Verify the data was restored

```bash
docker run --rm -v backup-demo-restored:/data alpine:latest cat /data/file.txt
```

You should see "Important data"! ✅

**Step 7:** Cleanup

```bash
docker volume rm backup-demo backup-demo-restored
rm backup-demo.tar.gz
```

### Exercise 8: Volume Drivers (Brief Introduction)

Docker supports different volume drivers for advanced use cases.

**Step 1:** List available volume drivers

```bash
docker volume create --help | grep driver
```

**Step 2:** Create a local volume (default)

```bash
docker volume create --driver local my-local-volume
```

**Advanced drivers** (not covered in detail):

- **NFS**: Network File System for shared storage
- **AWS EBS**: Amazon Elastic Block Store
- **Azure File Storage**: Microsoft Azure file shares
- **GlusterFS, Ceph**: Distributed storage systems

These are used in production for:

- Multi-host deployments
- Cloud storage integration
- High availability setups

We'll stick with the `local` driver for this lesson!

**Step 3:** Cleanup

```bash
docker volume rm my-local-volume
```

## Common Commands Reference

### Volume Management

```bash
docker volume ls                          # List volumes
docker volume create <name>               # Create a volume
docker volume inspect <name>              # View volume details
docker volume rm <name>                   # Remove a volume
docker volume prune                       # Remove unused volumes
docker volume prune -a                    # Remove all unused volumes
```

### Using Volumes with Containers

```bash
# Named volume
docker run -v volume-name:/container/path <image>

# Anonymous volume
docker run -v /container/path <image>

# Bind mount
docker run -v /host/path:/container/path <image>

# Read-only volume
docker run -v volume-name:/container/path:ro <image>

# Read-only bind mount
docker run -v /host/path:/container/path:ro <image>
```

### Volume Information

```bash
# See container's mounts
docker inspect <container> -f '{{.Mounts}}'

# See container's volumes only
docker inspect <container> -f '{{range .Mounts}}{{.Name}} {{.Destination}}{{"\n"}}{{end}}'

# Find which containers use a volume
docker ps -a --filter volume=<volume-name>
```

## Best Practices

### 1. Use Named Volumes for Data

```bash
# ✅ Good - easy to identify and reuse
docker run -v postgres-data:/var/lib/postgresql/data postgres:16

# ❌ Bad - anonymous volume hard to track
docker run -v /var/lib/postgresql/data postgres:16
```

### 2. Use Bind Mounts for Development

```bash
# ✅ Good - easy to edit code on host
docker run -v $(pwd)/src:/app/src node:18

# ❌ Bad - copying code into image for every change
COPY src /app/src  # in Dockerfile
```

### 3. Mount Read-Only When Possible

```bash
# ✅ Good - prevent accidental modifications
docker run -v /etc/config:/app/config:ro app:latest

# ⚠️ Risky - container can modify sensitive configs
docker run -v /etc/config:/app/config app:latest
```

### 4. Use Volumes for Databases

```dockerfile
# ✅ Good - persistent database
docker run -v db-data:/var/lib/postgresql/data postgres:16

# ❌ Bad - data lost on container removal
docker run postgres:16
```

### 5. Regularly Backup Important Volumes

```bash
# Backup
docker run --rm -v my-volume:/data -v $(pwd):/backup \
  alpine:latest tar czf /backup/my-volume-backup.tar.gz -C /data .

# Restore
docker run --rm -v my-volume:/data -v $(pwd):/backup \
  alpine:latest tar xzf /backup/my-volume-backup.tar.gz -C /data
```

### 6. Clean Up Unused Volumes

```bash
# Remove unused volumes periodically
docker volume prune

# Or with force (no confirmation)
docker volume prune -f
```

### 7. Use Volumes in docker-compose.yml

```yaml
version: "3.8"
services:
  db:
    image: postgres:16
    volumes:
      - postgres-data:/var/lib/postgresql/data

volumes:
  postgres-data:
```

## Challenges

### Challenge 1: Persistent Blog

Create a setup for a persistent blog:

1. Run a MySQL database with a volume named `blog-db`
2. Run a WordPress container connected to that database
3. Set up WordPress through the web interface
4. Remove both containers
5. Start them again and verify WordPress still has your data

### Challenge 2: Development Environment

Set up a Node.js development environment:

1. Create a simple Express app in `~/docker-lesson5/node-app`
2. Use a bind mount to mount your code into the container
3. Use a named volume for `node_modules` (to avoid overwriting)
4. Make changes to your code on the host
5. Verify changes appear in the container without rebuilding

### Challenge 3: Shared Logs

Create a logging system:

1. Create a volume named `app-logs`
2. Run 3 Alpine containers that write logs to `/logs/app.log` with different messages
3. Run a 4th container that tails the log file
4. Verify you see logs from all 3 containers

### Challenge 4: Backup and Migrate

Practice backup and migration:

1. Create a volume with sample data
2. Backup the volume to a tar file
3. Delete the original volume
4. Create a new volume with a different name
5. Restore the backup to the new volume
6. Verify data is identical

## Solutions

<details>
<summary>Challenge 1 Solution</summary>

```bash
# 1. Create volume and run MySQL
docker volume create blog-db

docker run -d --name mysql-blog \
  -e MYSQL_ROOT_PASSWORD=rootpass \
  -e MYSQL_DATABASE=wordpress \
  -e MYSQL_USER=wpuser \
  -e MYSQL_PASSWORD=wppass \
  -v blog-db:/var/lib/mysql \
  mysql:8.0

# Wait for MySQL to be ready
sleep 15

# 2. Run WordPress
docker run -d --name wordpress \
  --link mysql-blog:mysql \
  -e WORDPRESS_DB_HOST=mysql-blog \
  -e WORDPRESS_DB_USER=wpuser \
  -e WORDPRESS_DB_PASSWORD=wppass \
  -e WORDPRESS_DB_NAME=wordpress \
  -p 8080:80 \
  wordpress:latest

# 3. Visit http://localhost:8080 and set up WordPress

# 4. Remove containers
docker stop wordpress mysql-blog
docker rm wordpress mysql-blog

# 5. Start again
docker run -d --name mysql-blog \
  -e MYSQL_ROOT_PASSWORD=rootpass \
  -e MYSQL_DATABASE=wordpress \
  -e MYSQL_USER=wpuser \
  -e MYSQL_PASSWORD=wppass \
  -v blog-db:/var/lib/mysql \
  mysql:8.0

sleep 15

docker run -d --name wordpress \
  --link mysql-blog:mysql \
  -e WORDPRESS_DB_HOST=mysql-blog \
  -e WORDPRESS_DB_USER=wpuser \
  -e WORDPRESS_DB_PASSWORD=wppass \
  -e WORDPRESS_DB_NAME=wordpress \
  -p 8080:80 \
  wordpress:latest

# Visit http://localhost:8080 - your site is back!

# Cleanup
docker stop wordpress mysql-blog
docker rm wordpress mysql-blog
docker volume rm blog-db
```

</details>

<details>
<summary>Challenge 2 Solution</summary>

```bash
# 1. Create app
mkdir -p ~/docker-lesson5/node-app
cd ~/docker-lesson5/node-app

cat > package.json << 'EOF'
{
  "name": "dev-app",
  "version": "1.0.0",
  "dependencies": {
    "express": "^4.18.0"
  }
}
EOF

cat > server.js << 'EOF'
const express = require('express');
const app = express();

app.get('/', (req, res) => {
  res.send('<h1>Hello from Development!</h1>');
});

app.listen(3000, () => {
  console.log('Server running on port 3000');
});
EOF

# 2. Create named volume for node_modules
docker volume create node-app-modules

# 3. Run with bind mount + volume
docker run -d --name dev-app \
  -p 3000:3000 \
  -v $(pwd):/app \
  -v node-app-modules:/app/node_modules \
  -w /app \
  node:18 sh -c "npm install && node server.js"

# 4. Wait and test
sleep 5
curl localhost:3000

# 5. Edit server.js on host
sed -i 's/Development/Live Development/g' server.js

# 6. Restart container to see changes
docker restart dev-app
sleep 3
curl localhost:3000

# Cleanup
docker stop dev-app
docker rm dev-app
docker volume rm node-app-modules
```

</details>

<details>
<summary>Challenge 3 Solution</summary>

```bash
# 1. Create volume
docker volume create app-logs

# 2. Run 3 writer containers
docker run -d --name logger1 \
  -v app-logs:/logs \
  alpine:latest sh -c 'while true; do echo "[App1] $(date): Processing..." >> /logs/app.log; sleep 3; done'

docker run -d --name logger2 \
  -v app-logs:/logs \
  alpine:latest sh -c 'while true; do echo "[App2] $(date): Running task..." >> /logs/app.log; sleep 2; done'

docker run -d --name logger3 \
  -v app-logs:/logs \
  alpine:latest sh -c 'while true; do echo "[App3] $(date): Checking status..." >> /logs/app.log; sleep 4; done'

# 3. Run log reader
docker run -d --name log-reader \
  -v app-logs:/logs \
  alpine:latest tail -f /logs/app.log

# 4. View logs
docker logs log-reader -f

# You'll see interleaved logs from all 3 apps!
# Press Ctrl+C to stop

# Cleanup
docker stop logger1 logger2 logger3 log-reader
docker rm logger1 logger2 logger3 log-reader
docker volume rm app-logs
```

</details>

<details>
<summary>Challenge 4 Solution</summary>

```bash
# 1. Create volume with data
docker volume create original-data

docker run --rm -v original-data:/data alpine:latest sh -c '
  echo "User: Alice" > /data/users.txt
  echo "User: Bob" >> /data/users.txt
  echo "Config: production=true" > /data/config.txt
  mkdir /data/logs
  echo "Log entry 1" > /data/logs/app.log
'

# 2. Backup to tar
docker run --rm \
  -v original-data:/source \
  -v $(pwd):/backup \
  alpine:latest tar czf /backup/data-backup.tar.gz -C /source .

# 3. Delete original volume
docker volume rm original-data

# 4. Create new volume
docker volume create migrated-data

# 5. Restore to new volume
docker run --rm \
  -v migrated-data:/target \
  -v $(pwd):/backup \
  alpine:latest tar xzf /backup/data-backup.tar.gz -C /target

# 6. Verify data
docker run --rm -v migrated-data:/data alpine:latest cat /data/users.txt
docker run --rm -v migrated-data:/data alpine:latest cat /data/config.txt
docker run --rm -v migrated-data:/data alpine:latest ls -la /data/logs

# All data should match the original!

# Cleanup
docker volume rm migrated-data
rm data-backup.tar.gz
```

</details>

## Common Issues and Solutions

### Issue 1: "Error: Volume not found"

**Problem:** Trying to use a volume that doesn't exist

```bash
docker run -v nonexistent:/data alpine:latest
```

**Solution:** Docker creates volumes automatically, BUT when using `docker volume rm`, make sure the volume exists:

```bash
docker volume ls  # Check if volume exists
docker volume create my-volume  # Create if needed
```

### Issue 2: "Permission denied" when using bind mounts

**Problem:** Container can't write to host directory

**Solution:** Check directory permissions:

```bash
# Make directory writable
chmod -R 755 ~/docker-lesson5/data

# Or run container as specific user
docker run --user $(id -u):$(id -g) -v $(pwd):/data alpine:latest
```

### Issue 3: Volume data not showing up

**Problem:** Data written to volume doesn't appear

**Solution:** Check the mount path matches where the app writes data:

```bash
# Wrong - PostgreSQL writes to /var/lib/postgresql/data
docker run -v db:/data postgres:16

# Correct
docker run -v db:/var/lib/postgresql/data postgres:16
```

### Issue 4: Can't remove volume "volume is in use"

**Problem:** Volume in use by stopped container

```bash
docker volume rm my-volume
# Error: volume is in use
```

**Solution:** Remove containers using the volume first:

```bash
# Find containers using the volume
docker ps -a --filter volume=my-volume

# Remove them
docker rm <container-id>

# Now remove volume
docker volume rm my-volume
```

## Key Takeaways

1. **Containers are ephemeral** - data in the writable layer is lost when removed
2. **Use named volumes** for production data (databases, uploads, etc.)
3. **Use bind mounts** for development (live code editing)
4. **Mount read-only** when containers don't need to modify data
5. **Volumes are managed by Docker** and are the recommended approach
6. **Bind mounts** give direct access to host filesystem
7. **Multiple containers** can share the same volume
8. **Backup volumes regularly** using tar archives
9. **Clean up unused volumes** to save disk space
10. **Volume drivers** enable advanced storage backends (NFS, cloud storage)

## What's Next?

In **Lesson 06: Docker Compose Basics**, you'll learn:

- Defining multi-container applications in YAML
- Managing services, networks, and volumes together
- Development and production configurations
- Scaling services
- Docker Compose commands and workflows

Ready to continue? Let me know when you want to start Lesson 06!

---

**Status**: ✓ Lesson 05 Complete
**Next**: Lesson 06 - Docker Compose Basics
