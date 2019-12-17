blerg

todo:
- support multiple clients
- propagate context/trace
- add support for other streams
  - trace
  - traceheader
- web interface
- some kind of filtering?
- metrics
- streamers never removed from the slice in stream processor :)
- gracefully shutdown batch polling
- add opentelemetry tracing
- use uber atomic for shared maps.  note that stuff in stream.go is not concurrency safe.
- build query frontend.  query frontend only hits enough collectors to satisfy the requested rate limit