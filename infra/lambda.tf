variable package_type {
    description = "Name of the lambda function"
    type = string
    default = "Image"
}

variable lambda_function_name {
    default = "threepointcloud-playerlist-collection"
}

variable lambda_timeout {
    default = 300
}

variable memory_size {
    default = 1024
}

variable lambda_image_tag {
    type = string
}

resource "aws_lambda_function" "weather_alerts_lambda_function" {
  function_name = var.lambda_function_name
  role          = aws_iam_role.iam_for_lambda.arn
  image_uri     = "${aws_ecr_repository.three-point-cloud-player-list-container-repository.repository_url}:${var.lambda_image_tag}"
  package_type = var.package_type
  timeout = var.lambda_timeout
  memory_size = var.memory_size
}