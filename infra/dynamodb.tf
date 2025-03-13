resource "aws_dynamodb_table" "dynamodb-table" {
  name           = "player-database"
  billing_mode   = "PROVISIONED"
  hash_key       = "player-id"

  attribute {
    name = "player-id"
    type = "S"
  }
}