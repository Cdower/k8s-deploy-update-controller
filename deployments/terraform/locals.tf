locals {
  controller_ns = "deploy-update-controller"
  deployments = {
    home-assistant-complexity = {
        namespace = "home",
        channel   = "https://github.com/home-assistant/core/releases/latest",
        container = "home-assistant",
        deploy    = "home-assistant-complexity",
    },
    "ombi-complexity" = {
      namespace = "home",
      channel   = "https://github.com/linuxserver/docker-ombi/releases/tag/v4.3.3-ls111",
      container = "ombi",
      deploy    = "ombi-complexity",
    },
    "sonar-complexity" = {
      namespace = "home",
      channel   = "https://github.com/linuxserver/docker-sonarr/releases/latest",
      container = "sonar",
      deploy    = "sonar-complexity",
    }
  }
  regcred_enable = (var.registry_user != "" && var.registry_pass != "") ? 1 : 0
  regcred_secret_namespaces = local.regcred_enable == 1 ? flatten([local.controller_ns, var.regcred_add_ns]) : []
}
