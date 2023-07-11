variable "override_workspace_name" {
  description = "Override default workspace name derived from terraform"
  type        = string
  default     = "serverlessAndDuplicates"
}
variable "region" {
  description = ""
  type        = string
  default     = "eu-south-1"
}

variable "is_local" {
  description = "Used for sam local invoke"
  type        = bool
  default     = false
}

variable "calcOccurrences_version" {
  description = "Lambda version"
  type        = string
  default     = "latest"
}

variable "deploy_lambdas" {
  default = true
}
