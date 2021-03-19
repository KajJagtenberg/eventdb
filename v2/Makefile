run:
	go run cmd/server/main.go

pb:
	protoc --proto_path=proto --go_out=plugins=grpc:store --go_opt=paths=source_relative proto/store.proto