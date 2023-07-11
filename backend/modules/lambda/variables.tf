variable "solution_name" {}
variable "environment" {}

variable "archive_type" {
  default = "zip"
}

variable "function_name" {}
variable "role_id" {}

variable "source_path" {
  description = "relative path to the function source code"
}

variable "memory_size" {
  default = 128
}

variable "timeout" {
  default = 5
}

variable "runtime" {
  default = "go1.x"
}

variable "cloudwatch_log_retention_in_days" {
  default = 7
}
