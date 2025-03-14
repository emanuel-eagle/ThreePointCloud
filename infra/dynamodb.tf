resource "aws_dynamodb_table" "dynamodb-table" {
  name           = "celtics-player-database"
  billing_mode   = "PROVISIONED"
  read_capacity = 1000
  write_capacity = 1000
  hash_key       = "player-id"

  attribute {
    name = "player-id"
    type = "S"
  }
}

resource "aws_dynamodb_table" "dynamodb-table-career-data" {
  name           = "career-stats-database"
  billing_mode   = "PROVISIONED"
  read_capacity = 1000
  write_capacity = 1000
  hash_key       = "player-id"

  attribute {
    name = "player-id"
    type = "S"
  }
}