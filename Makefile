export MAIN_PATH = $$PWD
export DEV_INFRA_PATH = backend/environments/dev

tf-destroy-dev:
	@ cd "$$DEV_INFRA_PATH" && \
		terraform destroy

	@ cd "$$DEV_INFRA_PATH" && \
		terraform apply plan.out

tf-plan-dev: validate-envs
	@ cd "$$DEV_INFRA_PATH" && \
		terraform init && \
		terraform validate && \
		terraform plan -out=plan.out

validate-envs:
	@ if [ -z "$$AWS_PROFILE" ]; then \
			echo "check for empty AWS_PROFILE environment variable" && exit 1; \
		fi