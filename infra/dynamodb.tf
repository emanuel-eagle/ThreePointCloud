resource "aws_dynamodb_table" "dynamodb-table" {
  name           = "player-database"
  billing_mode   = "PROVISIONED"
  read_capacity = 1000
  write_capacity = 1000
  hash_key       = "player-id"

  attribute {
    name = "player-id"
    type = "S"
  }
}