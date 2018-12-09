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

