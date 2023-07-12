resource "aws_cloudwatch_log_group" "for_lambda" {
  name              = "/aws/lambda/${local.function_name}"
  retention_in_days = var.cloudwatch_log_retention_in_days
}

resource "aws_iam_policy" "for_logs_access" {
  name        = format("%s-%s-%s-%s", var.environment, var.solution_name, var.function_name, "for-log-access")
  description = "role policy for ${var.function_name} lambda for cloudwatch logs access"
  policy      = data.aws_iam_policy_document.for_logs_access.json
}

resource "aws_iam_role_policy_attachment" "for_logs_access" {
  role       = var.role_id
  policy_arn = aws_iam_policy.for_logs_access.arn
}

locals {
  binary_name   = "app"
  src_path      = "${var.source_path}/main.go"
  bin_path      = "${var.source_path}/bin/${local.binary_name}"
  archive_path  = "${var.source_path}/archive/${local.binary_name}.${var.archive_type}"
  function_name = format("%s-%s-%s", var.environment, var.solution_name, var.function_name)
}

resource "null_resource" "go_build" {
  provisioner "local-exec" {
    command = "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ${local.bin_path} ${local.src_path}"
  }
}

data "archive_file" "zip" {
  depends_on = [null_resource.go_build]

  type        = var.archive_type
  source_file = local.bin_path
  output_path = local.archive_path
}

resource "aws_lambda_function" "function" {
  function_name    = local.function_name
  description      = "${var.function_name} lambda"
  role             = data.aws_iam_role.role.arn
  handler          = local.binary_name
  runtime          = var.runtime
  memory_size      = var.memory_size
  timeout          = var.timeout
  filename         = data.archive_file.zip.output_path
  source_code_hash = data.archive_file.zip.output_base64sha256
}
