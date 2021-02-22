# gRPC-using-golang
## gRPC  client &amp; server

compile proto files using
* protoc calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:.
* protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.



set GOPATH using:
* export GO_PATH=~/go
* export PATH=$PATH:/$GO_PATH/bin
