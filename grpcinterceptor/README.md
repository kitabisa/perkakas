# Package grpcinterceptor

This package contains gRPC Interceptor (middleware).

## RequestID
RequestID interceptor is an interceptor to capture RequestID from the request metadata and inject it into context.

## Logger
Logger interceptor is an interceptor that will help logging grpc server

## How to use the interceptor

### Single interceptor
```go
func main(){
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", 50051))
	if err != nil {
		opt.Logger.Error(fmt.Sprintf("failed to listen %s:%d", host, port), err, nil)
	}

    // use the interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			requestid.UnaryServerInterceptor,
		),
	)

    // initialize grpc handler
	grpcHandler := grpcHandler.NewFlagHandler(opt)

	pb.RegisterFlagServer(grpcServer, grpcHandler)
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	opt.Logger.Info(fmt.Sprintf("GRPC serve at %s:%d", host, port), nil)

	grpcServer.Serve(lis)
}
```

### Chaining multiple interceptor
```go
func main(){
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", 50051))
	if err != nil {
		opt.Logger.Error(fmt.Sprintf("failed to listen %s:%d", host, port), err, nil)
	}

    // chaining multiple interceptor
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			requestid.UnaryServerInterceptor,
			logger.UnaryServerInterceptor,
		),
	)

    // initialize grpc handler
	grpcHandler := grpcHandler.NewFlagHandler(opt)

	pb.RegisterFlagServer(grpcServer, grpcHandler)
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	opt.Logger.Info(fmt.Sprintf("GRPC serve at %s:%d", host, port), nil)

	grpcServer.Serve(lis)
}
```