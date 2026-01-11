# Service Port Routing: How Traffic Flows from Service to Container

## Question
When the service receives a request, how does the routing happen to port 8080? Is it routed to the service or the containers?

## Answer

Traffic goes **directly to the container**, not to the service itself.

## Service Port Mapping Flow

```yaml
ports:
  - name: dashboard
    port: 8080          # Service listens on this
    targetPort: 8080    # Forwards to container's this port
```

### Step-by-Step Request Flow

**1. Request arrives at the Service:**
```
cloudflared → http://reverse-proxy:8080
              ↓
        Service IP: 10.43.152.220:8080
```

**2. Service does lookup using selector:**
```yaml
selector:
  app: traefik  # Find all pods with this label
```
The Service queries Kubernetes: "Give me all pods with label `app: traefik`"

Result: Pod IP = `10.42.0.15` (example)

**3. Service forwards to Pod (not to itself):**
```
Service translates the request:
FROM: http://10.43.152.220:8080 (Service IP)
TO:   http://10.42.0.15:8080    (Pod IP:targetPort)
```

**4. Container receives the request:**
The container running Traefik receives traffic directly on port 8080

## Visual Diagram

```
Request: http://reverse-proxy:8080
         ↓
    [Service: reverse-proxy]
    - ClusterIP: 10.43.152.220
    - Port: 8080 (listening port)
    - Selector: app=traefik
         ↓
    Kubernetes queries for pods
    with label app=traefik
         ↓
    Found: Pod "traefik-6df5979f88-2x7zw"
    Pod IP: 10.42.0.15
         ↓
    Service forwards to:
    10.42.0.15:8080 (targetPort)
         ↓
    [Container: traefik]
    - Receives request on port 8080
    - Traefik processes it
```

## Important Distinctions

**Service acts as a proxy/load balancer:**
- Service has its own IP address (ClusterIP)
- Service does **NOT** run as a process - it's implemented by `kube-proxy` using iptables/IPVS
- Service forwards traffic directly to pods (not through another hop)

**The `targetPort` is the container's port:**
- This must match the port the container is actually listening on
- If Traefik dashboard is on 8080, `targetPort` must be 8080
- You can map differently, e.g., `port: 80, targetPort: 8080`

## Example with Different Ports

Let's say you wanted to access the dashboard on port 9000 instead:

```yaml
ports:
  - name: dashboard
    port: 9000         # Service listens on 9000
    targetPort: 8080   # But forwards to container's 8080
```

**Flow:**
```
Request: http://reverse-proxy:9000
         ↓
Service receives on port 9000
         ↓
Service forwards to Pod port 8080
         ↓
Container receives on 8080 (where Traefik dashboard actually listens)
```

**In your current config:**
```yaml
port: 8080
targetPort: 8080
```
Both are the same, so no translation happens - just forwarding.

## Verify This Yourself

Check the actual IPs:

```bash
# Service IP
kubectl get svc reverse-proxy -n homelab -o wide

# Pod IP
kubectl get pods -n homelab -o wide

# See the iptables rules (advanced)
kubectl get endpoints reverse-proxy -n homelab
```

The **Endpoints** object shows exactly which Pod IPs the Service forwards to:

```bash
kubectl get endpoints reverse-proxy -n homelab -o yaml
```

You'll see something like:
```yaml
subsets:
- addresses:
  - ip: 10.42.0.15  # Pod IP
  ports:
  - port: 8080      # targetPort
```

This proves the Service forwards **directly to the Pod IP:targetPort**, not to itself!
