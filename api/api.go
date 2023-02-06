//go:generate protoc cpu.proto --go_out=../api --proto_path=../api/proto/
//go:generate protoc load-average.proto --go_out=../api --proto_path=../api/proto/
//go:generate protoc load-disk.proto --go_out=../api --proto_path=../api/proto/
//go:generate protoc dfsize.proto --go_out=../api --proto_path=../api/proto/
//go:generate protoc dfinode.proto --go_out=../api --proto_path=../api/proto/
//go:generate protoc stats.proto --go-grpc_out=../api --go_out=../api --proto_path=../api/proto/

package api
