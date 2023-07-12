output "api_id" {
  value = aws_apigatewayv2_api.main.id
}

output "api_execution_arn" {
  value = aws_apigatewayv2_api.main.execution_arn
}

output "stage_invoke_url" {
  value = aws_apigatewayv2_stage.stage.invoke_url
}
