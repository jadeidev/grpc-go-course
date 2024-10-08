# Generate proto files
cd into the greet-buf folder
`buf generate`

to start server:
`go run greet_server/server.go`

to start client:
`go run greet_client/client.go`


# How to setup repo for only protobuf files
proto buf repo structure (only stores proto files):
```
https://github.com/jadeidev/my-proto-repo
├── README.md
├── buf.gen.yaml
├── buf.lock
├── buf.yaml
└── proto
    ├── service1
    │   └── v1
    │       ├── service1_a.proto
    │       ├── service1_b.proto
    │       └── service1_c.proto
    └── service2
        └── v1
            ├── service2_a.proto
            ├── service2_b.proto
            └── service2_c.proto
```

You can create or update a `buf.lock` file for your module by running the buf dep update command

## Files
```yaml
# buf.yaml
version: v2
lint:
  use:
    - DEFAULT
  rpc_allow_google_protobuf_empty_requests: true
  rpc_allow_google_protobuf_empty_responses: true

breaking:
  use:
    - FILE

modules:
  - path: proto

deps:
  - buf.build/bufbuild/protovalidate
```

```yaml
# buf.gen.yaml
version: v2
plugins:
  - remote: buf.build/protocolbuffers/go
    out: gen/go
    opt: paths=source_relative
  - remote: buf.build/grpc/go
    out: gen/go
    opt: paths=source_relative
inputs:
  - directory: proto
```

```proto
# service1_a.proto

syntax = "proto3";

package service1.v1;

import "service1/v1/service1_b.proto";
import "service1/v1/service1_c.proto";

option go_package = "gen/service1/v1";
...
...
```


```proto
# service2_a.proto

syntax = "proto3";

package service1.v1;

import "service2/v1/service2_b.proto";
import "service2/v1/service2_c.proto";

option go_package = "gen/service2/v1";
...
...
```


## How do I generate code for implementation?
Now suppose you want to create a repo that implements server for `service1`
- create your repo `https://github.com/jadeidev/service1`
- create `buf.gen.yaml` file
```yaml
# buf.gen.yaml
version: v2
clean: true
plugins:
  - remote: buf.build/protocolbuffers/go
    out: gen
    opt: paths=source_relative
  - remote: buf.build/grpc/go
    out: gen
    opt: paths=source_relative
inputs:
  - git_repo: https://github.com/jadeidev/my-proto-repo
    branch: main
    subdir: proto
    paths:
      - service1
```
- Run in shell `buf generate --clean`
- You will see something similar to this added to the repo:
```
├── gen
│   └── ags
│       └── v1
│           ├── service1_a.pb.go
│           ├── service1_a_grpc.pb.go
│           ├── service1_b.pb.go
│           └── service1_c.pb.go
```
