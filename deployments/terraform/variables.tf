variable "image" {
  type = string
}

variable "ServiceAccountName" {
  type = string
  default = "deploy-update-controller"
}

variable "ControllerNamespace" {
  type = string
  default = "deploy-update-controller"
}

variable "registry_user" {
  type = string
  default = ""
}

variable "registry_pass" {
  type = string
  default = ""
}

variable "registry_server" {
  type = string
  default = "https://index.docker.io/v1/"
}