module github.com/safar/microservices-demo/services/cart

go 1.26.0

require (
	github.com/lib/pq v1.10.9
	github.com/safar/microservices-demo/proto v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.79.1
)

require (
	golang.org/x/net v0.50.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260217215200-42d3e9bedb6d // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/safar/microservices-demo/proto => ../../proto
