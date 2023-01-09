# Description: List of variables which can be provided ar runtime to override the specified defaults 

variable "project_id" {
  description = "GCP Project ID"
  type        = string
  nullable    = false
}

variable "name" {
  description = "Base name to derive everythign else from"
  default     = "disco"
  type        = string
  nullable    = false
}

variable "location" {
  description = "Deployment location"
  default     = "us-west1"
  type        = string
  nullable    = false
}

variable "git_repo" {
  description = "GitHub Repo"
  type        = string
  nullable    = false
}

variable "server_img" {
  description = "Image URI"
  default     = "us-west1-docker.pkg.dev/cloudy-demos/disco/disco"
  type        = string
  nullable    = false
}

variable "disco_schedule" {
  description = "Cron for disco service invocation"
  default     = "30 */5 * * *"
  type        = string
  nullable    = false
}

variable "runtime_only" {
  description = "Whether or not deploy the development resoruces"
  default     = true
  type        = bool
}
