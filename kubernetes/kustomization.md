● Great question! There are several approaches to parameterize the domain name in Kubernetes. Let me outline the best options:

1. Kustomize (Recommended for Kubernetes)

Kustomize is built into kubectl and is the most Kubernetes-native approach.

```txt
  Structure:
  kubernetes/
  ├── base/
  │   ├── kustomization.yaml
  │   ├── 0303-traefik-deployment.yaml
  │   ├── 0307-traefik-dashboard-route.yaml
  │   └── ...
  └── overlays/
      ├── production/
      │   ├── kustomization.yaml
      │   └── config.env
      └── staging/
          ├── kustomization.yaml
          └── config.env

```

base/kustomization.yaml:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - 0303-traefik-deployment.yaml
  - 0307-traefik-dashboard-route.yaml

# Define variables to be replaced
vars:
  - name: DOMAIN
    objref:
      kind: ConfigMap
      name: domain-config
      apiVersion: v1
    fieldref:
      fieldpath: data.domain

overlays/production/kustomization.yaml:
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

bases:
  - ../../base

# Create ConfigMap with your domain
configMapGenerator:
  - name: domain-config
    literals:
      - domain=mydomain.com

# Replace variables in manifests
replacements:
  - source:
      kind: ConfigMap
      name: domain-config
      fieldPath: data.domain
    targets:
      - select:
          kind: HTTPRoute
          name: traefik-dashboard
        fieldPaths:
          - spec.hostnames.0
```

Your HTTPRoute would use:
spec:
hostnames: - "traefik.$(DOMAIN)" # Gets replaced with traefik.mydomain.com

Deploy with:
kubectl apply -k overlays/production

---

2. Helm (If you want more complex templating)

Convert your manifests to Helm templates.

```txt
  Structure:
  helm-chart/
  ├── Chart.yaml
  ├── values.yaml
  ├── values-production.yaml
  └── templates/
      ├── 0303-traefik-deployment.yaml
      └── 0307-traefik-dashboard-route.yaml

```

values.yaml:

```yaml
domain: mydomain.com

traefik:
  subdomain: traefik
  image: traefik:v3.6
  replicas: 1

templates/0307-traefik-dashboard-route.yaml:
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: traefik-dashboard
  namespace: homelab
spec:
  hostnames:
    - "{{ .Values.traefik.subdomain }}.{{ .Values.domain }}"
  rules:
    - matches:
        - path:
            type: PathPrefix
            value: /
      backendRefs:
        - name: reverse-proxy
          port: 8080
```

Deploy with:

```sh
helm install homelab ./helm-chart -f values-production.yaml

```

---

3. Environment Variables via ConfigMap (Simplest)

Create a ConfigMap and reference it in your manifests.

domain-config.yaml:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: domain-config
  namespace: homelab
data:
  domain: "mydomain.com"
```

Problem: Kubernetes doesn't support direct ConfigMap substitution in resource fields like hostnames. You'd need an admission webhook or external tool.

---

4. Terraform/External Templating

Since you're already using Terraform (I see terraform/ in your repo), you could template the YAML files.

```terraform
variables.tf:
variable "domain" {
  description = "Base domain for services"
  type        = string
  default     = "mydomain.com"
}

main.tf:
resource "kubernetes_manifest" "traefik_dashboard_route" {
  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1"
    kind       = "HTTPRoute"
    metadata = {
      name      = "traefik-dashboard"
      namespace = "homelab"
    }
    spec = {
      hostnames = [
        "traefik.${var.domain}"
      ]
      parentRefs = [{
        name         = "homelab-gateway"
        namespace    = "homelab"
        sectionName  = "http"
      }]
      rules = [{
        matches = [{
          path = {
            type  = "PathPrefix"
            value = "/"
          }
        }]
        backendRefs = [{
          name = "reverse-proxy"
          port = 8080
        }]
      }]
    }
  }
}
```

---

5. Simple Script/envsubst (Quick & Dirty)

Use shell variable substitution.

Create templates:

# Save as 0307-traefik-dashboard-route.yaml.tmpl

```yaml
spec:
hostnames:
  - "traefik.${DOMAIN}"
```

deployment script:

```sh
#!/bin/bash
export DOMAIN="mydomain.com"
envsubst < 0307-traefik-dashboard-route.yaml.tmpl | kubectl apply -f -

```

---

My Recommendation

Based on your current setup, I'd suggest:

1. Kustomize - If you want to stay pure Kubernetes and have multiple environments
2. Terraform - If you want to manage everything in one place (since you're already using it for Cloudflare)
3. Helm - If you plan to share/reuse this setup or want the most flexibility

Given that you already have Terraform managing your Cloudflare tunnel, Terraform would be most consistent - you could manage both the DNS records and Kubernetes manifests with the same domain variable.

Would you like me to help you implement any of these approaches?
