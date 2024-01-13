# clarum-core

The core package of the Clarum framework.

This package contains:

- configuration management
- test action controlling
- shared validators
- utilities
- orchestration (although temporarily until docker and k8s actions are added)

## Beta warning

The whole framework is currently in beta, so be aware that the API may change.
The framework can be used to write tests, for examples check the **itests**
in [clarum-http](https://github.com/go-clarum/clarum-http) and
in [clarum-samples](https://github.com/go-clarum/samples).

Currently test actions are running with the go standard test runner, but the actions are not really aware of their
context.
The framework does not know which action runs in which test and also when an integration test actually starts or ends.
These features are required for the initial production version and they most probably will change how the test actions
will be used & configured.

## Roadmap

### v1.0

#### Runtime

- standard go test runtime
- orchestration:
    - run custom commands (to start apps)
    - beforeSuite
    - afterSuite
- execution context
    - log what tests have been executed
    - which were successful or failed

#### Packages

- core
- http
- json

#### Protocols

- HTTP1.1
    - actions
        - requests / responses
    - endpoints:
        - client
        - server
    - validation
        - methods
        - paths
        - headers
        - query params
        - full body validation (JSON, plain, form)
        - `@ignore@` in payload
            - on values

#### Configuration

- read CLI flags (may not be possible with current setup):
    - setup custom config location
    - control active profile
- configuration file types:
    - yaml
- allow configuration based on profiles

#### Logging

- log actions per test
- log received/sent messages

#### Metrics

- receive action times response times
- execution time per test
- per test suite


### v1.x

#### Runtime

- global test suite variables

#### HTTP

- actions:
    - expect errors:
        - timeout
        - connection exception
- use variables in payloads
- validation:
    - XML payload validation (with @ignore@ for values)
    - specific path validation/ignore for XPath & JsonPath on receive actions
    - schema validation (XML & JSON)
    - OpenApi validation
- actions:
    - retry until success
    - retry until failure (simulate DDOS)
    - send multiple requests
        - n times configurable - w/o delay
        - time range - w/o delay

#### Metrics

- how fast the system under test sends a request/response
- test actions:
    - fail test/test suite if metrics were/were not reached
- test suite metrics report:
    - with history of previous runs
    - with delta compared to previous runs

#### Features

- test reports:
    - configurable
- static response server:
    - for different protocols
    - fixed / random / round robin responses from configured data
- remote server & client endpoints:
    - clarum http-endpoint instance deployed in k8s and controlled from localhost/pipeline
- docker:
    - API support and actions
- k8s:
    - API support and actions

#### Communication

- gRPC
- HTTP2/HTTP3
- Websockets
- Kafka
- NATS
- SQL/NoSQL
- GraphQL
