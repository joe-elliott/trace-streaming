go run ./cmd/server --config ./cmd/server/config.yaml
go run ./cmd/client

docker run \
  --rm \
  --env JAEGER_ENDPOINT=http://192.168.1.124:14268/api/traces \
  -p8080-8083:8080-8083 \
  jaegertracing/example-hotrod:latest \
  all