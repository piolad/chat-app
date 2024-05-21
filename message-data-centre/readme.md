# This is normal for main.go in server folder to be red

To have files that are compiled from .proto (this should be only run to chehck if the synatx of the prorgam is correct, not red, otherwise they are created inside docker container)
```bash
protoc --go_out=. --go-grpc_out=. proto/service.proto
```