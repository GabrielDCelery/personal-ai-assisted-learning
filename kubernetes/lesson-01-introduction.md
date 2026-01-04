# Lesson 01: Introduction to Kubernetes

## Overview

In this lesson, you'll understand **why** Kubernetes exists and **what problems** it solves. You'll set up a local Kubernetes cluster, deploy your first application, and compare the Kubernetes approach to what you already know from Docker.

By the end of this lesson, you'll run the same application using both Docker and Kubernetes, understanding the key differences and when each tool is appropriate.

---

## The Problem: Docker in Production

Imagine you've built a fantastic web application using Docker. You've containerized it, tested it locally, and it works perfectly. Now you need to deploy it to production.

### Scenario: Your Application Needs to Scale

Your app becomes popular. You need to:

- Run 10 copies of your app container across 5 servers for high availability
- Automatically restart containers if they crash
- Load balance traffic across all 10 copies
- Deploy updates without downtime (rolling updates)
- Automatically scale to 20 copies during peak traffic
- Monitor which containers are healthy and route traffic only to them

### The Docker-Only Solution

```bash
# Server 1
docker run -d --name app-1 --restart always -p 8080:8080 myapp:v1
docker run -d --name app-2 --restart always -p 8081:8080 myapp:v1

# Server 2
docker run -d --name app-3 --restart always -p 8080:8080 myapp:v1
docker run -d --name app-4 --restart always -p 8081:8080 myapp:v1

# ... repeat for servers 3, 4, 5
# Then manually configure a load balancer to point to all 10 containers
# Write custom scripts to check health and restart failed containers
# Write more scripts to handle rolling updates
# Hope nothing breaks at 3 AM
```

**This is painful.** You need:

- SSH access to every server
- Custom scripts for health checks, scaling, updates
- Manual load balancer configuration
- Custom monitoring and alerting
- A plan for handling server failures

**This is where Kubernetes comes in.**

---

## What is Kubernetes?

**Kubernetes** (K8s) is a **container orchestration platform**. It automates the deployment, scaling, and management of containerized applications across clusters of machines.

### Key Concept: Orchestration

```
Docker:     🎻 (One musician playing a violin)
Kubernetes: 🎼 (A conductor coordinating an entire orchestra)
```

- **Docker** gives you containers (the musicians/instruments)
- **Kubernetes** orchestrates how containers run across many machines (the conductor)

### What Kubernetes Does

1. **Automated Scheduling**: Decides which machine runs which container
2. **Self-Healing**: Automatically restarts failed containers, replaces them, kills unhealthy ones
3. **Horizontal Scaling**: Adds or removes container replicas based on load
4. **Service Discovery**: Containers find each other automatically
5. **Load Balancing**: Distributes traffic across container replicas
6. **Rolling Updates**: Deploys new versions without downtime
7. **Secret Management**: Manages sensitive configuration securely
8. **Storage Orchestration**: Automatically mounts storage systems

---

## Kubernetes vs Docker: The Relationship

This is the most common confusion. Let's clarify:

### Kubernetes USES Docker (or other container runtimes)

```
┌─────────────────────────────────────────┐
│           Kubernetes                     │
│  (Orchestration & Management)            │
│                                          │
│  ┌──────────┐  ┌──────────┐            │
│  │  Pod 1   │  │  Pod 2   │            │
│  │┌────────┐│  │┌────────┐│            │
│  ││Container││  ││Container││            │
│  ││ (Docker)││  ││(Docker) ││            │
│  │└────────┘│  │└────────┘│            │
│  └──────────┘  └──────────┘            │
└─────────────────────────────────────────┘
```

**Docker**: Builds and runs containers
**Kubernetes**: Decides **where** and **how many** containers run, manages their lifecycle

### Comparison Table

| Feature            | Docker                               | Kubernetes                                  |
| ------------------ | ------------------------------------ | ------------------------------------------- |
| **Purpose**        | Run containers on a single machine   | Orchestrate containers across many machines |
| **Scope**          | Single host                          | Cluster of machines                         |
| **Scaling**        | Manual (`docker run` multiple times) | Declarative (`replicas: 10`)                |
| **Health Checks**  | Basic restart policies               | Sophisticated liveness/readiness probes     |
| **Load Balancing** | External tool needed                 | Built-in (Services)                         |
| **Updates**        | Manual stop/start                    | Rolling updates, rollbacks                  |
| **Configuration**  | Environment variables, compose files | ConfigMaps, Secrets, manifests              |
| **Networking**     | Docker networks                      | Cluster networking, Services, Ingress       |
| **Storage**        | Volumes                              | Persistent Volumes, Storage Classes         |
| **Best For**       | Development, simple deployments      | Production, complex distributed systems     |

---

## Core Kubernetes Concepts (Quick Introduction)

We'll dive deep into these in later lessons. For now, here's a mental model:

### 1. Cluster

A set of machines (nodes) that run your containers.

```
┌─────────────────────────────────────────────────────┐
│                 Kubernetes Cluster                   │
│                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────┐ │
│  │   Master     │  │   Worker     │  │  Worker   │ │
│  │   Node       │  │   Node 1     │  │  Node 2   │ │
│  │ (Control     │  │ (Runs Pods)  │  │ (Runs     │ │
│  │  Plane)      │  │              │  │  Pods)    │ │
│  └──────────────┘  └──────────────┘  └───────────┘ │
└─────────────────────────────────────────────────────┘
```

### 2. Node

A machine (physical or virtual) in the cluster. Can be a master (control plane) or worker.

### 3. Pod

The smallest deployable unit. Usually contains one container (but can have multiple).

```
Pod = A wrapper around your container(s)
```

### 4. Deployment

Manages multiple replicas of your pods. Handles scaling and updates.

```
Deployment → Creates ReplicaSet → Creates Pods → Run Containers
```

### 5. Service

A stable network endpoint to access your pods (which have dynamic IPs).

```
Service: "Talk to my-app on port 80"
         ↓ (routes traffic to)
      Pod 1, Pod 2, Pod 3, ... (whichever are healthy)
```

---

## Setting Up Your Local Kubernetes Environment

You need a Kubernetes cluster to practice. For local development, we'll use **minikube**.

### Option 1: Minikube (Recommended)

Minikube runs a single-node Kubernetes cluster inside a VM on your laptop.

#### Installation (Linux)

```bash
# Download minikube
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64

# Install
sudo install minikube-linux-amd64 /usr/local/bin/minikube

# Verify
minikube version
```

#### Installation (macOS)

```bash
# Using Homebrew
brew install minikube

# Verify
minikube version
```

#### Installation (Windows)

```powershell
# Using Chocolatey
choco install minikube

# Or download installer from:
# https://minikube.sigs.k8s.io/docs/start/
```

### Option 2: kind (Kubernetes in Docker)

Runs Kubernetes inside Docker containers (lightweight).

```bash
# Install kind
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

# Create cluster
kind create cluster --name my-cluster
```

### Option 3: Docker Desktop

Docker Desktop (Mac/Windows) includes Kubernetes. Enable it in Settings → Kubernetes → Enable Kubernetes.

---

## Hands-On Exercise 1: Start Your First Cluster

Let's start a local Kubernetes cluster.

```bash
# Start minikube
minikube start

# Expected output:
# 😄  minikube v1.32.0 on Ubuntu 22.04
# ✨  Automatically selected the docker driver
# 👍  Starting control plane node minikube in cluster minikube
# 🚜  Pulling base image ...
# 🔥  Creating docker container (CPUs=2, Memory=2200MB) ...
# 🐳  Preparing Kubernetes v1.28.3 on Docker 24.0.7 ...
# 🔎  Verifying Kubernetes components...
# 🌟  Enabled addons: storage-provisioner, default-storageclass
# 🏄  Done! kubectl is now configured to use "minikube" cluster
```

### Verify Your Cluster

```bash
# Check cluster info
kubectl cluster-info

# Expected output:
# Kubernetes control plane is running at https://192.168.49.2:8443
# CoreDNS is running at https://192.168.49.2:8443/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

# List nodes
kubectl get nodes

# Expected output:
# NAME       STATUS   ROLES           AGE   VERSION
# minikube   Ready    control-plane   1m    v1.28.3
```

**Congratulations!** You now have a Kubernetes cluster running on your machine.

---

## Hands-On Exercise 2: Your First Deployment (Imperative Way)

Let's deploy an nginx web server, similar to what you did with Docker.

### Docker Way (What You Know)

```bash
docker run -d --name my-nginx -p 8080:80 nginx
```

### Kubernetes Way

```bash
# Create a deployment (this creates pods)
kubectl create deployment my-nginx --image=nginx

# Verify the deployment
kubectl get deployments

# Expected output:
# NAME       READY   UP-TO-DATE   AVAILABLE   AGE
# my-nginx   1/1     1            1           10s

# See the pod that was created
kubectl get pods

# Expected output:
# NAME                        READY   STATUS    RESTARTS   AGE
# my-nginx-5b7d4c8d9b-x7k2l   1/1     Running   0          20s
```

### Expose the Deployment (Make it Accessible)

```bash
# Create a service to expose the deployment
kubectl expose deployment my-nginx --port=80 --type=NodePort

# Get the service details
kubectl get services

# Expected output:
# NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
# my-nginx     NodePort    10.96.100.123   <none>        80:32000/TCP   5s

# Access the nginx server (minikube provides a URL)
minikube service my-nginx --url

# Expected output: http://192.168.49.2:32000
# Open this URL in your browser to see nginx
```

### Understanding What Happened

1. **Deployment**: Created a blueprint for running nginx
2. **Pod**: Kubernetes created a pod containing the nginx container
3. **Service**: Created a stable endpoint to access the pod

---

## Hands-On Exercise 3: Declarative Approach (The Kubernetes Way)

The `kubectl create` commands are **imperative** (you tell K8s what to do step-by-step).
The **declarative** way is to define the desired state in YAML files.

### Create a Deployment YAML

```bash
# Create a file named nginx-deployment.yaml
```

**nginx-deployment.yaml**:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3 # Run 3 copies
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.25
          ports:
            - containerPort: 80
```

### Apply the Deployment

```bash
# Apply the YAML file
kubectl apply -f nginx-deployment.yaml

# Verify
kubectl get deployments
# NAME               READY   UP-TO-DATE   AVAILABLE   AGE
# nginx-deployment   3/3     3            3           10s

kubectl get pods
# NAME                                READY   STATUS    RESTARTS   AGE
# nginx-deployment-7d64f5b8c9-abc12   1/1     Running   0          15s
# nginx-deployment-7d64f5b8c9-def34   1/1     Running   0          15s
# nginx-deployment-7d64f5b8c9-ghi56   1/1     Running   0          15s
```

**You now have 3 nginx pods running!**

### Create a Service YAML

**nginx-service.yaml**:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  selector:
    app: nginx
  ports:
    - protocol: TCP
      port: 80
      targetPort: 80
  type: NodePort
```

### Apply the Service

```bash
kubectl apply -f nginx-service.yaml

kubectl get services
# NAME            TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
# nginx-service   NodePort    10.96.200.50    <none>        80:30123/TCP   5s

# Access the service
minikube service nginx-service --url
```

---

## Hands-On Exercise 4: Scaling Your Application

One of Kubernetes' superpowers is easy scaling.

### Scale Up

```bash
# Scale to 5 replicas
kubectl scale deployment nginx-deployment --replicas=5

# Watch pods being created
kubectl get pods -w
# (Press Ctrl+C to stop watching)

# Verify
kubectl get pods
# You should see 5 pods now
```

### Scale Down

```bash
# Scale back to 2 replicas
kubectl scale deployment nginx-deployment --replicas=2

kubectl get pods
# You should see 2 pods (Kubernetes terminated 3)
```

### Scaling with YAML (Declarative)

Edit `nginx-deployment.yaml` and change `replicas: 3` to `replicas: 10`:

```yaml
spec:
  replicas: 10
```

```bash
# Apply the change
kubectl apply -f nginx-deployment.yaml

# Kubernetes will create 8 more pods to reach 10
kubectl get pods
```

**This is the Kubernetes philosophy**: You declare the desired state, Kubernetes makes it happen.

---

## Hands-On Exercise 5: Self-Healing in Action

Kubernetes automatically restarts failed containers.

### Delete a Pod Manually

```bash
# List pods and copy one pod name
kubectl get pods

# Delete one pod (replace with actual pod name)
kubectl delete pod nginx-deployment-7d64f5b8c9-abc12

# Immediately check pods again
kubectl get pods
```

**What happened?**
Kubernetes immediately created a new pod to replace the deleted one. The deployment says "I need 10 replicas," so Kubernetes maintains that state.

### Simulate a Container Crash

```bash
# Exec into a pod (replace with actual pod name)
kubectl exec -it nginx-deployment-7d64f5b8c9-def34 -- /bin/bash

# Inside the container, kill nginx
root@nginx-deployment:/# kill 1

# Exit (the container will crash)
# Check pods
kubectl get pods

# You'll see the pod restarting:
# NAME                                READY   STATUS    RESTARTS   AGE
# nginx-deployment-7d64f5b8c9-def34   1/1     Running   1          5m
#                                                       ↑ (restart count increased)
```

Kubernetes detected the crash and restarted the container automatically.

---

## Common kubectl Commands Reference

### Cluster Management

| Command                | Description                       |
| ---------------------- | --------------------------------- |
| `kubectl cluster-info` | Display cluster information       |
| `kubectl get nodes`    | List all nodes in the cluster     |
| `kubectl version`      | Show kubectl and cluster versions |

### Working with Deployments

| Command                                        | Description                      |
| ---------------------------------------------- | -------------------------------- |
| `kubectl create deployment NAME --image=IMAGE` | Create a deployment              |
| `kubectl get deployments`                      | List all deployments             |
| `kubectl describe deployment NAME`             | Detailed info about a deployment |
| `kubectl delete deployment NAME`               | Delete a deployment              |
| `kubectl scale deployment NAME --replicas=N`   | Scale deployment to N replicas   |

### Working with Pods

| Command                                  | Description                            |
| ---------------------------------------- | -------------------------------------- |
| `kubectl get pods`                       | List all pods                          |
| `kubectl get pods -o wide`               | List pods with more details (IP, node) |
| `kubectl describe pod NAME`              | Detailed info about a pod              |
| `kubectl logs POD_NAME`                  | View pod logs                          |
| `kubectl logs POD_NAME -f`               | Stream pod logs (follow)               |
| `kubectl exec -it POD_NAME -- /bin/bash` | Execute shell in a pod                 |
| `kubectl delete pod NAME`                | Delete a pod                           |

### Working with Services

| Command                                    | Description                       |
| ------------------------------------------ | --------------------------------- |
| `kubectl get services`                     | List all services                 |
| `kubectl describe service NAME`            | Detailed info about a service     |
| `kubectl expose deployment NAME --port=80` | Create a service for a deployment |
| `kubectl delete service NAME`              | Delete a service                  |

### Declarative Management

| Command                       | Description                            |
| ----------------------------- | -------------------------------------- |
| `kubectl apply -f FILE.yaml`  | Create/update resources from YAML file |
| `kubectl delete -f FILE.yaml` | Delete resources defined in YAML file  |
| `kubectl get -f FILE.yaml`    | Get status of resources in YAML file   |

### Debugging

| Command                            | Description            |
| ---------------------------------- | ---------------------- |
| `kubectl logs POD_NAME`            | View logs from a pod   |
| `kubectl describe pod POD_NAME`    | See events and details |
| `kubectl get events`               | See cluster events     |
| `kubectl exec POD_NAME -- COMMAND` | Run command in a pod   |

---

## Kubernetes vs Docker Compose

You might be thinking: "Docker Compose also manages multiple containers. How is this different?"

| Feature              | Docker Compose                        | Kubernetes                            |
| -------------------- | ------------------------------------- | ------------------------------------- |
| **Scope**            | Single machine                        | Multiple machines (cluster)           |
| **Scaling**          | Limited to one host                   | Scales across cluster                 |
| **Self-Healing**     | Basic restart policies                | Automatic rescheduling, health checks |
| **Load Balancing**   | Manual setup                          | Built-in (Services)                   |
| **Production Ready** | Good for development                  | Designed for production               |
| **Rolling Updates**  | Manual                                | Built-in with rollback                |
| **Complexity**       | Simple                                | More complex, more powerful           |
| **Use Case**         | Local development, simple deployments | Production, distributed systems       |

**When to use Docker Compose**: Local development, simple applications
**When to use Kubernetes**: Production deployments, microservices, applications requiring high availability

---

## Challenges

### Challenge 1: Deploy a Custom Application

Deploy the Flask application you created in the Docker course to Kubernetes.

**Requirements:**

1. Build your Docker image (from Docker lesson 04 challenge-1)
2. Load the image into minikube: `minikube image load IMAGE_NAME`
3. Create a Deployment with 3 replicas
4. Expose it via a Service
5. Access the application using `minikube service`
6. Verify all 3 pods are running

### Challenge 2: Scale Based on "Load"

Simulate a traffic spike and scale your application.

**Requirements:**

1. Start with 2 replicas of nginx
2. "Simulate traffic" by scaling to 10 replicas
3. Verify all 10 pods are running and healthy
4. Scale back down to 2 replicas
5. Confirm pods are terminated gracefully

### Challenge 3: Self-Healing Test

Test Kubernetes' self-healing capabilities.

**Requirements:**

1. Deploy nginx with 5 replicas
2. Manually delete 2 pods
3. Observe Kubernetes recreating them
4. Document the time it takes for the pods to be replaced
5. Explain what component of Kubernetes is responsible for this behavior

### Challenge 4: Multi-Container Pod

Create a pod with two containers that share the same network namespace.

**Requirements:**

1. Create a YAML file defining a pod with:
   - Container 1: nginx (serves web traffic)
   - Container 2: busybox (runs `sleep 3600`)
2. Apply the pod
3. Exec into the busybox container
4. Use `wget localhost:80` to access nginx (they share localhost!)
5. Explain why this works (network namespace sharing)

---

## Solutions

<details>
<summary>Challenge 1 Solution: Deploy Custom Application</summary>

Assuming you have a Flask app image named `flask-api:latest`:

```bash
# Build the image (from Docker lesson 04 challenge-1 directory)
docker build -t flask-api:latest .

# Load image into minikube
minikube image load flask-api:latest

# Create deployment YAML
```

**flask-deployment.yaml**:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: flask-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: flask
  template:
    metadata:
      labels:
        app: flask
    spec:
      containers:
        - name: flask
          image: flask-api:latest
          imagePullPolicy: Never # Use local image
          ports:
            - containerPort: 5000
```

**flask-service.yaml**:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: flask-service
spec:
  selector:
    app: flask
  ports:
    - protocol: TCP
      port: 5000
      targetPort: 5000
  type: NodePort
```

```bash
# Apply the manifests
kubectl apply -f flask-deployment.yaml
kubectl apply -f flask-service.yaml

# Verify
kubectl get pods -l app=flask
kubectl get services flask-service

# Access the application
minikube service flask-service --url

# Test (assuming it has a /health endpoint)
curl $(minikube service flask-service --url)/health
```

</details>

<details>
<summary>Challenge 2 Solution: Scale Based on Load</summary>

```bash
# Create initial deployment
kubectl create deployment nginx-load-test --image=nginx --replicas=2

# Verify initial state
kubectl get pods -l app=nginx-load-test
# Should show 2 pods

# "Traffic spike" - scale to 10
kubectl scale deployment nginx-load-test --replicas=10

# Watch scaling in real-time
kubectl get pods -l app=nginx-load-test -w
# Press Ctrl+C when all 10 are Running

# Verify
kubectl get deployments nginx-load-test
# READY should show 10/10

# "Traffic decreased" - scale down
kubectl scale deployment nginx-load-test --replicas=2

# Watch scale-down
kubectl get pods -l app=nginx-load-test -w

# Cleanup
kubectl delete deployment nginx-load-test
```

**Observations:**

- Kubernetes creates new pods within seconds
- Scaling up and down is seamless
- No manual intervention needed on individual servers

</details>

<details>
<summary>Challenge 3 Solution: Self-Healing Test</summary>

```bash
# Create deployment
kubectl create deployment nginx-healing --image=nginx --replicas=5

# List pods and note their names
kubectl get pods -l app=nginx-healing

# Delete 2 pods (replace POD_NAME_1 and POD_NAME_2 with actual names)
time kubectl delete pod POD_NAME_1 POD_NAME_2

# Immediately check status
kubectl get pods -l app=nginx-healing

# Expected output: You'll see Terminating pods and new pods being created
```

**Expected behavior:**

- Within 1-2 seconds, Kubernetes detects pods are missing
- New pods are scheduled immediately
- Total time to full recovery: ~10-30 seconds (depending on image pull time)

**Explanation:**
The **ReplicaSet controller** (created by the Deployment) continuously monitors the desired vs actual state. When it detects fewer than 5 pods, it immediately creates new ones.

```bash
# Verify the ReplicaSet
kubectl get replicasets -l app=nginx-healing

# See events
kubectl get events --sort-by='.lastTimestamp' | grep nginx-healing

# Cleanup
kubectl delete deployment nginx-healing
```

</details>

<details>
<summary>Challenge 4 Solution: Multi-Container Pod</summary>

**multi-container-pod.yaml**:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: multi-container-pod
spec:
  containers:
    - name: nginx
      image: nginx
      ports:
        - containerPort: 80
    - name: busybox
      image: busybox
      command: ["sleep", "3600"]
```

```bash
# Create the pod
kubectl apply -f multi-container-pod.yaml

# Verify both containers are running
kubectl get pod multi-container-pod
# READY should show 2/2

# Exec into the busybox container
kubectl exec -it multi-container-pod -c busybox -- /bin/sh

# Inside the busybox container, access nginx via localhost
/ # wget -O- localhost:80
# You'll see the nginx welcome page HTML!

# Why does this work?
# Both containers share the same network namespace
# They both see "localhost" as the same network interface
# nginx is listening on port 80 in this shared namespace
# busybox can access it via localhost:80

# Exit
/ # exit

# Check logs from nginx container
kubectl logs multi-container-pod -c nginx
# You'll see the GET request from busybox

# Cleanup
kubectl delete pod multi-container-pod
```

**Explanation:**

- Containers in the same pod share the network namespace
- They communicate via `localhost`
- They also share the same storage volumes (we'll cover this later)
- Common pattern: Main app container + sidecar helper container (logging, monitoring, etc.)

</details>

---

## Best Practices (Introduction)

1. **Use Declarative YAML Files**
   - Always prefer `kubectl apply -f` over `kubectl create` commands
   - YAML files are version-controllable, reproducible, and self-documenting

2. **Label Everything**
   - Labels help you organize and select resources
   - Use meaningful labels: `app: nginx`, `env: production`, `version: v1.2.0`

3. **Never Use Latest Tag in Production**

   ```yaml
   # Bad
   image: nginx:latest

   # Good
   image: nginx:1.25.3
   ```

4. **Start with Small Replicas**
   - Don't deploy 100 replicas on day one
   - Start with 2-3, monitor, then scale as needed

5. **Use Services for Communication**
   - Never hardcode pod IPs (they change!)
   - Always access pods through Services (stable DNS names)

---

## Key Takeaways

1. **Kubernetes orchestrates containers** across multiple machines, while Docker runs containers on a single machine
2. **Declarative management** (YAML files) is the Kubernetes way: you define desired state, K8s makes it happen
3. **Pods** are the smallest unit, but you rarely create them directly
4. **Deployments** manage replicas of your pods and handle scaling/updates
5. **Services** provide stable networking to access your pods
6. **Self-healing** is automatic: Kubernetes recreates failed containers
7. **Scaling** is trivial: change `replicas` value, apply YAML
8. **kubectl** is your primary tool for interacting with Kubernetes

---

## Next Steps

In [Lesson 02: Understanding Pods](lesson-02-understanding-pods.md), you'll dive deep into pods:

- Pod lifecycle and states
- Multi-container pod patterns
- Init containers
- Pod design best practices
- Resource requests and limits

---

## Questions to Ponder

1. If Kubernetes automatically restarts failed containers, why would you ever need to SSH into a server?
2. You have a web app with a database. Should they be in the same pod or separate pods? Why?
3. How does Kubernetes know which pods belong to which Deployment?
4. What happens if you scale a deployment to 0 replicas? Why might you do this?
5. If you have 3 nginx pods, and you access the Service 10 times, which pod serves each request?

---

## Cleanup

When you're done experimenting:

```bash
# Delete all resources we created
kubectl delete deployment nginx-deployment my-nginx
kubectl delete service nginx-service my-nginx

# Stop minikube (optional)
minikube stop

# Delete minikube cluster (if you want to start fresh later)
minikube delete
```

---

**Congratulations!** You've taken your first steps into Kubernetes. You've deployed applications, scaled them, and witnessed self-healing in action. In the next lesson, we'll explore pods in depth and understand why they're designed the way they are.
