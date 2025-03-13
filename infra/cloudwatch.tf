resource "aws_cloudwatch_log_group" "career_stats_cloudwatch_group" {
  name              = "/aws/lambda/${var.lambda_function_name_career_stats}"
  retention_in_days = 14
}