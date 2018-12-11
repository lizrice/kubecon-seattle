# Simpler with OPA

```bash
kubectl create configmap -n opa sa --from-file=opa/sa.rego

# Check there are no errors reported
kubectl describe configmap -n opa sa

# Try applying the YAML again
kubectl apply -f this-is-fine.yaml
```