# With OPA

Follow [OPA instructions](https://www.openpolicyagent.org/docs/kubernetes-admission-control.html) to install OPA and get it up and running

## Install rules to prevent service account

```bash
kubectl create configmap -n opa sa --from-file=opa/sa.rego

# Check there are no errors reported
kubectl describe configmap -n opa sa

# Try applying the YAML again
kubectl apply -f this-is-fine.yaml
```
