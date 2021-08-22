provider "kubernetes" {
  config_path    = "~/.kube/config"
  config_context = "default"
}

terraform {
  backend "local" {
    path = "./terraform.tfstate"
  }
}
