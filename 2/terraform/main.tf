terraform {
  required_version = "~>1.5.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.63.0"
    }
  }
}

provider "aws" {
  region = var.region

  default_tags {
    tags = {
      team      = "amir"
      terraform = true
      env       = local.workspace
    }
  }
  ignore_tags {
    # Ignore tags added by "Resource Auto Tagger" tool
    keys = ["CreatedBy", "CreationDate", "IamRoleName"]
  }
}

locals {
  repository_prefix          = "cnapp-graph"
  workspace                  = lower(coalesce(var.override_workspace_name, terraform.workspace))
}

data "aws_caller_identity" "current" {}
