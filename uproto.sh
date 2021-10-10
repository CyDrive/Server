protoc --proto_path=rpc --go_out=rpc --go_opt=paths=source_relative --go-grpc_out=rpc --go-grpc_opt=paths=source_relative manage.proto
protoc --proto_path=model --go_out=model --go_opt=paths=source_relative --csharp_out=model storage.proto account.proto
protoc --proto_path=consts --go_out=consts --go_opt=paths=source_relative --csharp_out=consts status_code.proto