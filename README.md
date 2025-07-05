create a scaler based on KEDA, the functionalities include:
- stockprice scaler
  - get stock price from yahoo finance
  - compare the price with the target price
    - if the price is higher than the target price, scale up the replica count
    - if the price is lower than the target price, scale down the replica count
- ingress scaler(TBD)
- rocketmq scaler(TBD)



How to fetch grpc API
```sh
# install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# fetch the grpc API
➜  keda-external-scaler grpcurl -plaintext localhost:6000 list 
externalscaler.ExternalScaler
grpc.reflection.v1.ServerReflection
grpc.reflection.v1alpha.ServerReflection
➜  keda-external-scaler grpcurl -plaintext localhost:6000 list externalscaler.ExternalScaler         
externalscaler.ExternalScaler.GetMetricSpec
externalscaler.ExternalScaler.GetMetrics
externalscaler.ExternalScaler.IsActive
externalscaler.ExternalScaler.StreamIsActive
➜  keda-external-scaler grpcurl -plaintext localhost:6000 describe externalscaler.ExternalScaler
externalscaler.ExternalScaler is a service:
service ExternalScaler {
  rpc GetMetricSpec ( .externalscaler.ScaledObjectRef ) returns ( .externalscaler.GetMetricSpecResponse );
  rpc GetMetrics ( .externalscaler.GetMetricsRequest ) returns ( .externalscaler.GetMetricsResponse );
  rpc IsActive ( .externalscaler.ScaledObjectRef ) returns ( .externalscaler.IsActiveResponse );
  rpc StreamIsActive ( .externalscaler.ScaledObjectRef ) returns ( stream .externalscaler.IsActiveResponse );
}
# example
grpcurl -plaintext -d '{
  "name": "stock-scaler",
  "namespace": "default",
  "scalerMetadata": {
    "symbol": "AAPL",
    "targetPrice": "150.00"
  }
}' localhost:6000 externalscaler.ExternalScaler.IsActive
```