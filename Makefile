STACK = go-log-trace
PROFILE = mine

.PHONY: build deploy validate destroy

deploy:
	sam build
	AWS_PROFILE=$(PROFILE) sam deploy --stack-name $(STACK) \
	--resolve-s3 --no-confirm-changeset --no-fail-on-empty-changeset --capabilities CAPABILITY_IAM

validate:
	AWS_PROFILE=$(PROFILE) aws cloudformation validate-template --template-body file://template.yml

destroy:
	AWS_PROFILE=$(PROFILE) aws cloudformation delete-stack --stack-name $(STACK)

clean:
	rm -rf main gin-bin

.PHONY: build-BMainFunction
build-BMainFunction:
	cd b && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION}" -o ${ARTIFACTS_DIR}/bmain

.PHONY: build-AMainFunction
build-AMainFunction:
	cd a && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION}" -o ${ARTIFACTS_DIR}/amain
