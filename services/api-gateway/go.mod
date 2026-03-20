module github.com/rado31/rabbit/api-gateway

go 1.23

require (
	github.com/gin-gonic/gin           v1.12.0
	github.com/rado31/rabbit/proto     v0.0.0
	google.golang.org/grpc             v1.79.3
)

replace github.com/rado31/rabbit/proto => ../../proto
