# Lesson 03: Kubernetes Architecture

## Overview

In Lessons 01 and 02, you deployed applications to Kubernetes and worked with pods. But what's actually happening behind the scenes? When you run `kubectl apply`, which components handle your request? How does Kubernetes keep your desired state running?

In this lesson, you'll understand the Kubernetes architecture: the control plane components that manage the cluster, the node components that run your workloads, and how they work together to maintain your applications.

By the end of this lesson, you'll be able to explain how Kubernetes achieves self-healing, how the scheduler decides where to place pods, and how the declarative model works.

---

## The Big Picture

Kubernetes follows a **master-worker architecture** (also called control plane and worker nodes).

```
┌─────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                        │
│                                                              │
│  ┌────────────────────────────────────────────────────┐    │
│  │          Control Plane (Master Node)                │    │
│  │  ┌──────────┐  ┌──────────┐  ┌────────────────┐   │    │
│  │  │   API    │  │   etcd   │  │   Scheduler    │   │    │
│  │  │  Server  │  │          │  │                │   │    │
│  │  └──────────┘  └──────────┘  └────────────────┘   │    │
│  │  ┌──────────────────────────────────────────────┐ │    │
│  │  │       Controller Manager                     │ │    │
│  │  └──────────────────────────────────────────────┘ │    │
│  └────────────────────────────────────────────────────┘    │
│                            │                                │
│                            │ (manages)                      │
│                            ↓                                │
│  ┌─────────────────┐  ┌─────────────────┐  ┌────────────┐ │
│  │  Worker Node 1  │  │  Worker Node 2  │  │  Worker    │ │
│  │  ┌───────────┐  │  │  ┌───────────┐  │  │  Node 3    │ │
│  │  │  kubelet  │  │  │  │  kubelet  │  │  │  ┌──────┐  │ │
│  │  ├───────────┤  │  │  ├───────────┤  │  │  │Pods  │  │ │
│  │  │kube-proxy │  │  │  │kube-proxy │  │  │  └──────┘  │ │
│  │  ├───────────┤  │  │  ├───────────┤  │  │            │ │
│  │  │   Pods    │  │  │  │   Pods    │  │  │            │ │
│  │  └───────────┘  │  │  └───────────┘  │  │            │ │
│  └─────────────────┘  └─────────────────┘  └────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

**Key concepts**:

- **Control Plane**: The brain that makes decisions about the cluster
- **Worker Nodes**: The muscle that runs your actual workloads
- **Communication**: Control plane tells workers what to do; workers report their status

---

## Control Plane Components

The control plane manages the entire cluster. It makes global decisions about scheduling, responds to cluster events, and maintains the desired state.

### 1. API Server (kube-apiserver)

**The front door to Kubernetes.**

```
You → kubectl → API Server → etcd/scheduler/controllers
```

**What it does**:

- Exposes the Kubernetes API (RESTful)
- Validates and processes all requests
- Only component that talks directly to etcd
- Authentication and authorization gateway

**Example flow**:

```bash
kubectl create deployment nginx --image=nginx

# What happens:
# 1. kubectl sends HTTP POST to API Server
# 2. API Server authenticates you
# 3. API Server validates the request
# 4. API Server writes to etcd
# 5. API Server notifies controllers about new deployment
```

**Key characteristics**:

- Stateless (all state stored in etcd)
- Horizontally scalable (can run multiple replicas)
- Every operation goes through the API server

---

### 2. etcd

**The cluster's database and source of truth.**

```
┌──────────────────────────┐
│         etcd             │
│  (Key-Value Store)       │
│                          │
│  deployments/nginx: {...}│
│  pods/nginx-abc: {...}   │
│  services/nginx: {...}   │
└──────────────────────────┘
```

**What it does**:

- Stores all cluster data (configuration, state, metadata)
- Distributed and consistent (uses Raft consensus algorithm)
- Strongly consistent key-value store
- Watches for changes and notifies API server

**What's stored in etcd**:

- Cluster configuration
- Current state of all resources (pods, deployments, services, etc.)
- Secrets and ConfigMaps
- Node information
- Resource quotas

**Critical point**: If you lose etcd data, you lose your entire cluster state. Always back up etcd in production.

---

### 3. Scheduler (kube-scheduler)

**The matchmaker: finds the right node for each pod.**

```
Scheduler receives: "Need to run pod X"

Scheduler evaluates:
├─ Node 1: 2GB available, CPU 50% ✓
├─ Node 2: 500MB available, CPU 90% ✗ (not enough resources)
└─ Node 3: 4GB available, CPU 20% ✓✓ (best fit!)

Decision: Schedule pod X on Node 3
```

**What it does**:

- Watches for newly created pods with no assigned node
- Selects the best node for each pod
- Updates the pod definition with the node assignment
- Does NOT actually run the pod (kubelet does that)

**Scheduling algorithm considers**:

- Resource requests (CPU, memory)
- Node affinity/anti-affinity rules
- Taints and tolerations
- Data locality
- Inter-pod affinity
- Hardware constraints

**Example**:

```yaml
# Pod with resource requests
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  containers:
    - name: app
      image: nginx
      resources:
        requests:
          memory: "128Mi"
          cpu: "500m"

# Scheduler will only place this on nodes with at least 128Mi memory and 500m CPU available
```

---

### 4. Controller Manager (kube-controller-manager)

**The autopilot: ensures desired state matches actual state.**

The controller manager runs multiple controllers, each watching a specific resource type.

```
┌────────────────────────────────────┐
│     Controller Manager             │
│                                    │
│  ┌──────────────────────────────┐ │
│  │  Deployment Controller       │ │  Watches: Deployments
│  │  (Manages ReplicaSets)       │ │  Ensures: Correct ReplicaSets exist
│  └──────────────────────────────┘ │
│  ┌──────────────────────────────┐ │
│  │  ReplicaSet Controller       │ │  Watches: ReplicaSets
│  │  (Manages Pods)              │ │  Ensures: Correct number of pods
│  └──────────────────────────────┘ │
│  ┌──────────────────────────────┐ │
│  │  Node Controller             │ │  Watches: Nodes
│  │  (Monitors node health)      │ │  Evicts pods from failed nodes
│  └──────────────────────────────┘ │
│  ┌──────────────────────────────┐ │
│  │  Service Controller          │ │  Watches: Services
│  │  (Manages endpoints)         │ │  Updates endpoint lists
│  └──────────────────────────────┘ │
└────────────────────────────────────┘
```

**Key controllers**:

| Controller         | Responsibility                                      |
| ------------------ | --------------------------------------------------- |
| **Deployment**     | Manages ReplicaSets for deployments                 |
| **ReplicaSet**     | Ensures correct number of pod replicas are running  |
| **Node**           | Monitors node health, evicts pods from failed nodes |
| **Service**        | Updates endpoints when pods come and go             |
| **Endpoint**       | Populates Endpoint objects (links Services to Pods) |
| **Namespace**      | Deletes all resources when a namespace is deleted   |
| **ServiceAccount** | Creates default service accounts for namespaces     |

**How controllers work (control loop)**:

```
Loop forever:
  1. Watch for changes to resources
  2. Read desired state (from etcd via API server)
  3. Read actual state (from etcd via API server)
  4. If desired != actual:
     Take action to reconcile (create/update/delete resources)
  5. Sleep briefly
  6. Repeat
```

**Example: ReplicaSet Controller**:

```
Desired state: 3 nginx pods
Actual state: 2 nginx pods running

Action: Create 1 new pod
Result: Desired state achieved
```

---

### 5. Cloud Controller Manager (cloud-controller-manager)

**The cloud integrator (optional, only for cloud providers).**

**What it does**:

- Integrates Kubernetes with cloud provider APIs (AWS, GCP, Azure)
- Manages cloud-specific resources
- Separates cloud-specific code from core Kubernetes

**Responsibilities**:

- **Node Controller**: Checks if nodes exist in cloud provider
- **Route Controller**: Sets up network routes in cloud infrastructure
- **Service Controller**: Creates cloud load balancers for LoadBalancer services
- **Volume Controller**: Creates and attaches cloud volumes

---

## Worker Node Components

Worker nodes run your actual workloads (pods). Each node has three main components.

### 1. kubelet

**The node agent: ensures containers are running.**

```
┌────────────────────────────────────┐
│           Worker Node              │
│                                    │
│  ┌──────────────────────────────┐ │
│  │         kubelet              │ │
│  │  ┌─────────────────────────┐ │ │
│  │  │ 1. Watch API Server     │ │ │
│  │  │ 2. Get pod specs        │ │ │
│  │  │ 3. Tell container       │ │ │
│  │  │    runtime to start     │ │ │
│  │  │    containers           │ │ │
│  │  │ 4. Monitor pods         │ │ │
│  │  │ 5. Report status back   │ │ │
│  │  └─────────────────────────┘ │ │
│  └──────────────────────────────┘ │
│                                    │
│  ┌──────────────────────────────┐ │
│  │  Container Runtime           │ │
│  │  (Docker/containerd/CRI-O)   │ │
│  │                              │ │
│  │  ┌────┐ ┌────┐ ┌────┐       │ │
│  │  │Pod1│ │Pod2│ │Pod3│       │ │
│  │  └────┘ └────┘ └────┘       │ │
│  └──────────────────────────────┘ │
└────────────────────────────────────┘
```

**What it does**:

- Registers node with API server
- Watches for pods assigned to its node
- Starts/stops containers (via container runtime)
- Mounts volumes
- Reports pod and node status back to API server
- Runs liveness/readiness probes

**Important**: kubelet does NOT manage pods not created by Kubernetes (e.g., if you manually `docker run` on the node).

---

### 2. kube-proxy

**The network proxy: handles networking and load balancing.**

```
Request to Service "nginx" (ClusterIP: 10.96.1.100)
                │
                ↓
         ┌─────────────┐
         │ kube-proxy  │ (Uses iptables/ipvs rules)
         └─────────────┘
                │
     ┌──────────┼──────────┐
     ↓          ↓          ↓
  Pod 1      Pod 2      Pod 3
  (10.1.1.5) (10.1.1.6) (10.1.1.7)
```

**What it does**:

- Maintains network rules on each node
- Implements Kubernetes Service abstraction
- Routes traffic to the correct pods
- Load balances across pod replicas

**How it works**:

1. Watches API server for Service and Endpoint changes
2. Updates iptables/ipvs rules on the node
3. When traffic hits a Service IP, rules redirect to a pod IP

**Example**:

```bash
# When you create a Service
kubectl expose deployment nginx --port=80

# kube-proxy creates rules like:
# Traffic to 10.96.1.100:80 → Forward to one of [10.1.1.5:80, 10.1.1.6:80, 10.1.1.7:80]
```

**Proxy modes**:

- **iptables** (default): Uses Linux iptables rules
- **ipvs**: Uses IP Virtual Server (more efficient for large clusters)
- **userspace** (legacy): Proxies traffic in userspace (slower)

---

### 3. Container Runtime

**The container executor: actually runs containers.**

```
kubelet → Container Runtime Interface (CRI) → Container Runtime
                                                    │
                                    ┌───────────────┼───────────────┐
                                    ↓               ↓               ↓
                                 Docker        containerd        CRI-O
```

**Supported container runtimes**:

- **containerd** (most common, default in many distros)
- **CRI-O** (lightweight, designed for Kubernetes)
- **Docker** (via dockershim, deprecated in Kubernetes 1.24+)

**What it does**:

- Pulls container images
- Starts and stops containers
- Manages container lifecycle
- Provides container logs

**Note**: Kubernetes talks to container runtimes via the **Container Runtime Interface (CRI)**, not directly.

---

## Add-ons (Optional but Common)

These are not part of the core Kubernetes architecture but are commonly deployed.

### 1. DNS (CoreDNS)

**Provides service discovery within the cluster.**

```bash
# Pods can find each other by DNS name
curl http://nginx-service.default.svc.cluster.local
```

**What it does**:

- Runs as a pod in the cluster
- Provides DNS resolution for Services
- Enables service discovery

---

### 2. Dashboard

**Web-based UI for managing the cluster.**

---

### 3. Metrics Server

**Collects resource metrics (CPU, memory) from nodes and pods.**

Needed for:

- `kubectl top nodes`
- `kubectl top pods`
- Horizontal Pod Autoscaler (HPA)

---

### 4. Ingress Controller

**Manages external HTTP/HTTPS routing.**

Examples: nginx-ingress, Traefik, HAProxy

---

## How It All Works Together

Let's walk through what happens when you create a deployment.

### Scenario: You Run `kubectl create deployment nginx --image=nginx --replicas=3`

```
┌─────────────────────────────────────────────────────────────────┐
│  Step 1: kubectl → API Server                                   │
│  kubectl sends HTTP POST to API Server with deployment spec     │
└─────────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────────┐
│  Step 2: API Server → etcd                                      │
│  API Server validates request, writes deployment to etcd        │
└─────────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────────┐
│  Step 3: Deployment Controller watches API Server              │
│  Sees new deployment, creates a ReplicaSet                      │
│  Writes ReplicaSet to etcd (via API Server)                     │
└─────────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────────┐
│  Step 4: ReplicaSet Controller watches API Server              │
│  Sees new ReplicaSet with replicas=3                            │
│  Creates 3 pod definitions, writes to etcd (via API Server)     │
└─────────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────────┐
│  Step 5: Scheduler watches API Server                          │
│  Sees 3 pods without node assignment                            │
│  Evaluates nodes, assigns each pod to a node                    │
│  Updates pod specs with node assignment (via API Server)        │
└─────────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────────┐
│  Step 6: kubelet (on each node) watches API Server             │
│  Sees pods assigned to its node                                 │
│  Tells container runtime to pull nginx image and start container│
└─────────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────────┐
│  Step 7: kubelet reports pod status to API Server              │
│  Pod status: ContainerCreating → Running                        │
│  API Server updates etcd                                        │
└─────────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────────┐
│  Step 8: You check the status                                  │
│  kubectl get pods → Shows 3 pods running                        │
└─────────────────────────────────────────────────────────────────┘
```

**Key takeaway**: Every component only talks to the API server. The API server is the central hub.

---

## The Declarative Model: Desired State vs Actual State

This is the fundamental concept behind Kubernetes.

### Imperative (How to do it)

```bash
# Traditional approach (imperative)
docker run nginx
docker run nginx
docker run nginx
# If one crashes, you manually restart it
```

### Declarative (What you want)

```yaml
# Kubernetes approach (declarative)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 3
  # ... rest of spec
```

**Declarative philosophy**:

```
You declare: "I want 3 nginx pods running"
Kubernetes ensures: 3 nginx pods are always running

If a pod dies:
  - Controller notices: Actual (2) != Desired (3)
  - Controller creates a new pod
  - Actual state returns to desired state
```

**Why this matters**:

- **Self-healing**: Kubernetes automatically fixes issues
- **Idempotent**: Applying the same YAML multiple times has the same effect
- **Version control**: Your cluster state is defined in files
- **Repeatable**: Same YAML produces same results every time

---

## Hands-On Exercise 1: Explore Control Plane Components

Let's see these components in action.

### View Control Plane Pods

```bash
# Control plane components run as pods in the kube-system namespace
kubectl get pods -n kube-system

# Expected output:
# NAME                               READY   STATUS    RESTARTS   AGE
# coredns-5d78c9869d-abc12           1/1     Running   0          10m
# etcd-minikube                      1/1     Running   0          10m
# kube-apiserver-minikube            1/1     Running   0          10m
# kube-controller-manager-minikube   1/1     Running   0          10m
# kube-proxy-xyz45                   1/1     Running   0          10m
# kube-scheduler-minikube            1/1     Running   0          10m
# storage-provisioner                1/1     Running   0          10m
```

### Inspect API Server

```bash
# Get detailed info about the API server
kubectl describe pod -n kube-system kube-apiserver-minikube

# Look for:
# - Command: kube-apiserver with various flags
# - Mounts: Volumes for certificates, etcd access
```

### Check etcd

```bash
# View etcd pod
kubectl get pod -n kube-system -l component=etcd

# See etcd configuration
kubectl describe pod -n kube-system etcd-minikube
```

---

## Hands-On Exercise 2: Watch the Reconciliation Loop

Let's observe controllers in action.

### Create a Deployment

```bash
# Create a deployment with 3 replicas
kubectl create deployment test-nginx --image=nginx --replicas=3

# Immediately watch pods
kubectl get pods -l app=test-nginx -w
```

**You'll see**:

```
NAME                          READY   STATUS              RESTARTS   AGE
test-nginx-7d64f5b8c9-abc12   0/1     ContainerCreating   0          1s
test-nginx-7d64f5b8c9-def34   0/1     ContainerCreating   0          1s
test-nginx-7d64f5b8c9-ghi56   0/1     ContainerCreating   0          1s
test-nginx-7d64f5b8c9-abc12   1/1     Running             0          5s
test-nginx-7d64f5b8c9-def34   1/1     Running             0          6s
test-nginx-7d64f5b8c9-ghi56   1/1     Running             0          7s
```

### Test Self-Healing

```bash
# Delete one pod
kubectl delete pod test-nginx-7d64f5b8c9-abc12

# Immediately check
kubectl get pods -l app=test-nginx

# You'll see the deleted pod terminating and a NEW pod being created
# NAME                          READY   STATUS              RESTARTS   AGE
# test-nginx-7d64f5b8c9-abc12   1/1     Terminating         0          2m
# test-nginx-7d64f5b8c9-def34   1/1     Running             0          2m
# test-nginx-7d64f5b8c9-ghi56   1/1     Running             0          2m
# test-nginx-7d64f5b8c9-xyz99   0/1     ContainerCreating   0          1s
```

**What happened**:

1. ReplicaSet controller noticed: Desired (3) != Actual (2)
2. Controller created a new pod
3. Scheduler assigned it to a node
4. kubelet started the container

### Scale the Deployment

```bash
# Scale to 5 replicas
kubectl scale deployment test-nginx --replicas=5

# Watch the controller create 2 new pods
kubectl get pods -l app=test-nginx -w
```

**Cleanup**:

```bash
kubectl delete deployment test-nginx
```

---

## Hands-On Exercise 3: View Node Information

Let's explore worker nodes.

### List Nodes

```bash
# List all nodes
kubectl get nodes

# Expected output (minikube has 1 node):
# NAME       STATUS   ROLES           AGE   VERSION
# minikube   Ready    control-plane   1h    v1.28.3

# Get detailed node information
kubectl get nodes -o wide

# Shows: IP, OS, kernel version, container runtime
```

### Describe a Node

```bash
# Get detailed node information
kubectl describe node minikube

# Important sections:
# - Conditions: Ready, MemoryPressure, DiskPressure
# - Addresses: InternalIP, Hostname
# - Capacity: Total CPU and memory
# - Allocatable: Available CPU and memory for pods
# - Allocated resources: How much is currently used
# - Non-terminated Pods: List of pods on this node
```

**Example output**:

```
Capacity:
  cpu:                2
  memory:             2048Mi
Allocatable:
  cpu:                2
  memory:             1946Mi
Allocated resources:
  Resource           Requests     Limits
  --------           --------     ------
  cpu                250m (12%)   500m (25%)
  memory             128Mi (6%)   256Mi (13%)
```

### View kubelet Logs (if accessible)

```bash
# SSH into minikube node
minikube ssh

# View kubelet logs
sudo journalctl -u kubelet -f

# You'll see logs like:
# - Pulling images
# - Starting containers
# - Reporting pod status
# - Health check results

# Exit
exit
```

---

## Hands-On Exercise 4: Observe kube-proxy

Let's see how kube-proxy manages networking.

### Create a Service

```bash
# Create a deployment
kubectl create deployment web --image=nginx --replicas=3

# Expose it as a service
kubectl expose deployment web --port=80 --type=ClusterIP

# Get the service IP
kubectl get service web

# Expected output:
# NAME   TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
# web    ClusterIP   10.96.100.123   <none>        80/TCP    10s
```

### Examine kube-proxy

```bash
# Find kube-proxy pod
kubectl get pods -n kube-system -l k8s-app=kube-proxy

# View iptables rules (minikube)
minikube ssh

# View iptables rules created by kube-proxy
sudo iptables-save | grep web

# You'll see rules that redirect traffic from Service IP (10.96.100.123)
# to pod IPs

exit
```

**Cleanup**:

```bash
kubectl delete deployment web
kubectl delete service web
```

---

## Communication Patterns

### Component Communication

```
┌─────────────┐
│   kubectl   │
└──────┬──────┘
       │ HTTPS
       ↓
┌─────────────────────────────────────────┐
│           API Server                    │
│  (All communication goes through here)  │
└────┬────────┬────────┬─────────┬────────┘
     │        │        │         │
     ↓        ↓        ↓         ↓
  ┌──────┐ ┌────┐ ┌──────┐  ┌─────────┐
  │ etcd │ │Sched│ │Contr.│  │ kubelet │
  │      │ │uler│ │ Mgr  │  │(on node)│
  └──────┘ └────┘ └──────┘  └─────────┘
```

**Key patterns**:

1. **Everything goes through API server**: No component talks directly to another
2. **Watch mechanism**: Components watch API server for changes (long-polling)
3. **No direct etcd access**: Only API server reads/writes etcd
4. **Pull-based**: kubelet pulls pod specs from API server (not pushed)

---

## Challenges

### Challenge 1: Component Identification

Identify which component is responsible for each action.

**Actions**:

1. Decides which node should run a new pod
2. Notices a pod has crashed and creates a replacement
3. Starts a container on a worker node
4. Routes traffic from a Service to pods
5. Stores the cluster state

<details>
<summary>Solution</summary>

1. **Scheduler** (kube-scheduler)
2. **ReplicaSet Controller** (part of controller-manager)
3. **kubelet** (calls container runtime)
4. **kube-proxy**
5. **etcd**

</details>

---

### Challenge 2: Trace a Deployment Creation

Create a deployment and trace which components are involved at each step.

**Requirements**:

1. Create a deployment with 2 replicas
2. List the components involved in order
3. Explain what each component does
4. Verify the deployment is running
5. Delete the deployment

<details>
<summary>Solution</summary>

```bash
# 1. Create deployment
kubectl create deployment trace-test --image=nginx --replicas=2

# Components involved:
# a. kubectl → API Server (HTTP POST request)
# b. API Server → etcd (writes deployment object)
# c. Deployment Controller (watches API Server, creates ReplicaSet)
# d. API Server → etcd (writes ReplicaSet object)
# e. ReplicaSet Controller (watches API Server, creates 2 pod objects)
# f. API Server → etcd (writes pod objects)
# g. Scheduler (watches API Server, assigns pods to nodes)
# h. API Server → etcd (updates pod objects with node assignments)
# i. kubelet (watches API Server, sees pods assigned to its node)
# j. kubelet → Container Runtime (pulls image, starts containers)
# k. kubelet → API Server (reports pod status)
# l. API Server → etcd (updates pod status)

# 2. Verify
kubectl get deployments,replicasets,pods -l app=trace-test

# 3. Cleanup
kubectl delete deployment trace-test
```

</details>

---

### Challenge 3: Test Self-Healing

Test the self-healing capabilities and identify which component performs the healing.

**Requirements**:

1. Create a deployment with 4 replicas
2. Delete 2 pods
3. Observe new pods being created
4. Explain which component detected the issue
5. Explain which component created the new pods

<details>
<summary>Solution</summary>

```bash
# 1. Create deployment
kubectl create deployment heal-test --image=nginx --replicas=4

# 2. Wait for pods to be running
kubectl get pods -l app=heal-test

# 3. Delete 2 pods (replace with actual pod names)
kubectl delete pod heal-test-xxxxx-abc12 heal-test-xxxxx-def34

# 4. Immediately check
kubectl get pods -l app=heal-test -w

# You'll see:
# - Deleted pods: Terminating
# - New pods: ContainerCreating → Running

# 5. Explanation:
# - ReplicaSet Controller detected the issue
#   (It continuously watches: Desired=4, Actual=2, creates 2 new pods)
# - Scheduler assigned new pods to nodes
# - kubelet started the containers

# Cleanup
kubectl delete deployment heal-test
```

</details>

---

### Challenge 4: Explore Node Resources

Examine how Kubernetes tracks node resources.

**Requirements**:

1. Get node capacity and allocatable resources
2. Create a deployment that uses significant resources
3. Check how resources are allocated
4. Explain the difference between capacity and allocatable

<details>
<summary>Solution</summary>

```bash
# 1. View node resources
kubectl describe node minikube | grep -A 5 "Capacity:"
kubectl describe node minikube | grep -A 5 "Allocatable:"

# Capacity: Total hardware resources
# Allocatable: Resources available for pods (Capacity - system reservations)

# 2. Create resource-intensive deployment
cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: resource-test
spec:
  replicas: 2
  selector:
    matchLabels:
      app: resource-test
  template:
    metadata:
      labels:
        app: resource-test
    spec:
      containers:
      - name: nginx
        image: nginx
        resources:
          requests:
            memory: "128Mi"
            cpu: "250m"
          limits:
            memory: "256Mi"
            cpu: "500m"
EOF

# 3. Check allocated resources
kubectl describe node minikube | grep -A 10 "Allocated resources:"

# You'll see:
# Resource           Requests          Limits
# cpu                500m (25%)        1000m (50%)
# memory             256Mi (13%)       512Mi (26%)

# 4. Cleanup
kubectl delete deployment resource-test
```

</details>

---

### Challenge 5: View Component Logs

Access logs from various Kubernetes components.

**Requirements**:

1. View API server logs
2. View scheduler logs
3. View controller manager logs
4. Identify an interesting event in each log

<details>
<summary>Solution</summary>

```bash
# 1. API server logs
kubectl logs -n kube-system kube-apiserver-minikube | tail -20

# Look for: API requests, authentication events

# 2. Scheduler logs
kubectl logs -n kube-system kube-scheduler-minikube | tail -20

# Look for: Pod scheduling decisions

# 3. Controller manager logs
kubectl logs -n kube-system kube-controller-manager-minikube | tail -20

# Look for: Reconciliation loops, resource creation

# To follow logs in real-time:
kubectl logs -n kube-system kube-scheduler-minikube -f

# Create a deployment in another terminal
kubectl create deployment log-test --image=nginx

# You'll see the scheduler logs showing the pod being scheduled

# Cleanup
kubectl delete deployment log-test
```

</details>

---

## Best Practices

1. **Understand the Flow**
   - Always remember: All components talk to API server
   - API server is the only component that writes to etcd
   - Controllers continuously reconcile desired vs actual state

2. **Back Up etcd**
   - etcd contains all cluster state
   - Losing etcd means losing the cluster
   - Always have etcd backups in production

3. **Monitor Control Plane**
   - API server, scheduler, controller manager are critical
   - If control plane fails, existing pods keep running
   - But new pods can't be scheduled

4. **Resource Planning**
   - Understand node capacity vs allocatable
   - Set resource requests and limits on pods
   - Monitor node resource usage

5. **HA Control Plane (Production)**
   - Run multiple replicas of control plane components
   - Use load balancer for API server
   - Run etcd cluster (3 or 5 nodes)

---

## Key Takeaways

1. **Control Plane** manages the cluster: API Server, etcd, Scheduler, Controller Manager
2. **API Server** is the central hub - all components communicate through it
3. **etcd** is the source of truth - all cluster state is stored here
4. **Scheduler** assigns pods to nodes based on resources and constraints
5. **Controllers** continuously reconcile desired state with actual state
6. **Worker Nodes** run workloads: kubelet, kube-proxy, container runtime
7. **kubelet** ensures containers are running on its node
8. **kube-proxy** handles networking and service load balancing
9. **Declarative model**: You declare desired state, Kubernetes maintains it
10. **Self-healing**: Controllers automatically fix issues (pod crashes, node failures)

---

## Next Steps

In [Lesson 04: Working with kubectl](lesson-04-working-with-kubectl.md), you'll learn:

- Essential kubectl commands for daily work
- Imperative vs declarative resource management
- YAML manifest deep dive
- Advanced debugging techniques
- kubectl productivity tips and tricks

---

## Questions to Ponder

1. What happens if the API server crashes? Will existing pods continue running?
2. If etcd is the only stateful component, what happens to controllers if they crash?
3. Why does every component communicate through the API server instead of directly?
4. How does Kubernetes handle the situation when a worker node goes offline?
5. What's the difference between a node being "NotReady" vs being completely offline?

---

## Cleanup

```bash
# No specific cleanup needed for this lesson
# All exploration was read-only

# If you created any test resources:
kubectl delete deployment --all
kubectl delete service --all
```

---

**Congratulations!** You now understand how Kubernetes works internally. You know the role of each component, how they communicate, and how the declarative model enables self-healing. In the next lesson, we'll dive deep into kubectl to master interacting with your cluster.
