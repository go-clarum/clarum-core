# clarum
TODO:


### v1.0

#### Runtime
- standard go test runtime & orchestration

#### Packages
- core
- http

#### Protocols
- HTTP
    - actions
        - requests / responses
        - expect errors:
            - timeout
            - connection exception
    - endpoints:
        - client
        - server
    - validation
        - header validation
        - query params
        - full body validation (JSON, XML, plain, form)
        - `@ignore@` in payload
            - on values
            - on entire elements (JSON & XML)

#### API
- endpoints:
    - example APIs:
        - `.Receive(t *testing.T, action *Action)`
        - `.Receive(t *testing.T, action *Action) *http.Response`
        - `.Receive(action *Action) (*http.Response, types.Error)`

#### Configuration
- read CLI flags:
    - setup custom config location
    - control active profile
- different configuration file types:
    - json, yaml, properties
- allow configuration based on profiles

#### Logging
- log actions per test
- log received/sent messages

#### Metrics
- action times


### v1.1

#### Runtime
- global test suite variables

#### HTTP
- use variables in payloads
- validation
    - path validation/ignore (XPath, JsonPath)
    - schema check
    - openapi check

#### Metrics
- how fast the system under test sends a request/response
- test actions:
    - fail test if metrics were/were not reached


### v1.x

### HTTP

- Actions:
    - retry until success
    - retry until failure (simulate DDOS)
    - send multiple requests
        - n times configurable - w/o delay
        - time range - w/o delay

#### Features
- test reports:
  - configurable
- static response server:
    - for different protocols
    - fixed / random / round robin responses from configured data
    - should be deployable in the cloud
- docker:
    - API support and actions
- k8s:
    - API support and actions

#### Protocols
Websockets
Kafka
SQL
GraphQL
gRPC?
