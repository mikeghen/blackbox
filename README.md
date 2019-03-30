<img src="https://raw.githubusercontent.com/crypdex/blackbox/master/resources/images/logo2.png" width=300>

# BlackboxOS

The BlackboxOS is a pluggable platform for deploying multi-chain applications. It is used as the basis for all [Crypdex's](https://crypdex.io) Blackbox devices.

Documentation is currently accruing on the [Wiki](https://github.com/crypdex/blackbox/wiki).

The CLI is available [here](https://github.com/crypdex/blackbox-cli).

## Features

<img src="https://raw.githubusercontent.com/crypdex/blackbox/master/resources/images/screenshot.png" width=330 align="left">

- Update framework
- Portable, Docker-based services
- Multiarch support: Runs on x86_64 and arm64 devices
- Optimized for multiple full nodes
- Unified multi-chain deterministic wallet
- Expandable with new chains
- Accessible via CLI, HTTP API, native RPCs, and GUI (under development)

## About

The BlackboxOS is being developed to support devices running blockchain systems. Many of the design decisions are based on the needs of running applications that require maintaining multiple blockchains and their supporting services, primarily in a cloudless environment.

The core of the system uses Docker containers and Docker Compose to orchestrate the "app". With this foundation, adapting to use Docker swarm mode or K8s is a possible direction based on needs.

There are companion admin and api services that support updates and multichain wallet functionality as well as the addition of new chains as they become available.

## Supported Chains

|          | service | wallet | features            |
| -------- | ------- | ------ | ------------------- |
| PIVX     | ✓       | ✓      | Masternode, staking |
| Dash     | ✓       | ✓      | Masternode          |
| ZCoin    | ✓       | ✓      | Masternode          |
| Blocknet | ✓       | ✓      | Masternode, staking |
| Bitcoin  |         |        |                     |
| Litecoin | ✓       |        |                     |
