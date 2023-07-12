resource "aws_apigatewayv2_api" "main" {
  name          = format("%s-%s", var.environment, var.solution_name)
  protocol_type = "HTTP"
}

resource "aws_cloudwatch_log_group" "log" {
  name              = "/aws/api-gw/${aws_apigatewayv2_api.main.name}"
  retention_in_days = var.cloudwatch_log_retention_in_days
}

resource "aws_apigatewayv2_stage" "stage" {
  api_id = aws_apigatewayv2_api.main.id

  name        = var.environment
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.log.arn

    format = jsonencode({
      requestId               = "$context.requestId"
      sourceIp                = "$context.identity.sourceIp"
      requestTime             = "$context.requestTime"
      protocol                = "$context.protocol"
      httpMethod              = "$context.httpMethod"
      resourcePath            = "$context.resourcePath"
      routeKey                = "$context.routeKey"
      status                  = "$context.status"
      responseLength          = "$context.responseLength"
      integrationErrorMessage = "$context.integrationErrorMessage"
      }
    )
  }
}
