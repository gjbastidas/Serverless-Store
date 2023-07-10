######################
##     Database     ##
######################

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
