# Access service account from inside a pod

```bash
TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
NAMESPACE=$(cat /var/run/secrets/kubernetes.io/serviceaccount/namespace)

curl -k -H "Authorization: Bearer $TOKEN" https://kubernetes.default/api/v1/
# gets a list of API resources

curl -k -H "Authorization: Bearer $TOKEN" https://kubernetes.default/api/v1/namespaces/default/pods
# We don't have access to everything with the default service account
```

## Service account version

```bash
kubectl apply -f https://raw.githubusercontent.com/lizrice/kubecon-seattle/master/this-is-even-better.yaml?token=AAb_eGlPYRK4_K0pJdeGL7j_NT_dQodJks5cFcQVwA%3D%3D
```


[Next](./PAGE-3.md)
