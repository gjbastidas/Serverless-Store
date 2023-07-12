data "aws_iam_role" "role" {
  name = var.role_id
}

data "aws_iam_policy_document" "for_logs_access" {
  statement {
    effect = "Allow"
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = [
      aws_cloudwatch_log_group.for_lambda.arn,
    ]
  }
}
