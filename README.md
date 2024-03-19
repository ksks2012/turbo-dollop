# Introduce

- Implements an easy limiter to control and limit actions or operations in application.
- In order to study system design. 

# Command

## HTTP with Redis

```
go run ./example/redis_http/main.go
```

## HTTP with memory
```
go run ./example/memory_http/main.go
```

## Test

- memory
    - test
        ```
        go test ./drivers/storage/memory/
        ```
    - benchmark
        ```
        go test -bench -v ./drivers/storage/memory/
        ```

- redis
    - test
        ```
        go test ./drivers/storage/redis/
        ```
    - benchmark
        ```
        go test -bench -v ./drivers/storage/redis/
        ```

# TODO

## Storage

- [x] Memeory: with go-cache
- [ ] Redis

## Middleware

- [ ] HTTP
- [ ] Fasthttp
- [ ] gin

## Algorithm

- [ ] Token Bucket
- [ ] Leaking Bucket
- [ ] Fixed Window
- [ ] Sliding Window
- [ ] Sliding Window Counter