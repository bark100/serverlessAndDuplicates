variable "override_workspace_name" {
  description = "Override default workspace name derived from terraform"
  type        = string
  default     = "serverlessAndDuplicates"
}
variable "region" {
  description = "AWS region."
  type        = string
  default     = "eu-south-1" // Closest to Israel
}

variable "deploy_lambdas" {
  default = true
}

variable "deploy_apigw" {
  default = true
}

variable "deploy_dynamodb" {
  default = true
}
