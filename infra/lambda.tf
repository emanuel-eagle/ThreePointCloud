variable package_type {
    description = "Name of the lambda function"
    type = string
    default = "Image"
}

variable lambda_function_name {
    default = "threepointcloud-playerlist-collection"
}

variable lambda_function_name_career_stats {
    default = "threepointcloud-careerstats-collection"
}

variable lambda_function_name_career_stats_coordinator {
    default = "threepointcloud-careerstats-coordinator"
}

variable lambda_timeout {
    default = 120
}

variable career_stats_lambda_timeout {
    default = 900
}

variable career_stats_lambda_memory_size {
    default = 1024
}

variable memory_size {
    default = 512
}

variable lambda_image_tag {
    type = string
}

variable career_stats_lambda_image_tag {
    type = string
}

variable lambda_image_tag_careerStatsCoordinatorLambda {
    type = string
}

variable chunk_size {
    type = string
    default = "200"
}

resource "aws_lambda_function" "threepointcloud_playerlist_collection" {
  function_name = var.lambda_function_name
  role          = aws_iam_role.iam_for_lambda.arn
  image_uri     = "${aws_ecr_repository.three-point-cloud-player-list-container-repository.repository_url}:${var.lambda_image_tag}"
  package_type = var.package_type
  timeout = var.lambda_timeout
  memory_size = var.memory_size
}

resource "aws_lambda_function" "threepointcloud_careerstats_collection" {
  function_name = var.lambda_function_name_career_stats
  role          = aws_iam_role.iam_for_lambda.arn
  image_uri     = "${aws_ecr_repository.three-point-cloud-career-stats-container-repository.repository_url}:${var.career_stats_lambda_image_tag}"
  package_type = var.package_type
  timeout = var.career_stats_lambda_timeout
  memory_size = var.career_stats_lambda_memory_size
}

resource "aws_lambda_function" "threepointcloud_careerstats_coordinator_collection" {
  function_name = var.lambda_function_name_career_stats_coordinator
  role          = aws_iam_role.iam_for_lambda.arn
  image_uri     = "${aws_ecr_repository.three-point-cloud-career-stats-coordinator-container-repository.repository_url}:${var.lambda_image_tag_careerStatsCoordinatorLambda}"
  package_type = var.package_type
  timeout = var.lambda_timeout
  memory_size = var.memory_size
  environment {
    variables = {
        TABLE_NAME = aws_dynamodb_table.dynamodb-table.name
        HASH_KEY = aws_dynamodb_table.dynamodb-table.hash_key
        CAREER_STATS_LAMBDA = aws_lambda_function.threepointcloud_careerstats_collection.arn
        CHUNK_SIZE = var.chunk_size
    }
  }
}