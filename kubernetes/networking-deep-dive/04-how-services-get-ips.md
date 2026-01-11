# How Services Get IP Addresses

## Question
How come the service got an IP address?

## Answer

When you create a Service with `type: ClusterIP`, Kubernetes **automatically assigns** it an IP address from a pre-configured IP range called the **Service CIDR**.

## Step-by-Step Process

**1. You create the Service:**
```bash
kubectl apply -f kubernetes/03-traefik.yaml
```

**2. Kubernetes API Server receives the request:**
- Sees `type: ClusterIP` (or no type, since ClusterIP is default)
- Allocates an IP from the Service CIDR pool
- Stores this in etcd

**3. Service gets a ClusterIP:**
```
Service: reverse-proxy
ClusterIP: 10.43.152.220  ← Automatically assigned
```

**4. kube-proxy implements the routing:**
- Runs on every node in the cluster
- Watches for Service changes
- Creates iptables/IPVS rules to make that IP work

## What is the Service CIDR?

It's an IP range reserved for Services, configured when the cluster is created.

**Check your cluster's Service CIDR:**
```bash
# For k3s (your setup)
kubectl cluster-info dump | grep service-cluster-ip-range

# Or check a Service to see the pattern
kubectl get svc -A -o wide
```

In your case, `10.43.152.220` suggests your Service CIDR is probably `10.43.0.0/16`.

## The IP is "Virtual" - Not Real!

**Important**: The Service IP doesn't exist on any network interface. It's a **virtual IP** implemented by networking rules.

### Verify this:

```bash
# SSH to your k3s node
# Try to find the IP on any interface
ip addr show | grep 10.43.152.220
# Result: Nothing! It doesn't exist as a real interface
```

**How does it work then?**

`kube-proxy` creates iptables rules that intercept traffic to `10.43.152.220` and rewrite it to go to Pod IPs.

## Behind the Scenes: iptables/IPVS

When a request goes to `http://10.43.152.220:8080`:

**1. Packet enters the network stack:**
```
Destination: 10.43.152.220:8080
```

**2. iptables intercepts it (before routing):**
```bash
# kube-proxy created rules like:
iptables -t nat -A KUBE-SERVICES -d 10.43.152.220/32 -p tcp --dport 8080 \
  -j KUBE-SVC-XXXXX
```

**3. Rule rewrites the destination:**
```
Before: 10.43.152.220:8080
After:  10.42.0.15:8080  (Pod IP)
```

**4. Packet is routed to the Pod:**
The rewritten packet goes directly to the Pod IP.

## Visualization

```
┌─────────────────────────────────────────────┐
│         Kubernetes Cluster                  │
│                                             │
│  Service CIDR: 10.43.0.0/16                │
│  ├─ 10.43.0.1      (kube-dns)              │
│  ├─ 10.43.152.220  (reverse-proxy) ← Auto assigned │
│  └─ 10.43.xxx.xxx  (other services)        │
│                                             │
│  Pod CIDR: 10.42.0.0/16                    │
│  ├─ 10.42.0.11     (cloudflared pod)       │
│  ├─ 10.42.0.15     (traefik pod)           │
│  └─ 10.42.x.x      (other pods)            │
│                                             │
│  [kube-proxy on each node]                 │
│   - Watches Services                        │
│   - Creates iptables rules                  │
│   - Makes virtual IPs work                  │
└─────────────────────────────────────────────┘

Request Flow:
1. cloudflared: "Connect to 10.43.152.220:8080"
2. iptables: "Rewrite to 10.42.0.15:8080"
3. traefik pod receives request
```

## Who Assigns the IP?

**Kubernetes API Server** has a component called the **Service Controller** that:

1. Watches for new Service resources
2. Allocates an IP from the available pool
3. Stores it in the Service object
4. Notifies kube-proxy on all nodes

**Allocation is automatic** - you can't choose the IP (unless you specify `clusterIP: x.x.x.x` explicitly).

## Different Service Types

```yaml
type: ClusterIP
```
- **Gets**: Virtual IP from Service CIDR
- **Accessible**: Only inside the cluster

```yaml
type: NodePort
```
- **Gets**: ClusterIP + opens a port (30000-32767) on every node
- **Accessible**: From outside via `<NodeIP>:<NodePort>`

```yaml
type: LoadBalancer
```
- **Gets**: ClusterIP + NodePort + external IP from cloud provider
- **Accessible**: From internet via external load balancer

```yaml
type: ExternalName
```
- **Gets**: No ClusterIP at all!
- **Behavior**: Just a DNS CNAME to an external service

## See it in Action

```bash
# Create a Service and watch it get an IP
kubectl create service clusterip test-svc --tcp=80:80 -n homelab

# Immediately check - IP is assigned
kubectl get svc test-svc -n homelab
# NAME       TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)
# test-svc   ClusterIP   10.43.241.185   <none>        80/TCP

# The IP came from nowhere - Kubernetes assigned it!

# Clean up
kubectl delete svc test-svc -n homelab
```

## Summary

1. **Service CIDR** is a pre-configured IP range (e.g., `10.43.0.0/16`)
2. **Kubernetes API Server** automatically assigns an IP from this range when you create a Service
3. **The IP is virtual** - implemented by kube-proxy using iptables/IPVS rules
4. **Pods can reach it** because kube-proxy rewrites the destination to actual Pod IPs
5. **You don't need to configure anything** - it's all automatic!
