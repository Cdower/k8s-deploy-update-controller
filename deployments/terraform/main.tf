resource "kubernetes_namespace" "controller" {
  metadata {
    name = local.controller_ns
  }
}

resource "kubernetes_service_account" "controller" {
  for_each = local.deployments
  metadata {
    name = "deploy-update-controller-${each.key}"
    namespace = local.controller_ns
  }
}

resource "kubernetes_role" "controller_target_role" {
    for_each = local.deployments
    metadata {
      name = "deploy-update-controller-${each.key}"
      namespace = each.value.namespace
    }
    rule {
      api_groups  = ["extensions", "apps"]
      resources   = ["deployments"]
      verbs       = ["get", "list", "watch", "create", "update", "patch", "delete"]
    }
    rule {
      api_groups  = [""]
      resources   = ["pods"]
      verbs       = ["get", "list", "watch"]
    }
}

resource "kubernetes_role_binding" "controller" {
  for_each = local.deployments
  metadata {
    name = "deploy-update-controller-${each.key}"
    namespace = each.value.namespace
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "Role"
    name      = "deploy-update-controller-${each.key}"
  }
  subject {
    kind      = "ServiceAccount"
    name      = "deploy-update-controller-${each.key}"
    namespace = local.controller_ns
  }
}

resource "kubernetes_cron_job" "update" {
  for_each = local.deployments
  metadata {
    name = "deploy-update-controller-${each.key}"
    namespace = local.controller_ns
  }
  spec {
    concurrency_policy            = "Replace"
    failed_jobs_history_limit     = 5
    schedule                      = "15 4 * * *"
    starting_deadline_seconds     = 30
    successful_jobs_history_limit = 10
    job_template {
      metadata {}
      spec {
        backoff_limit              = 2
        ttl_seconds_after_finished = 900
        template {
          metadata {}
          spec {
            restart_policy = "OnFailure"
            image_pull_secrets {
              name = kubernetes_secret.regcred[local.controller_ns].metadata.0.name
            }
            service_account_name = "deploy-update-controller-${each.key}"
            container {
              name    = "deploy-update-controller-${each.key}"
              image   = var.image
              image_pull_policy = "Always"
              command = ["./deploy-update-controller"]
              env {
                name  = "IN_CLUSTER"
                value = "true"
              }
              env {
                name  = "CONTAINER"
                value = each.key
              }
              env {
                name  = "DEPLOYMENT"
                value = each.value["deploy"]
              } 
              env {
                name  = "NAMESPACE"
                value = each.value["namespace"]
              } 
              env {
                name  = "CHANNEL"
                value = each.value["channel"]
              }
              env {
                name  = "SERVICE_ACCOUNT_NAME"
                value = var.ServiceAccountName
              }
              env {
                name  = "CONTROLLER_NAMESPACE"
                value = var.ControllerNamespace
              }
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_secret" "regcred" {
  for_each = toset(local.regcred_secret_namespaces)
  metadata {
    name = "regcred"
    namespace = each.key
  }

  data = {
    ".dockerconfigjson" = <<DOCKER
{
  "auths": {
    "${var.registry_server}": {
      "auth": "${base64encode("${var.registry_user}:${var.registry_pass}")}"
    }
  }
}
DOCKER
  }

  type = "kubernetes.io/dockerconfigjson"
}
