package kubernetes.admission

import data.kubernetes.namespaces

deny[msg] {
    input.request.kind.kind = "RoleBinding"
    input.request.operation = "CREATE"
    rolename = input.request.object.roleRef.name
    rolename = "admin"
    msg = sprintf("not letting you create role binding with role %s", [rolename])
}

