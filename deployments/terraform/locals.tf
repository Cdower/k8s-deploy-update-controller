locals {
  controller_ns = "deploy-update-controller"
  deployments = {
    home-assistant-complexity = {
        namespace = "home",
        channel   = "https://github.com/home-assistant/core/releases/latest"
        container = "home-assistant"
        deploy    = "home-assistant-complexity"
    },
  }
  regcred_enable = (var.registry_user != "" && var.registry_pass != "") ? 1 : 0
}