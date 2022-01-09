locals {
  controller_ns = "deploy-update-controller"
  deployments = {
    home-assistant-complexity = {
        namespace = "home",
        channel   = "https://github.com/home-assistant/core/releases/latest",
        container = "home-assistant",
        deploy    = "home-assistant-complexity",
        schedule  = "15 4 * * *",
    },
    # "ombi-complexity" = {
    #   namespace = "home",
    #   channel   = "https://github.com/linuxserver/docker-ombi/releases/latest",
    #   container = "ombi",
    #   deploy    = "ombi-complexity",
    #   schedule  = "15 4 * * *",
    # },
    "sonar-complexity" = {
      namespace = "home",
      channel   = "https://github.com/linuxserver/docker-sonarr/releases/latest",
      container = "sonar",
      deploy    = "sonar-complexity",
      schedule  = "15 4 * * *",
    }
    "transmission-complexity" = {
      namespace = "home",
      channel   = "https://github.com/bubuntux/nordvpn/releases/latest",
      container = "vpn",
      deploy    = "transmission-complexity",
      schedule  = "15 4 * * *",
    }
    "transmission-complexity2" = {
      namespace = "home",
      channel   = "https://github.com/linuxserver/docker-transmission/releases/latest",
      container = "transmission",
      deploy    = "transmission-complexity",
      schedule  = "15 5 * * *",
    }
  }
  regcred_enable = (var.registry_user != "" && var.registry_pass != "") ? 1 : 0
  regcred_secret_namespaces = local.regcred_enable == 1 ? flatten([local.controller_ns, var.regcred_add_ns]) : []
}
