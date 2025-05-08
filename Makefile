generate:
	cd proto && buf generate
update:
	cd proto && buf mod update

run:
	go run cmd/grpc-gw/main.go