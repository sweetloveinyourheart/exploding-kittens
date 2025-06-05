# Exploding Kittens

Exploding Kittens is a distributed, service-oriented implementation of the popular card game. This repository contains the backend services, client SDK, protocol definitions, and documentation for running and developing the game.

---
## Table Of Contents
- [Getting Started](#getting-started)
- [Testing](#testing)

---
## Getting Started

This section provides an overview of how the game is structured and how its core components interact. It’s a good place to begin if you're new to the project or planning to contribute.

- [Game Play](./docs/gameplay.md)
- [Architecture Overview](./docs/architecture_overview.md)
- [Services](./docs/services.md)
- [Data Flow](./docs/data_flow.md)

---
## Testing

Any change to the to the repository should be validated in CI before
a pull requests is merged into the main branch.

### Unit tests

Run check for all services
```shell
make test
```
... Or a specific service
```shell
make ut-[service_name] # eg. ut-userserver, ut-gameengineserver
```
... will run the linter and unit tests.

### Github CI

This repository contains the following GitHub CI workflows:
- The [env.yaml](./.github/workflows/env.yaml) is a reusable workflow used to setup an environment for running builds or tests.  Specific test workflows inject their work into the `commands-to-execute` variable.
- The [tests.yaml](./.github/workflows/tests.yaml) defines all lower-level tests that will be run in CI, as well as an action to generate a code-coverage report.