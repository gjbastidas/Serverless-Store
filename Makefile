export MAIN_PATH = $$PWD
export DEV_INFRA_PATH = backend/environments/dev
export AWS_PKG_PATH = backend/store_apis/pkg/aws

ddb-mock-generate:
	@ if [ -z "$$GOPATH" ]; then \
			echo "check for empty GOPATH environment variable" && exit 1; \
		fi && \
		MOCKGEN="$$GOPATH/bin/mockgen" && \
		if [ ! -f "$$MOCKGEN" ]; then \
			echo "check gomock installation" && exit 1; \
		fi && \
		cd "$$AWS_PKG_PATH" && \
		$$MOCKGEN -source=dynamodb.go -destination=mocks/dynamodb.go

tf-destroy-dev:
	@ cd "$$DEV_INFRA_PATH" && \
		terraform destroy

tf-apply-dev:
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