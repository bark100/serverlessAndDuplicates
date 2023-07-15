resource "aws_iam_role" "lambda_execution" {
  name = "${local.workspace}-apigw-lambda-execution"

  assume_role_policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Sid" : "",
        "Effect" : "Allow",
        "Principal" : {
          "Service" : ["apigateway.amazonaws.com"]
        },
        "Action" : "sts:AssumeRole"
      }
    ]
  })
  inline_policy {
    name = "LambdaAccess"
    policy = jsonencode(
      {
        "Version" : "2012-10-17",
        "Statement" : [
          {
            "Action" : "lambda:*",
            "Resource" : "*",
            "Effect" : "Allow"
          }
        ]
      }
    )
  }
}

resource "aws_cloudwatch_log_group" "this" {
  name = "${local.workspace}-apigw-access-log"
}

module "api_gateway" {
  source  = "terraform-aws-modules/apigateway-v2/aws"
  version = "~> 2.2.2"
  count   = var.deploy_apigw ? 1 : 0

  name                   = local.workspace
  protocol_type          = "HTTP"
  create_api_domain_name = false

  default_route_settings = {
    detailed_metrics_enabled = false
    throttling_burst_limit   = 1250
    throttling_rate_limit    = 2500
  }

  default_stage_access_log_destination_arn = aws_cloudwatch_log_group.this.arn
  default_stage_access_log_format = jsonencode(
    {
      caller           = "$context.identity.caller"
      httpMethod       = "$context.httpMethod"
      ip               = "$context.identity.sourceIp"
      protocol         = "$context.protocol"
      requestId        = "$context.requestId"
      requestTime      = "$context.requestTime"
      resourcePath     = "$context.resourcePath"
      responseLength   = "$context.responseLength"
      status           = "$context.status"
      user             = "$context.identity.user"
      userAgent        = "$context.identity.userAgent"
      errorMessage     = "$context.error.message"
      integrationError = "$context.integration.error"
    }
  )

  body = templatefile("${path.module}/files/openapi.yaml.tftpl", {
    calc_function_arn         = one(module.calcOccurrences[*].lambda_function_arn)
    getResult_function_arn    = one(module.getResult[*].lambda_function_arn)
    name                      = local.workspace
    version                   = "0.0.1"
    lambda_execution_role_arn = aws_iam_role.lambda_execution.arn
  })
}
