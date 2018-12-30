# lightning-rebalancer
⚡️ Rebalance your Lightning Network channel using onchain BTC

This server implements a REST API interface for given LND gRPC daemon 
It provides endpoint to:
* Recieve a LN payment using BTC onchain 
* Pay a LN invoice to receive BTC onchain

## Development

These instructions are intended to facilitate the development and testing of Rebalancer. Operators interested in deploying Rebalancer should install the appropriate binary as
[released](#release)

### Prerequisites

* [Go 1.10+](https://golang.org/dl/)
* [Dep 0.5+](https://github.com/golang/dep#installation)

### Installing

* Clone the git repository

```bash
$ git clone git@github.com:vulpemventures/lightning-rebalancer.git
$ cd lightning-rebalancer
```

* Install development/test dependencies

```bash
$ ./scripts/install
```

* Build lightning-rebalancer (on MacOSX)
```bash
$ ./scripts/build darwin amd64
```

* Build lightning-rebalancer (on Linux amd64 platform)

```bash
$ ./scripts/build linux amd64
```

* Build lightning-rebalancer for Linux armv7 using Docker 

```bash
$ ./scripts/buildarm
```

#### Run Server
```bash
$ ./build/rebalancer-linux-amd64
```

The following is the list of variables that can be set to change the default configuration:

  * `LND_HOST`  - LND gRPC host
  * `LND_PORT`  - LND gRPC port
  * `HTTP_PORT` - HTTP server port to run
  * `TLS_PATH`  - absolute path of the TLS Certificate
  * `MAC_PATH`  -  absolute path of the Macaroon path


## API

TBD



