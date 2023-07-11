resource "aws_dynamodb_table" "table" {
  name           = local.workspace
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "Id"
  range_key      = "Timestamp"

  attribute {
    name = "Id"
    type = "S"
  }

  attribute {
    name = "Timestamp"
    type = "N"
  }
}
