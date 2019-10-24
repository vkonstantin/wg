**Install**

```
make build
make test
make docker
docker run --name todo -p 8080:8080 -d wg/todo-image
```

**Try**

```
# create new user
curl -i -d '{"requestID":"r1"}' -H "Content-Type: application/json" -X POST http://localhost:8080/user
# list of TODOs
curl -i -H "token: {\"id\":1}" -X GET http://localhost:8080/todo
# Add TODO
curl -i -d '{"requestID":"r2", "text":"DOTO"}' -H "token: {\"id\":1}" -H "Content-Type: application/json" -X POST http://localhost:8080/todo
# list of TODOs
curl -i -H "token: {\"id\":1}" -X GET http://localhost:8080/todo
```

**Structure**

1. common - it is a package that should be shared between all other services.
2. message - this package for communication messages. In case of REST it is a request and response objects, that should be serialized as JSON. But it also can be protobuf, for example. In this case, it is more convenient to keep this package separate to easy sharing it between different services and to describe API and documentation.
3. controller - package for business logic
4. model - internal entities
5. server - types of external communications of service, aka REST, WS, gRPC.
6. storage - types of storages, aka Memory, some databases...

**Remarks**
1. JSON serialization is better to do using precompiled code, like easyjson library do.
2. I decided to do deduplication of request to do the test task more interesting. That is a simple implementation of this but it is very useful. It can help with client retries, saga-patterns and saces like this...
3. I decided to do a garbage collection function to keep fit the heap in case of a mass delete of elements.
