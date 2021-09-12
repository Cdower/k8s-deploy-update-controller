output "regcred_name" {
  value = kubernetes_secret.regcred[local.controller_ns].metadata.0.name
}

output "regcred_enabled" {
  value = local.regcred_enable
}