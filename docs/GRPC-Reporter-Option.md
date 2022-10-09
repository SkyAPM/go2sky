### GRPCReporterOption

`GRPCReporterOption` allows for functional options to adjust behaviour of a `gRPC` reporter to be created by `NewGRPCReporter`.

| Function                            | Describe                                                                                         |
|-------------------------------------|--------------------------------------------------------------------------------------------------|
| `reporter.WithLog`                  | setup log for gRPC reporter                                                                      |
| `reporter.WithCheckInterval`        | setup service and endpoint registry check interval                                               |
| `reporter.WithMaxSendQueueSize`     | setup send span queue buffer length                                                              |
| `reporter.WithInstanceProps`        | setup service instance properties eg: org=SkyAPM                                                 |
| `reporter.WithTransportCredentials` | setup transport layer security                                                                   |
| `reporter.WithAuthentication`       | used Authentication for gRPC                                                                     |
| `reporter.WithCDS`                  | setup CDS service                                                                                |
| `reporter.WithLayer`                | setup layer                                                                                      |
| `reporter.WithFAASLayer`            | setup layer to FAAS                                                                              |
| `reporter.WithProcessLabels`        | setup labels bind to process                                                                     |
| `reporter.WithProcessStatusHook`    | setup is enabled the process status                                                              |
| `reporter.WithMeterCollectPeriod`   | setup meter collection interval, if input is <= 0, go2sky will not collect meter, default is 15s |
