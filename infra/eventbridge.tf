resource "aws_cloudwatch_event_rule" "weekly_threpointcloud_player_update" {
  name                = "weekly-playerlist-lambda-runner"
  description         = "Trigger player data update once a week"
  schedule_expression = "cron(0 0 ? * SUN *)"
}

resource "aws_cloudwatch_event_target" "run_lambda_weekly" {
  rule      = aws_cloudwatch_event_rule.weekly_threpointcloud_player_update.name
  target_id = "threepointcloud-playerlist-lambda"
  arn       = aws_lambda_function.threepointcloud_playerlist_collection.arn
}

resource "aws_lambda_permission" "allow_cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.threepointcloud_playerlist_collection.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.weekly_threpointcloud_player_update.arn
}