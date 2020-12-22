# biubiu
devops system server built by golang


### generate grpc file
```shell
protoc -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I. \
--go_out=. --go_opt=paths=source_relative \
--go-grpc_out=. --go-grpc_opt=paths=source_relative \
user.proto 
```

### generate grpc-gateway file
```shell
protoc -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I. \
--grpc-gateway_out . \
--grpc-gateway_opt logtostderr=true \
--grpc-gateway_opt paths=source_relative \
user.proto 
```

### generate grpc-validators
```shell
protoc -I${GOPATH}/src \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
-I${GOPATH}/src/github.com/mwitkow/go-proto-validators -I. \
--go_out=. --go_opt=paths=source_relative \
--govalidators_out=gogoimport=true:. \
auth.proto 
```