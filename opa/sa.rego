package kubernetes.admission

import data.kubernetes.namespaces

deny[msg] {
    input.request.kind.kind = "ServiceAccount"
    input.request.operation = "CREATE"
    sa = input.request.object.metadata.name
    msg = "not letting you create a service account"
}