# Traefik Configuration Breakdown

## Question
How does the Traefik config work - all the fields and their behavior?

## Answer

## Deployment Section

```yaml
apiVersion: apps/v1
kind: Deployment
```
- **What it is**: Declares this as a Kubernetes Deployment resource
- **Behavior**: Creates and manages pods; ensures desired number of replicas are running

```yaml
metadata:
  name: traefik
  namespace: homelab
```
- **name**: The deployment is named "traefik"
- **namespace**: Lives in the "homelab" namespace
- **Behavior**: All commands need `-n homelab` to access this deployment

```yaml
spec:
  replicas: 1
```
- **What it does**: Run exactly 1 copy of the Traefik pod
- **Behavior**: If the pod crashes, Kubernetes automatically creates a new one
- **Scale**: Change to `2` or `3` for high availability

```yaml
  selector:
    matchLabels:
      app: traefik
```
- **What it does**: Tells the Deployment which pods it manages
- **Behavior**: Looks for pods with label `app: traefik`
- **Critical**: This must match the labels in `template.metadata.labels` below

```yaml
  template:
    metadata:
      labels:
        app: traefik
```
- **What it does**: Defines the pod template and applies labels to created pods
- **Behavior**: Every pod created gets the label `app: traefik`
- **Why important**: Services use these labels to find pods to route traffic to

```yaml
    spec:
      containers:
        - name: traefik
          image: traefik:v3.6
```
- **name**: Container name within the pod (just for identification)
- **image**: Uses Traefik version 3.6 from Docker Hub
- **Behavior**: Kubernetes pulls this image and runs it

## Traefik Args - Where Traefik's Behavior is Configured

```yaml
          args:
            - --api.dashboard=true
```
- **What it does**: Enables the Traefik web dashboard UI
- **Behavior**: Dashboard becomes available at `/dashboard/` endpoint
- **Access**: By default on port 8080 (the API port)

```yaml
            - --api.insecure=true
```
- **What it does**: Allows dashboard access without authentication
- **Behavior**: Dashboard is publicly accessible (no username/password required)
- **Security**: ⚠️ Use only behind Cloudflare Access or similar protection
- **Production**: Set to `false` and configure proper auth

```yaml
            - --entrypoints.web.address=:80
```
- **What it does**: Creates an entrypoint (listening port) named "web" on port 80
- **Behavior**: Traefik listens for HTTP traffic on port 80
- **Routing**: Routes and Ingresses can specify which entrypoint to use
- **Multiple entrypoints**: You could add `--entrypoints.websecure.address=:443` for HTTPS

```yaml
            - --log.level=INFO
```
- **What it does**: Sets log verbosity
- **Options**: `DEBUG` (very verbose) → `INFO` (normal) → `WARN` → `ERROR` (minimal)
- **Behavior**: Shows request routing, middleware application, errors
- **View logs**: `kubectl logs -n homelab deployment/traefik`

## Container Ports

```yaml
          ports:
            - name: web
              containerPort: 80
```
- **name**: Labels this port as "web" (just documentation)
- **containerPort**: Tells Kubernetes the container listens on port 80
- **Behavior**: Kubernetes can route traffic to this port
- **Matches**: The `--entrypoints.web.address=:80` arg above

```yaml
            - name: dashboard
              containerPort: 8080
```
- **name**: Labels this port as "dashboard"
- **containerPort**: Traefik's API/dashboard port
- **Behavior**: Dashboard and API are accessible here
- **Default**: Port 8080 is Traefik's built-in API port (can't change without config)

## Service Section

```yaml
apiVersion: v1
kind: Service
metadata:
  name: reverse-proxy
  namespace: homelab
```
- **What it is**: Creates a stable network endpoint
- **name**: DNS name will be `reverse-proxy.homelab.svc.cluster.local`
- **Behavior**: Provides a single IP address even if pods restart (ClusterIP doesn't change)

```yaml
spec:
  type: ClusterIP
```
- **What it does**: Creates an internal-only IP address
- **Behavior**: Only accessible from within the cluster (not from outside)
- **Alternatives**:
  - `NodePort`: Exposes on each node's IP
  - `LoadBalancer`: Gets external IP (cloud providers)
- **Your setup**: Perfect because Cloudflare tunnel connects from inside the cluster

```yaml
  selector:
    app: traefik
```
- **Critical**: This is how the Service finds pods to route traffic to
- **Behavior**: Looks for any pod with label `app: traefik`
- **Load balancing**: If multiple pods exist (replicas > 1), distributes traffic among them
- **Dynamic**: If pods are added/removed, Service automatically updates

```yaml
  ports:
    - name: web
      port: 80
      targetPort: 80
```
- **name**: Labels this port mapping as "web" (just documentation)
- **port: 80**: The Service listens on port 80
- **targetPort: 80**: Forwards traffic to port 80 on the pod
- **Behavior**: `http://reverse-proxy:80` → forwards to container port 80
- **Flow**: `cloudflared → Service (port 80) → Pod (targetPort 80)`

```yaml
    - name: dashboard
      port: 8080
      targetPort: 8080
```
- **port: 8080**: The Service listens on port 8080
- **targetPort: 8080**: Forwards to container port 8080 (dashboard)
- **Behavior**: `http://reverse-proxy:8080` → forwards to Traefik dashboard
- **Your setup**: Cloudflare tunnel now points here to show the dashboard

## Traffic Flow Diagram

```
Internet
  ↓
Cloudflare (homelab-dev.gaborzeller.com)
  ↓
cloudflared pod (in cluster)
  ↓ DNS lookup: reverse-proxy → 10.43.152.220
Service "reverse-proxy:8080"
  ↓ selector: app=traefik → finds pod
Traefik pod (port 8080)
  ↓ Traefik dashboard
User sees dashboard UI
```

## What's Missing from Current Config

**No Providers Configured**: Traefik doesn't know how to discover services yet. To route traffic to other apps, you need to add:

```yaml
args:
  - --providers.kubernetesIngress=true  # Read Kubernetes Ingress resources
  - --providers.kubernetesCRD=true      # Read Traefik's IngressRoute CRDs
```

**No Service Account**: Traefik can't read Kubernetes API yet (needs RBAC permissions)

**Current state**: The dashboard works, but Traefik has no routes to other services
