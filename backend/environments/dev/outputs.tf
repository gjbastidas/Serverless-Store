output "products_table_arn" {
  value = module.products_table.dynamodb_table_arn
}

output "products_lambda_arn" {
  value = module.products_lambda.function_arn
}

output "stage_invoke_url" {
  value = module.api_gw.stage_invoke_url
}
