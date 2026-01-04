# Kubernetes Learning Path

**From Containers to Orchestration**

## Overview

This Kubernetes course builds on your Docker knowledge to help you master container orchestration at scale. While Docker taught you how to package and run individual containers, Kubernetes teaches you how to manage hundreds or thousands of containers across multiple machines in production.

**Prerequisites:** Complete the Docker course first. You should be comfortable with containers, images, volumes, and networking.

---

## Why Kubernetes?

Docker gets you from "works on my machine" to "works in a container."
Kubernetes gets you from "works in a container" to "works at scale in production."

### Problems Kubernetes Solves

- **Manual scaling:** Manually starting 10 copies of your app across different servers
- **No auto-healing:** Containers crash and stay down until you manually restart them
- **Complex networking:** Managing load balancing across multiple container instances
- **Rolling updates:** Deploying new versions without downtime requires complex scripting
- **Resource management:** Efficiently packing containers onto servers to maximize resource usage
- **Service discovery:** Containers need to find and communicate with each other dynamically
- **Configuration management:** Managing environment variables and secrets across many containers

---

## Learning Path Structure

This course is divided into 5 modules with 20 comprehensive lessons.

### Module 1: Kubernetes Foundations (Lessons 1-4)

Building blocks of Kubernetes

- **Lesson 01: Introduction to Kubernetes** ✓
  - What is Kubernetes and why do we need it?
  - Kubernetes vs Docker: Understanding the relationship
  - Setting up your local Kubernetes environment (minikube/kind)
  - Your first deployment: From docker run to kubectl run

- **Lesson 02: Understanding Pods**
  - Pods: The smallest deployable unit
  - Single vs multi-container pods
  - Pod lifecycle and states
  - Hands-on: Creating and inspecting pods

- **Lesson 03: Kubernetes Architecture**
  - Control plane components (API Server, Scheduler, Controller Manager, etcd)
  - Node components (kubelet, kube-proxy, container runtime)
  - How components work together
  - The declarative model: Desired state vs actual state

- **Lesson 04: Working with kubectl**
  - Essential kubectl commands
  - Imperative vs declarative management
  - YAML manifests deep dive
  - Debugging techniques (logs, describe, exec)

---

### Module 2: Workload Management (Lessons 5-8)

Running and managing applications

- **Lesson 05: Deployments**
  - Beyond pods: Why we need Deployments
  - ReplicaSets under the hood
  - Scaling applications
  - Rolling updates and rollbacks

- **Lesson 06: StatefulSets and DaemonSets**
  - StatefulSets for stateful applications (databases)
  - DaemonSets for node-level services (monitoring, logging)
  - When to use each workload type
  - Hands-on: Deploying a MySQL StatefulSet

- **Lesson 07: Jobs and CronJobs**
  - Running batch workloads
  - One-time jobs vs scheduled jobs
  - Job patterns (parallel processing, work queues)
  - Hands-on: Data processing pipeline

- **Lesson 08: Application Configuration**
  - Environment variables in Kubernetes
  - ConfigMaps for configuration data
  - Secrets for sensitive information
  - Best practices for configuration management

---

### Module 3: Networking (Lessons 9-12)

Connecting containers and exposing services

- **Lesson 09: Services Fundamentals**
  - The problem: Dynamic pod IPs
  - ClusterIP: Internal communication
  - NodePort: External access (development)
  - LoadBalancer: External access (production)

- **Lesson 10: Service Discovery and DNS**
  - How pods find each other
  - Kubernetes DNS system
  - Headless services
  - Cross-namespace communication

- **Lesson 11: Ingress Controllers**
  - Beyond LoadBalancers: HTTP/HTTPS routing
  - Path-based and host-based routing
  - TLS/SSL termination
  - Hands-on: Deploying nginx-ingress

- **Lesson 12: Network Policies**
  - Default: All pods can talk to all pods
  - Implementing network segmentation
  - Ingress and egress rules
  - Security best practices

---

### Module 4: Storage and State (Lessons 13-16)

Managing persistent data

- **Lesson 13: Volumes in Kubernetes**
  - Docker volumes vs Kubernetes volumes
  - Volume types: emptyDir, hostPath, configMap, secret
  - Ephemeral vs persistent storage
  - Hands-on: Shared storage between containers

- **Lesson 14: Persistent Volumes and Claims**
  - The PV/PVC model
  - Storage classes and dynamic provisioning
  - Access modes and reclaim policies
  - Hands-on: Database with persistent storage

- **Lesson 15: StatefulSet Storage**
  - VolumeClaimTemplates
  - Stable network identities and storage
  - Ordered deployment and scaling
  - Hands-on: Multi-replica database cluster

- **Lesson 16: Storage Best Practices**
  - Choosing the right storage solution
  - Backup and disaster recovery
  - Performance considerations
  - Cloud provider storage options

---

### Module 5: Production Readiness (Lessons 17-20)

Taking applications to production

- **Lesson 17: Health Checks and Self-Healing**
  - Liveness probes: Detecting broken containers
  - Readiness probes: Managing traffic flow
  - Startup probes: Handling slow-starting apps
  - Hands-on: Building resilient applications

- **Lesson 18: Resource Management**
  - CPU and memory requests vs limits
  - Quality of Service (QoS) classes
  - ResourceQuotas and LimitRanges
  - Preventing resource starvation

- **Lesson 19: Autoscaling**
  - Horizontal Pod Autoscaler (HPA)
  - Vertical Pod Autoscaler (VPA)
  - Cluster Autoscaler
  - Hands-on: Auto-scaling web application

- **Lesson 20: Security and RBAC**
  - Security contexts and pod security
  - Role-Based Access Control (RBAC)
  - Service accounts
  - Security best practices checklist

---

## Hands-On Philosophy

Each lesson includes:

1. **Core Concepts**: Clear explanations with diagrams
2. **Hands-On Exercises**: Step-by-step tutorials
3. **Challenges**: Real-world problems to solve
4. **Solutions**: Complete working examples
5. **Best Practices**: Production-ready patterns
6. **Key Takeaways**: Summary of critical points

---

## Local Development Environment

We'll primarily use **minikube** for local Kubernetes development:

```bash
# Install minikube (example for Linux)
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube

# Start your cluster
minikube start

# Verify installation
kubectl cluster-info
kubectl get nodes
```

Alternative options:

- **kind** (Kubernetes in Docker): Lightweight, runs in Docker
- **Docker Desktop**: Built-in Kubernetes (Mac/Windows)
- **k3s**: Lightweight Kubernetes for resource-constrained environments

---

## Tools You'll Use

- **kubectl**: The Kubernetes command-line tool
- **minikube**: Local Kubernetes cluster
- **k9s** (optional): Terminal-based UI for Kubernetes
- **helm** (later lessons): Package manager for Kubernetes

---

## Learning Strategy

1. **Complete Docker course first**: Kubernetes builds on container knowledge
2. **Type every command**: Don't copy-paste; muscle memory matters
3. **Break things**: Delete resources, watch what happens, understand recovery
4. **Read the YAML**: Every manifest teaches you Kubernetes concepts
5. **Use kubectl explain**: Built-in documentation (e.g., `kubectl explain pod.spec`)
6. **Check logs and events**: Debugging is a critical skill

---

## From Docker to Kubernetes: Key Differences

| Docker Concept           | Kubernetes Equivalent            | Why Different?                       |
| ------------------------ | -------------------------------- | ------------------------------------ |
| `docker run`             | Pod + Deployment                 | Kubernetes manages multiple replicas |
| `docker-compose.yml`     | Multiple manifests or Helm chart | More powerful orchestration          |
| Container restart policy | Liveness/Readiness probes        | More sophisticated health checking   |
| Docker network           | Service + NetworkPolicy          | Designed for distributed systems     |
| Docker volume            | PersistentVolume + PVC           | Supports distributed storage         |
| `docker ps`              | `kubectl get pods`               | Cluster-wide visibility              |
| Environment variables    | ConfigMaps + Secrets             | Centralized configuration            |

---

## Course Progression

### Weeks 1-2: Foundations

Master the basics: pods, deployments, services

### Weeks 3-4: Core Concepts

Storage, configuration, networking patterns

### Weeks 5+: Production Skills

Scaling, monitoring, security, real-world projects

---

## Next Steps

Start with [Lesson 01: Introduction to Kubernetes](lesson-01-introduction.md) to begin your journey from containers to orchestration.

By the end of this course, you'll be able to:

- Deploy production applications on Kubernetes
- Scale services automatically based on load
- Manage stateful applications like databases
- Implement zero-downtime deployments
- Secure and monitor your clusters
- Debug issues in distributed systems

---

## Questions Before You Start?

- "Do I need a powerful computer?" → No, minikube runs well on 4GB RAM
- "Do I need cloud credits?" → No, everything works locally
- "How is this different from Docker Compose?" → You'll learn this in Lesson 01
- "Should I learn Kubernetes if I only use Docker locally?" → If you plan to work with production systems, yes

Let's build production-grade container orchestration skills together!
