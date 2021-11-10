protoc --proto_path=. -I"third_party/proto" \
    --dart_out=models \
    models/account.proto models/file_info.proto models/message.proto models/http_models.proto

protoc --proto_path=. -I"third_party/proto" \
    --dart_out=consts \
    consts/enums.proto

protoc --dart_out=. \
    google/protobuf/timestamp.proto