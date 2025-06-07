# Exploding Kittens

Exploding Kittens is a distributed, service-oriented implementation of the popular card game. This repository contains the backend services, client SDK, protocol definitions, and documentation for running and developing the game.

---
## Getting Started

This section provides an overview of how the game is structured and how its core components interact. It’s a good place to begin if you're new to the project or planning to contribute.

- [Architecture Overview](./docs/architecture_overview.md)
- [Recommended Tooling](#recommended-tooling)
- [Running Locally](#running-locally)
- [Testing](#testing)

---
## Recommended Tooling

The recommended way to deal with tooling and versions is to use [asdf](https://asdf-vm.com/#/). This will allow you to install and manage multiple versions of the same tool on your machine. 
Additionally, [direnv](https://direnv.net/) is also recommended to manage local environment variables using .env and .env.local files (See samples)

### Installation
First install `asdf` and `direnv`, then reload your profile for changes to take effect.
Either restart your terminal application or run `source ~/.*rc`.

* MacOS:
```shell
brew install asdf direnv
source ~/.zshrc

asdf plugin-add direnv
asdf direnv setup --shell bash --version latest

cut -d' ' -f1 .tool-versions|xargs -I{} asdf plugin add {}
asdf install
asdf direnv allow
```

* Linux:
```shell
asdf plugin-add direnv
asdf direnv setup --shell bash --version latest

cut -d' ' -f1 .tool-versions|xargs -i asdf plugin add  {}
asdf install
asdf direnv allow
```

---
## Running Locally
<a name="running-locally"></a>

First copy `.env.sample` to `.env`, and fill in the appropriate / missing values.

```
cp .env.sample .env
```

---
### Building Containers
Build the Docker image if any changes have been made, or if it has never been built, in order to make it available for container creation. Build all of the images using the following command: 

```shell
make build
```

This will build all of the game services.

### Running all services locally with docker compose
Be sure to have successfully built the `exploding-kittens:latest` Docker image in the previous section before running Compose. If building for the first time, run the following to build all of the containers:
```shell
make compose-up
```

---
## Development
Here you'll find information on setting up your environment, auto-generating files for contributing to the project.

### Local Development

Run the following command to setup all necessary packages and tools
```shell
make go-deps (optional setup)
```

When new features are coming, don't forget to run the following command to auto-generate necessary files (proto code, API docs, and service config schema)
```shell
make gen
```

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
- The [unit-tests.yaml](./.github/workflows/unit-tests.yaml) defines all lower-level tests that will be run in CI, as well as an action to generate a code-coverage report.