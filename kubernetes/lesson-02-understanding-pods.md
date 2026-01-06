# Lesson 02: Understanding Pods

## Overview

In Lesson 01, you learned that pods are the smallest deployable unit in Kubernetes. You created pods indirectly through Deployments, but didn't explore what pods actually are or why they're designed this way.

In this lesson, you'll understand pods deeply: their lifecycle, design patterns, how multiple containers share resources within a pod, and when to use single vs multi-container pods.

By the end of this lesson, you'll be able to create pods directly, debug pod issues, and make informed decisions about pod design.

---

## What is a Pod?

A **pod** is the smallest deployable unit in Kubernetes. It's a wrapper around one or more containers that share:

- **Network namespace**: Same IP address, can communicate via `localhost`
- **Storage volumes**: Shared file systems
- **IPC namespace**: Can use inter-process communication
- **UTS namespace**: Same hostname

### Docker Container vs Kubernetes Pod

```
Docker World:
┌──────────────┐
│  Container   │ ← You deploy this directly
└──────────────┘

Kubernetes World:
┌─────────────────────────────┐
│          Pod                │ ← You deploy this
│  ┌──────────────┐           │
│  │  Container   │           │
│  └──────────────┘           │
└─────────────────────────────┘
```

**Key insight**: In Kubernetes, you never run a container directly. You always create a pod, which then runs your container(s).

### Why Pods Exist

**Question**: Why not just deploy containers directly like Docker?

**Answer**: Pods solve the "tightly coupled containers" problem.

#### Scenario: Web App + Log Collector

Imagine you have:

- A web application container
- A log collection sidecar that reads logs and ships them to a central system

These containers need to:

- Share a volume (the log directory)
- Start and stop together
- Run on the same machine (for efficient communication)

**Docker solution**: Manual orchestration, custom scripts
**Kubernetes solution**: Put both containers in the same pod

---

## Anatomy of a Pod

### Single Container Pod (Most Common)

This is what you'll use 90% of the time.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
  labels:
    app: nginx
    env: dev
spec:
  containers:
    - name: nginx
      image: nginx:1.25
      ports:
        - containerPort: 80
```

**Breakdown**:

- `apiVersion: v1`: Pods are part of the core Kubernetes API
- `kind: Pod`: We're defining a pod resource
- `metadata`: Name and labels for organization
- `spec.containers`: List of containers (here, just one)

### Multi-Container Pod

Multiple containers sharing the same network and storage.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: web-with-logging
spec:
  containers:
    # Main application
    - name: web-app
      image: myapp:1.0
      ports:
        - containerPort: 8080
      volumeMounts:
        - name: shared-logs
          mountPath: /var/log/app

    # Sidecar: Log shipper
    - name: log-shipper
      image: fluent/fluentd:latest
      volumeMounts:
        - name: shared-logs
          mountPath: /var/log/app
          readOnly: true

  # Shared volume
  volumes:
    - name: shared-logs
      emptyDir: {}
```

**What's happening**:

1. Both containers run in the same pod
2. They share the `shared-logs` volume
3. `web-app` writes logs to `/var/log/app`
4. `log-shipper` reads from the same directory
5. They can communicate via `localhost`

---

## Pod Lifecycle

Pods go through several states during their lifetime.

### Pod Phases

| Phase       | Description                                                  |
| ----------- | ------------------------------------------------------------ |
| `Pending`   | Pod accepted by Kubernetes, but container(s) not yet created |
| `Running`   | Pod bound to a node, at least one container is running       |
| `Succeeded` | All containers terminated successfully (for Jobs)            |
| `Failed`    | All containers terminated, at least one failed               |
| `Unknown`   | Pod state cannot be determined (communication error)         |
| `Completed` | Pod has completed its work (used with Jobs and CronJobs)     |

### Container States (within a Pod)

| State        | Description                                   |
| ------------ | --------------------------------------------- |
| `Waiting`    | Container is waiting to start (pulling image) |
| `Running`    | Container is executing                        |
| `Terminated` | Container finished execution or crashed       |

### Pod Lifecycle Diagram

```
Created → Pending → ContainerCreating → Running → [Succeeded/Failed]
                                            ↓
                                       (if crashes)
                                            ↓
                                       CrashLoopBackOff
```

---

## Hands-On Exercise 1: Create Your First Pod Directly

Previously, you created pods via Deployments. Now you'll create a pod directly.

### Create a Simple Pod

**nginx-pod.yaml**:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-first-pod
  labels:
    app: nginx
    tier: frontend
spec:
  containers:
    - name: nginx-container
      image: nginx:1.25
      ports:
        - containerPort: 80
          protocol: TCP
```

### Deploy the Pod

```bash
# Apply the YAML
kubectl apply -f nginx-pod.yaml

# Check pod status
kubectl get pods

# Expected output:
# NAME            READY   STATUS    RESTARTS   AGE
# my-first-pod    1/1     Running   0          10s

# Get detailed information
kubectl describe pod my-first-pod
```

### Understanding the Output

```bash
kubectl get pods my-first-pod

# READY: 1/1 means 1 container ready out of 1 total
# STATUS: Current phase (Running, Pending, etc.)
# RESTARTS: How many times containers have restarted
# AGE: Time since pod creation
```

### Access the Pod

```bash
# Get the pod's IP address
kubectl get pod my-first-pod -o wide

# Expected output:
# NAME           READY   STATUS    IP           NODE
# my-first-pod   1/1     Running   10.244.0.5   minikube

# Port-forward to access from your machine
kubectl port-forward pod/my-first-pod 8080:80

# In another terminal or browser, visit:
# http://localhost:8080
# You should see the nginx welcome page

# Press Ctrl+C to stop port-forwarding
```

---

## Hands-On Exercise 2: Inspect Pod Details

### View Pod Logs

```bash
# View logs from the nginx container
kubectl logs my-first-pod

# Expected output: nginx access logs
# 10.244.0.1 - - [05/Jan/2026:10:23:45 +0000] "GET / HTTP/1.1" 200 615

# Follow logs in real-time
kubectl logs my-first-pod -f

# (Press Ctrl+C to stop)
```

### Execute Commands Inside the Pod

```bash
# Get a shell inside the container
kubectl exec -it my-first-pod -- /bin/bash

# You're now inside the container
root@my-first-pod:/# ls
bin  boot  dev  docker-entrypoint.d  docker-entrypoint.sh  etc  home  lib  ...

# Check nginx process
root@my-first-pod:/# ps aux | grep nginx

# View nginx config
root@my-first-pod:/# cat /etc/nginx/nginx.conf

# Exit the container
root@my-first-pod:/# exit
```

### Inspect Pod with kubectl describe

```bash
kubectl describe pod my-first-pod
```

**Key sections in the output**:

- **Name, Namespace, Labels**: Identification
- **Status**: Current phase
- **IP**: Pod's cluster IP
- **Containers**: Details of each container
- **Conditions**: Pod readiness checks
- **Events**: Recent activities (image pull, container start, etc.)

---

## Hands-On Exercise 3: Multi-Container Pod (Sidecar Pattern)

Let's create a pod with two containers that work together.

### Scenario: Web Server + Request Logger

**multi-container-pod.yaml**:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: web-logger-pod
spec:
  containers:
    # Main container: nginx web server
    - name: nginx
      image: nginx:1.25
      ports:
        - containerPort: 80
      volumeMounts:
        - name: shared-logs
          mountPath: /var/log/nginx

    # Sidecar container: log watcher
    - name: log-watcher
      image: busybox
      command: ["sh", "-c", "tail -f /var/log/nginx/access.log"]
      volumeMounts:
        - name: shared-logs
          mountPath: /var/log/nginx

  # Shared volume between containers
  volumes:
    - name: shared-logs
      emptyDir: {}
```

### Deploy and Test

```bash
# Create the pod
kubectl apply -f multi-container-pod.yaml

# Check pod status (should show 2/2 containers ready)
kubectl get pod web-logger-pod

# Expected output:
# NAME              READY   STATUS    RESTARTS   AGE
# web-logger-pod    2/2     Running   0          15s

# Port-forward to access nginx
kubectl port-forward pod/web-logger-pod 8080:80

# In another terminal, make a request
curl http://localhost:8080

# View logs from the log-watcher container
kubectl logs web-logger-pod -c log-watcher

# You'll see the access log from the curl request!
```

### Why This Works

1. Both containers share the `shared-logs` volume
2. nginx writes access logs to `/var/log/nginx/access.log`
3. log-watcher reads from the same file via the shared volume
4. This pattern is called a **sidecar**: a helper container alongside the main app

---

## Common Multi-Container Patterns

### 1. Sidecar Pattern

A helper container that extends the main container's functionality.

**Use cases**:

- Log shipping (Fluentd, Filebeat)
- Metrics collection (Prometheus exporters)
- Configuration synchronization

```
┌──────────────────────────┐
│         Pod              │
│  ┌────────┐  ┌────────┐ │
│  │  App   │  │Sidecar │ │
│  └────────┘  └────────┘ │
│       Shared Volume      │
└──────────────────────────┘
```

### 2. Ambassador Pattern

A proxy container that handles communication to external services.

**Use cases**:

- Database connection pooling
- API rate limiting
- Service mesh proxies (Istio, Linkerd)

```
┌─────────────────────────────┐
│          Pod                │
│  ┌────────┐   ┌──────────┐ │
│  │  App   │──→│Ambassador│─┼─→ External Service
│  └────────┘   └──────────┘ │
└─────────────────────────────┘
```

### 3. Adapter Pattern

A container that transforms the main container's output to a standard format.

**Use cases**:

- Log format standardization
- Metrics format conversion
- Data transformation

```
┌──────────────────────────────┐
│          Pod                 │
│  ┌────────┐    ┌─────────┐  │
│  │  App   │───→│ Adapter │──┼─→ Monitoring System
│  │(custom)│    │(standard)│  │
│  └────────┘    └─────────┘  │
└──────────────────────────────┘
```

---

## Pod Configuration Options

### Environment Variables

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: env-pod
spec:
  containers:
    - name: app
      image: nginx
      env:
        - name: ENVIRONMENT
          value: "production"
        - name: LOG_LEVEL
          value: "info"
```

### Resource Requests and Limits

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: resource-pod
spec:
  containers:
    - name: app
      image: nginx
      resources:
        requests:
          memory: "64Mi"
          cpu: "250m"
        limits:
          memory: "128Mi"
          cpu: "500m"
```

**Explanation**:

- `requests`: Guaranteed resources (scheduler uses this)
- `limits`: Maximum resources allowed
- `cpu`: In millicores (1000m = 1 CPU core)
- `memory`: In Mi (Mebibytes) or Gi (Gibibytes)

### Restart Policy

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: restart-pod
spec:
  restartPolicy: Always # Always (default), OnFailure, Never
  containers:
    - name: app
      image: nginx
```

| Policy      | Behavior                                   |
| ----------- | ------------------------------------------ |
| `Always`    | Always restart container (default)         |
| `OnFailure` | Restart only if container exits with error |
| `Never`     | Never restart, even if container fails     |

---

## Hands-On Exercise 4: Pod Lifecycle States

Let's observe different pod states.

### Pending State (Image Pull)

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: pending-pod
spec:
  containers:
    - name: app
      image: nginx:this-tag-does-not-exist # Invalid tag
```

```bash
# Create the pod
kubectl apply -f pending-pod.yaml

# Check status immediately
kubectl get pods pending-pod

# Expected output:
# NAME          READY   STATUS              RESTARTS   AGE
# pending-pod   0/1     ErrImagePull        0          10s

# After a few retries:
# NAME          READY   STATUS              RESTARTS   AGE
# pending-pod   0/1     ImagePullBackOff    0          1m

# See detailed error
kubectl describe pod pending-pod

# Look for the Events section:
# Failed to pull image "nginx:this-tag-does-not-exist": rpc error: code = Unknown

# Cleanup
kubectl delete pod pending-pod
```

### CrashLoopBackOff State

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: crash-pod
spec:
  containers:
    - name: app
      image: busybox
      command: ["sh", "-c", "exit 1"] # Immediately exit with error
```

```bash
# Create the pod
kubectl apply -f crash-pod.yaml

# Watch the pod status
kubectl get pods crash-pod -w

# You'll see:
# crash-pod   0/1   ContainerCreating   0     1s
# crash-pod   0/1   Error               0     2s
# crash-pod   0/1   CrashLoopBackOff    0     3s
# crash-pod   0/1   Error               1     15s  (restart count increased)

# Kubernetes is trying to restart the container with exponential backoff

# Describe to see events
kubectl describe pod crash-pod

# Cleanup
kubectl delete pod crash-pod
```

### Completed State (Successful Job)

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: completed-pod
spec:
  restartPolicy: Never # Important for one-time tasks
  containers:
    - name: app
      image: busybox
      command: ["sh", "-c", "echo 'Task completed'; sleep 5; exit 0"]
```

```bash
# Create the pod
kubectl apply -f completed-pod.yaml

# Check status after ~10 seconds
kubectl get pod completed-pod

# Expected output:
# NAME            READY   STATUS      RESTARTS   AGE
# completed-pod   0/1     Completed   0          15s

# View logs
kubectl logs completed-pod
# Task completed

# Cleanup
kubectl delete pod completed-pod
```

---

## Init Containers

**Init containers** run before the main containers start. They're used for setup tasks.

### Use Cases

- Wait for a service to be available
- Fetch configuration from a remote source
- Initialize database schemas
- Set up file permissions

### Example: Wait for Service

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: init-demo
spec:
  initContainers:
    - name: wait-for-database
      image: busybox
      command:
        [
          "sh",
          "-c",
          "echo 'Waiting for database...'; sleep 10; echo 'Database ready'",
        ]

  containers:
    - name: app
      image: nginx
      ports:
        - containerPort: 80
```

```bash
# Create the pod
kubectl apply -f init-demo.yaml

# Watch the pod status
kubectl get pod init-demo -w

# You'll see:
# init-demo   0/1   Init:0/1   0     2s   (init container running)
# init-demo   0/1   PodInitializing   0     12s  (init completed, starting main container)
# init-demo   1/1   Running   0     15s

# View init container logs
kubectl logs init-demo -c wait-for-database
# Waiting for database...
# Database ready

# Cleanup
kubectl delete pod init-demo
```

---

## Debugging Pods

### Common Issues and Solutions

#### 1. ImagePullBackOff

**Problem**: Can't pull the container image

```bash
kubectl describe pod POD_NAME

# Look for:
# Failed to pull image "imagename:tag": ... not found
```

**Solutions**:

- Check image name and tag are correct
- Verify image exists in the registry
- Check image pull secrets (for private registries)

#### 2. CrashLoopBackOff

**Problem**: Container keeps crashing

```bash
# Check logs
kubectl logs POD_NAME

# Check previous container's logs (if it restarted)
kubectl logs POD_NAME --previous
```

**Solutions**:

- Fix application crash (check logs)
- Adjust resource limits
- Fix liveness probe configuration

#### 3. Pending State

**Problem**: Pod stuck in Pending

```bash
kubectl describe pod POD_NAME

# Look for:
# 0/1 nodes are available: 1 Insufficient cpu
```

**Solutions**:

- Not enough resources on nodes
- Node selectors don't match any node
- PersistentVolumeClaim not bound

#### 4. Container Not Starting

```bash
# Get detailed info
kubectl describe pod POD_NAME

# Check events section for errors
kubectl get events --sort-by='.lastTimestamp' | grep POD_NAME
```

---

## Challenges

### Challenge 1: Single Container Pod with Environment Variables

Create a pod that runs a busybox container with custom environment variables.

**Requirements**:

1. Create a pod named `env-test-pod`
2. Use the `busybox` image
3. Set environment variables:
   - `APP_NAME=MyApp`
   - `VERSION=1.0.0`
4. Command: `sh -c "echo APP_NAME=$APP_NAME VERSION=$VERSION; sleep 3600"`
5. Verify the environment variables are set correctly by checking logs

### Challenge 2: Multi-Container Pod with Shared Volume

Create a pod with two containers sharing a volume.

**Requirements**:

1. Pod name: `shared-volume-pod`
2. Container 1: `writer`
   - Image: `busybox`
   - Command: Write timestamp to `/data/log.txt` every 5 seconds
3. Container 2: `reader`
   - Image: `busybox`
   - Command: Read and display `/data/log.txt` every 5 seconds
4. Shared volume mounted at `/data` in both containers
5. Verify both containers can read/write the same file

### Challenge 3: Init Container Pattern

Create a pod that uses an init container to set up the environment.

**Requirements**:

1. Init container: Downloads a file or creates a configuration
2. Main container: Uses the file created by the init container
3. Use a shared volume between init and main containers
4. Verify the init container completes before the main container starts

### Challenge 4: Debug a Broken Pod

Create a pod that fails and debug it.

**Requirements**:

1. Create a pod with an intentional error (wrong image, crash, etc.)
2. Use `kubectl describe` to identify the issue
3. Use `kubectl logs` to check container output
4. Fix the issue and redeploy
5. Document the debugging steps you took

### Challenge 5: Resource Limits Test

Create a pod with resource constraints and observe behavior.

**Requirements**:

1. Create a pod with very low memory limit (e.g., 10Mi)
2. Run a memory-intensive command
3. Observe the pod being OOMKilled (Out Of Memory)
4. Check events to see the kill reason
5. Adjust limits and redeploy

---

## Solutions

<details>
<summary>Challenge 1 Solution: Single Container Pod with Environment Variables</summary>

**env-test-pod.yaml**:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: env-test-pod
spec:
  containers:
    - name: busybox
      image: busybox
      env:
        - name: APP_NAME
          value: "MyApp"
        - name: VERSION
          value: "1.0.0"
      command:
        ["sh", "-c", "echo APP_NAME=$APP_NAME VERSION=$VERSION; sleep 3600"]
```

```bash
# Create the pod
kubectl apply -f env-test-pod.yaml

# Check logs
kubectl logs env-test-pod

# Expected output:
# APP_NAME=MyApp VERSION=1.0.0

# Verify environment variables inside the container
kubectl exec env-test-pod -- env | grep -E "APP_NAME|VERSION"

# Cleanup
kubectl delete pod env-test-pod
```

</details>

<details>
<summary>Challenge 2 Solution: Multi-Container Pod with Shared Volume</summary>

**shared-volume-pod.yaml**:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: shared-volume-pod
spec:
  containers:
    - name: writer
      image: busybox
      command:
        - sh
        - -c
        - |
          while true; do
            date >> /data/log.txt
            sleep 5
          done
      volumeMounts:
        - name: shared-data
          mountPath: /data

    - name: reader
      image: busybox
      command:
        - sh
        - -c
        - |
          while true; do
            echo "=== Reading log ==="
            cat /data/log.txt 2>/dev/null || echo "File not yet created"
            sleep 5
          done
      volumeMounts:
        - name: shared-data
          mountPath: /data

  volumes:
    - name: shared-data
      emptyDir: {}
```

```bash
# Create the pod
kubectl apply -f shared-volume-pod.yaml

# Check pod status (should be 2/2 running)
kubectl get pod shared-volume-pod

# View writer logs
kubectl logs shared-volume-pod -c writer

# View reader logs
kubectl logs shared-volume-pod -c reader

# You should see timestamps being written and read

# Cleanup
kubectl delete pod shared-volume-pod
```

</details>

<details>
<summary>Challenge 3 Solution: Init Container Pattern</summary>

**init-container-pod.yaml**:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: init-setup-pod
spec:
  initContainers:
    - name: setup
      image: busybox
      command:
        - sh
        - -c
        - |
          echo "Initializing configuration..."
          echo "DATABASE_URL=postgresql://localhost:5432/mydb" > /config/app.conf
          echo "API_KEY=secret-key-12345" >> /config/app.conf
          echo "Configuration created successfully"
      volumeMounts:
        - name: config-volume
          mountPath: /config

  containers:
    - name: app
      image: busybox
      command:
        - sh
        - -c
        - |
          echo "Starting application..."
          echo "Loading configuration:"
          cat /config/app.conf
          echo "Application running..."
          sleep 3600
      volumeMounts:
        - name: config-volume
          mountPath: /config

  volumes:
    - name: config-volume
      emptyDir: {}
```

```bash
# Create the pod
kubectl apply -f init-container-pod.yaml

# Watch the pod startup
kubectl get pod init-setup-pod -w

# Check init container logs
kubectl logs init-setup-pod -c setup

# Check main container logs
kubectl logs init-setup-pod -c app

# Expected output shows config being created then loaded

# Cleanup
kubectl delete pod init-setup-pod
```

</details>

<details>
<summary>Challenge 4 Solution: Debug a Broken Pod</summary>

**broken-pod.yaml** (intentionally broken):

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: broken-pod
spec:
  containers:
    - name: app
      image: nginx:nonexistent-tag # Wrong tag
      ports:
        - containerPort: 80
```

**Debugging steps**:

```bash
# Create the broken pod
kubectl apply -f broken-pod.yaml

# Step 1: Check pod status
kubectl get pod broken-pod

# Output:
# NAME         READY   STATUS              RESTARTS   AGE
# broken-pod   0/1     ImagePullBackOff    0          30s

# Step 2: Describe the pod
kubectl describe pod broken-pod

# Look for Events section:
# Failed to pull image "nginx:nonexistent-tag": ... not found

# Step 3: Identify the issue
# The image tag doesn't exist

# Step 4: Fix the YAML
# Change image to nginx:1.25

# Step 5: Delete and recreate
kubectl delete pod broken-pod

# Create fixed version
```

**fixed-pod.yaml**:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: broken-pod
spec:
  containers:
    - name: app
      image: nginx:1.25 # Fixed
      ports:
        - containerPort: 80
```

```bash
kubectl apply -f fixed-pod.yaml

# Verify it's running
kubectl get pod broken-pod
# NAME         READY   STATUS    RESTARTS   AGE
# broken-pod   1/1     Running   0          10s

# Cleanup
kubectl delete pod broken-pod
```

</details>

<details>
<summary>Challenge 5 Solution: Resource Limits Test</summary>

**memory-limit-pod.yaml**:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: memory-limit-pod
spec:
  containers:
    - name: memory-hog
      image: polinux/stress
      resources:
        limits:
          memory: "10Mi" # Very low limit
        requests:
          memory: "5Mi"
      command: ["stress"]
      args: ["--vm", "1", "--vm-bytes", "50M", "--vm-hang", "1"]
      # This tries to allocate 50M but limit is 10Mi
```

```bash
# Create the pod
kubectl apply -f memory-limit-pod.yaml

# Watch the pod
kubectl get pod memory-limit-pod -w

# You'll see:
# memory-limit-pod   0/1   OOMKilled   0     10s
# memory-limit-pod   0/1   CrashLoopBackOff   1     20s

# Describe to see the OOMKilled event
kubectl describe pod memory-limit-pod

# Look for:
# Last State:     Terminated
#   Reason:       OOMKilled
#   Exit Code:    137

# Fix by increasing memory limit
```

**memory-limit-pod-fixed.yaml**:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: memory-limit-pod-fixed
spec:
  containers:
    - name: memory-hog
      image: polinux/stress
      resources:
        limits:
          memory: "100Mi" # Increased
        requests:
          memory: "50Mi"
      command: ["stress"]
      args: ["--vm", "1", "--vm-bytes", "50M", "--vm-hang", "1"]
```

```bash
# Apply fixed version
kubectl apply -f memory-limit-pod-fixed.yaml

# Verify it runs successfully
kubectl get pod memory-limit-pod-fixed

# Cleanup
kubectl delete pod memory-limit-pod memory-limit-pod-fixed --ignore-not-found
```

</details>

---

## Best Practices

1. **Use Deployments, Not Bare Pods**
   - Pods don't self-heal if deleted
   - Always use Deployments, StatefulSets, or Jobs
   - Create pods directly only for debugging/testing

2. **One Primary Process Per Container**
   - Don't run multiple unrelated processes in one container
   - Use multi-container pods if you need multiple processes
   - Exception: Helper processes (log rotation, etc.)

3. **Design Pods for Horizontal Scaling**
   - Assume your pod can be killed and recreated at any time
   - Don't store critical data in the container filesystem
   - Use external storage for persistence

4. **Use Init Containers for Setup**
   - Separate initialization logic from main application
   - Makes startup sequence clearer
   - Easier to debug and maintain

5. **Set Resource Requests and Limits**

   ```yaml
   resources:
     requests:
       memory: "64Mi"
       cpu: "250m"
     limits:
       memory: "128Mi"
       cpu: "500m"
   ```

   - Prevents pods from consuming all node resources
   - Helps scheduler make better placement decisions

6. **Use Meaningful Labels**

   ```yaml
   metadata:
     labels:
       app: myapp
       version: v1.2.0
       environment: production
       tier: backend
   ```

7. **Implement Proper Health Checks** (covered in Lesson 17)
   - Liveness probes: Detect broken containers
   - Readiness probes: Control traffic flow
   - Startup probes: Handle slow-starting apps

---

## Key Takeaways

1. **Pods are the atomic unit** in Kubernetes - you never deploy containers directly
2. **Most pods contain one container**, but multi-container patterns solve specific problems
3. **Containers in a pod share network and storage**, enabling tight coupling
4. **Pod lifecycle phases**: Pending → Running → Succeeded/Failed
5. **Init containers run before main containers**, useful for setup tasks
6. **Always use higher-level resources** (Deployments) instead of bare pods in production
7. **Resource requests and limits** prevent resource starvation
8. **Multi-container patterns**: Sidecar (helper), Ambassador (proxy), Adapter (transform)

---

## Next Steps

In [Lesson 03: Kubernetes Architecture](lesson-03-kubernetes-architecture.md), you'll learn:

- How the Kubernetes control plane works
- Node components and their responsibilities
- How components communicate
- The declarative model: desired vs actual state
- How Kubernetes achieves self-healing

---

## Questions to Ponder

1. Why would you ever use a multi-container pod instead of separate single-container pods?
2. What happens to a pod's IP address when it restarts?
3. If init containers fail, will the main containers start?
4. How does Kubernetes decide which node to place a pod on?
5. Can containers in different pods communicate via `localhost`? Why or why not?

---

## Cleanup

```bash
# Delete all pods we created
kubectl delete pod my-first-pod web-logger-pod --ignore-not-found

# Verify all pods are deleted
kubectl get pods
```

---

**Congratulations!** You now understand pods deeply. You know when to use single vs multi-container pods, how to debug pod issues, and the various lifecycle states. In the next lesson, we'll explore the Kubernetes architecture to understand how all these components work together behind the scenes.
