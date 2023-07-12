resource "aws_apigatewayv2_integration" "lambda_handler" {
  api_id = var.api_id

  integration_type = var.integration_type
  integration_uri  = var.integration_uri
}

resource "aws_apigatewayv2_route" "handler" {
  api_id    = var.api_id
  route_key = var.route_key

  target = "integrations/${aws_apigatewayv2_integration.lambda_handler.id}"
}

resource "aws_lambda_permission" "api_gw" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = var.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${var.api_execution_arn}/*/*"
}
