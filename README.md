# Go Deployment Updater
Run as a k8s cron job checking open source projects for new releases.
A cheap and lazy generic knock-off inspired by https://github.com/rancher/system-upgrade-controller for keeping random deployments up to date.

## Configs
All configs set as envvar 
* Channel: The value for `channel` is assumed to be a URL that returns HTTP 302 with the last path element of the value returned in the Location header assumed to be an image tag (after munging "+" to "-").
  * channel: https://github.com/home-assistant/core/releases/latest
* Namespace: the namespace of the deployment
  * namespace: home-assistant
* Deployment: the name of the deployment to update
  * deployment: home-assistant
* serviceAccountName: The service account for the pod to use. As with normal pods, if not specified the `deploy-update-controller` service account from the namespace will be assigned.
  * See https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
* ControllerNamespace: Namespace for deploy updater to run in 
  * default `ControllerNamespace: deploy-update-controller`

## Multi Stage Builds
Uses buildx using docs and reccomendations from  https://github.com/docker/buildx#building-multi-platform-images and https://github.com/tonistiigi/xx
TODO: Checkout buildx bake https://docs.docker.com/buildx/working-with-buildx/

## Structure
Porject layout based on [golang-standards/project-layout](https://github.com/golang-standards/project-layout)

## Secrets
secret_envfile sets default empty strings as defaults for the registry credentials, currently dockerhub only, and then imports secret.env to override them.
If blank no secret will be created.