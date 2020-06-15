### GRPCReporterOption

`GRPCReporterOption` allows for functional options to adjust behaviour of a `gRPC` reporter to be created by `NewGRPCReporter`.

|    Function    | Describe |
| ---------- | --- |
| `grpc.WithLogger` |  setup logger for gRPC reporter |
| `grpc.WithCheckInterval` |  setup service and endpoint registry check interval |
| `grpc.WithInstanceProps` |  setup service instance properties eg: org=SkyAPM |
| `grpc.WithTransportCredentials` |  setup transport layer security |
| `grpc.WithAuthentication` |  used Authentication for gRPC |