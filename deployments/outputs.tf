output "name" {
  secret_metadata = kubernetes_secret.regcred.metadata
}