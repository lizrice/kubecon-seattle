# Save yourselves

Here's a really great idea: run this.

```bash
kubectl apply -f https://raw.githubusercontent.com/lizrice/kubecon-seattle/master/save-yourselves.yaml?token=AAb_eGlPYRK4_K0pJdeGL7j_NT_dQodJks5cFcQVwA%3D%3D
```

## Access service account from inside a pod

```bash
TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
NAMESPACE=$(cat /var/run/secrets/kubernetes.io/serviceaccount/namespace)

curl -k -H "Authorization: Bearer $TOKEN" https://kubernetes.default/api/v1/
# gets a list of API resources

curl -k -H "Authorization: Bearer $TOKEN" https://kubernetes.default/api/v1/namespaces/default/pods
# We don't have access to everything with the default service account
```

## Admission control

Validating webhook admission control should be enabled by default in 1.13
Run webhook-certs.sh to generate the secret containing the certificates that the webhook server uses

## Opting out of service account tokens

Opt out of automounting API credentials for a service account:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: my-sa
automountServiceAccountToken: false
```

Opt out of automounting API credentials for a particular pod:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  serviceAccountName: my-sa
  automountServiceAccountToken: false
```

