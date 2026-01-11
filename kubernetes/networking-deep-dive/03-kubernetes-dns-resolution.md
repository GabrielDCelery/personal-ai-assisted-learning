# How Kubernetes Knows to Direct Traffic to reverse-proxy

## Question
How does Kubernetes know that it should direct traffic to http://reverse-proxy?

## Answer

Through **Kubernetes DNS resolution** - every Service gets an automatic DNS entry.

## Kubernetes DNS Resolution

When the cloudflared container (running in the `homelab` namespace) tries to connect to `http://reverse-proxy:80`, Kubernetes automatically resolves it through its internal DNS (CoreDNS).

### DNS Resolution Rules

**Same namespace (simple):**
```
reverse-proxy  →  resolves to the Service IP in the same namespace
```

**Different namespace (full name):**
```
reverse-proxy.homelab.svc.cluster.local  →  resolves to Service IP in homelab namespace
```

### In Your Setup

1. **Cloudflared pod** is in `homelab` namespace
2. **Service `reverse-proxy`** is in `homelab` namespace
3. When cloudflared tries `http://reverse-proxy:80`, Kubernetes DNS resolves it to the Service's ClusterIP
4. The Service forwards traffic to pods with label `app: traefik`

## Testing DNS Resolution

```bash
# Check the service exists and has an IP
kubectl get svc -n homelab

# Exec into the cloudflared pod and test DNS
kubectl exec -n homelab deployment/cloudflared-tunnel -- nslookup reverse-proxy

# Or test the connection directly
kubectl exec -n homelab deployment/cloudflared-tunnel -- wget -O- http://reverse-proxy:80
```

## Flow Diagram

```
Internet
  ↓
Cloudflare Tunnel (manages routing)
  ↓
cloudflared pod (in homelab namespace)
  ↓ connects to http://reverse-proxy:80
Kubernetes DNS resolves "reverse-proxy" → ClusterIP (e.g., 10.43.123.45)
  ↓
Service "reverse-proxy" (in homelab namespace)
  ↓ forwards to pods matching selector: app=traefik
Traefik pod
```

## DNS Name Formats

Kubernetes creates these DNS names for every service:

```bash
# Short name (same namespace only)
reverse-proxy

# Namespace qualified
reverse-proxy.homelab

# Fully qualified domain name (FQDN)
reverse-proxy.homelab.svc.cluster.local
```

All three resolve to the same Service IP when queried from within the cluster.

## How CoreDNS Works

**CoreDNS** is the DNS server running in your cluster:

```bash
# Check CoreDNS
kubectl get pods -n kube-system | grep coredns
```

**When a pod makes a DNS query:**
1. Pod's `/etc/resolv.conf` points to CoreDNS
2. CoreDNS watches Kubernetes API for Services
3. CoreDNS returns the Service's ClusterIP
4. Pod connects to that IP
5. iptables rules forward to actual Pod

## Example Resolution

```bash
# Inside a pod
$ nslookup reverse-proxy

Server:         10.43.0.10
Address:        10.43.0.10#53

Name:   reverse-proxy.homelab.svc.cluster.local
Address: 10.43.152.220
```

The magic is **Kubernetes DNS** - every service gets an automatic DNS entry that pods can use to find each other!
