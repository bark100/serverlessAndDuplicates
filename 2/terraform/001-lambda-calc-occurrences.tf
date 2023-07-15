locals {
  calcOccurrences_name = "${local.workspace}-calcOccurrences"
}

module "calcOccurrences" {
  count  = var.deploy_lambdas ? 1 : 0
  source = "terraform-aws-modules/lambda/aws"

  function_name = local.calcOccurrences_name
  description   = "todo"
  handler       = "main"
  runtime       = "go1.x"
  timeout       = 5
  memory_size   = 128

  source_path = [
    {
      path = "../lambdas/calcOccurrences",
      commands = [
        "GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o .build/main",
        ":zip .build/main ."
      ],
      patterns = [".build/main"]
    }
  ]

  environment_variables = {
    DYNAMODB_TABLE = one(aws_dynamodb_table.table[*].name)
  }

  attach_policies    = true
  number_of_policies = 1
  policies = [
    "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess" // TODO: reduce this to the table level.
  ]
}
