resource "aws_dynamodb_table" "table" {
  count = var.deploy_dynamodb ? 1 : 0

  name         = local.workspace
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "Id"

  attribute {
    name = "Id"
    type = "S"
  }
}
