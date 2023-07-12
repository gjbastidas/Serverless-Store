####################
#     Database     #
####################

module "products_table" {
  source  = "terraform-aws-modules/dynamodb-table/aws"
  version = "3.3.0"

  name         = format("%s-%s-%s", var.environment, var.solution_name, "products")
  hash_key     = "id"
  range_key    = "dateModified"
  billing_mode = "PAY_PER_REQUEST"

  attributes = [
    {
      name = "id",
      type = "N"
    },
    {
      name = "dateModified",
      type = "S"
    }
  ]
}

####################
#   Permissions    #
####################

data "aws_iam_policy_document" "for_products_lambda" {
  statement {
    effect = "Allow"
    actions = [
      "dynamodb:PutItem",
    ]

    resources = [
      module.products_table.dynamodb_table_arn,
    ]
  }
}

module "role_for_products_lambda" {
  source = "../../modules/lambda_role"

  environment                 = var.environment
  solution_name               = var.solution_name
  function_name               = "products"
  assume_role_policy_document = data.aws_iam_policy_document.to_assume_lambda_service_role.json
  role_policy_document        = data.aws_iam_policy_document.for_products_lambda.json
}

####################
#    Functions     #
####################

module "products_lambda" {
  source = "../../modules/lambda"

  environment   = var.environment
  solution_name = var.solution_name
  role_id       = module.role_for_products_lambda.id
  function_name = "products"
  source_path   = "../../store_apis/cmd/lambdas/products"
}

####################
#   API Gateway    #
####################

module "api_gw" {
  source = "../../modules/api_gateway"

  environment   = var.environment
  solution_name = var.solution_name
}

################################
#   API Gateway Integrations   #
################################

module "products_lambda_integration" {
  source = "../../modules/api_gateway_lambda_integration"

  api_id            = module.api_gw.api_id
  api_execution_arn = module.api_gw.api_execution_arn
  integration_type  = "AWS_PROXY"
  integration_uri   = module.products_lambda.invoke_arn
  function_name     = module.products_lambda.function_name
  route_key         = "POST /products"
}
