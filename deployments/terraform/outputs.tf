output "regcred_name" {
  value = kubernetes_secret.regcred.0.metadata.0.name
}

output "regcred_enabled" {
  value = local.regcred_enable
}