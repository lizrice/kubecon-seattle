# Access service account from inside a pod

```bash
TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)

curl -k -H "Authorization: Bearer $TOKEN" https://kubernetes.default/api/v1/
# gets a list of API resources

curl -k -H "Authorization: Bearer $TOKEN" https://kubernetes.default/api/v1/namespaces/default/pods
# We wouldn't have access to everything with the default service account

# Install kubectl
curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
chmod +x kubectl
kubectl get pods
kubectl auth can-i create pods
```

## Delete this-is-fine

```bash
kubectl delete -f https://raw.githubusercontent.com/lizrice/kubecon-seattle/master/this-is-fine.yaml
```

[Next](./PAGE-3.md)
