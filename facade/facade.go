package facade

// Provide quick registration and get methods

const (
	noRegister = " No Register"
)

const (
	compDb            = "_db"
	serviceGrpc       = "_grpc"
	compRedis         = "_redis"
	serviceHttpServer = "_http"
	compMutex         = "_mutex"
	compCache         = "_cache"
	compQueue         = "_queue"
)

func buildID(cType, id string) string {
	return id + cType
}
