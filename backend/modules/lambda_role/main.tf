resource "aws_iam_role" "for_lambda" {
  name               = format("%s-%s-%s", var.environment, var.solution_name, var.function_name)
  description        = "service role for ${var.function_name} lambda"
  assume_role_policy = var.assume_role_policy_document
}

resource "aws_iam_policy" "for_lambda" {
  name        = format("%s-%s-%s", var.environment, var.solution_name, var.function_name)
  description = "role policy for ${var.function_name} lambda"
  policy      = var.role_policy_document
}

resource "aws_iam_role_policy_attachment" "for_lambda" {
  role       = aws_iam_role.for_products_lambda.id
  policy_arn = aws_iam_policy.for_products_lambda.arn
}
