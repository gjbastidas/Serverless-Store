resource "aws_apigatewayv2_route" "route" {
  api_id    = var.api_id
  route_key = var.route_key

  target = "integrations/${var.integration_id}"
}
