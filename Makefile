test:
	@echo "=============Running unit tests============="
	go test ./... -cover -coverprofile cover.out

lint:
	@echo "=============Linting============="
	golint -set_exit_status ./...
