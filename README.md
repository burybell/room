# room
Easy to expand message distribution component, can be used for chat rooms, and has achieved horizontal expansion of clusters with long connections to ws or sse using nats go

# Example
a sample chat room, you can run it directly

## start nats
```shell
docker run -d -p 4222:4222 nats:latest
```

## start chat room
```shell
cd ./examples/chat-room/ && go build && ./chat-room
```

## open browser
```shell
http://localhost:8080/index
```