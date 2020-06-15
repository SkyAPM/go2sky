### GRPCReporterOption

`GRPCReporterOption` allows for functional options to adjust behaviour of a `gRPC` reporter to be created by `NewGRPCReporter`.

|    Function    | Describe |
| ---------- | --- |
| `reporter.WithLogger` |  setup logger for gRPC reporter |
| `reporter.WithCheckInterval` |  setup service and endpoint registry check interval |
| `reporter.WithInstanceProps` |  setup service instance properties eg: org=SkyAPM |
| `reporter.WithTransportCredentials` |  setup transport layer security |
| `reporter.WithAuthentication` |  used Authentication for gRPC |