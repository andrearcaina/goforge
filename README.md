# goforge

Opinionated CLI/TUI tool for generating Go backend service boilerplate

### Overview

Inspired by the fact that I hate rewriting and copy pasting code that I have already written before, so I decided to create a CLI tool to generate the needed files and folder structure to quickly create a new backend service.

The idea is that these are backend services (like microservices), not an entire backend monolith implementation. For starters though this generates a good strtucture, but for it to be an entire monolith it would need to require extra setup from yourself.

P.S. This was encountered while working on [fáfnir](https://github.com/andrearcaina/fafnir).

### Installation

`goforge` can be installed on Windows or Linux, WSL, Git Bash:

Install via cURL on Linux, WSL, or Git Bash:

```bash
curl -fsSL https://raw.githubusercontent.com/andrearcaina/goforge/main/install.sh | sh
```

Install via Powershell on Windows:

```bash
irm https://raw.githubusercontent.com/andrearcaina/goforge/main/install.ps1 | iex
```

Or via Go (requires Go 1.24.5 or later):

```bash
go install github.com/andrearcaina/goforge@latest
```

### CLI Usage

This CLI can generate three server types (REST, gRPC, GraphQL), along with optional database scaffolding and configuration (Makefile) scripts.

To generate a Go REST server, run:

```
> goforge generate -p rest-server -n my-server -s REST -m
```

This command does the following:

- Creates a directory named `rest-server`
- Initializes a `go.mod` file with the module name `my-server`
- Configures the server type as `REST`
- Generates a `Makefile`

### Interactive Mode (TUI)

If any required flags are missing, an interactive TUI will appear:
![goforge no db flag](images/goforge-no-db-flag.png)

You can see that the TUI appears, asking if you should generate the database files.

The following fields must be provided either via CLI flags or through the TUI (check [`internal/config/config.go`](internal/config/config.go)):

```go
type ServerTypeFlag string

const (
	REST    ServerTypeFlag = "rest"
	GRPC    ServerTypeFlag = "grpc"
	GraphQL ServerTypeFlag = "graphql"
)

type Form struct {
    Name           string
    ServerTypeFlag ServerTypeFlag
    DatabaseFlag   bool
    MakefileFlag   bool
    DockerFlag     bool
}
```

For non-interactive use (without opening the TUI options), explicitly provide each boolean option, including false values:

```bash
goforge generate \
  --path rest-server \
  --name my-server \
  --server rest \
  --database=false \
  --makefile=true \
  --docker=false
```

Docker generation requires database generation. If a generated file already exists, `goforge` preserves it and returns an error; pass `--force` to explicitly replace generated files.

For example, running with all required flags will produce:
![goforge all flags](images/goforge-all-flags.png)

Here is an example where you can trigger the TUI:\
![goforge TUI](images/goforge-tui.png)

### Roadmap

- [x] Allow different server type boilerplates
    - [x] REST (`go-chi`)
    - [x] gRPC (`grpc-go`)
    - [x] GraphQL (`gqlgen`)
- [x] Add Makefile support flag
- [x] Add SQLc/database support flag
    - [x] Potentially add a docker-compose file to spin up a PostgreSQL instance
- [x] Interactive UI for when certain flags are not provided
- [x] Add installation scripts for both Windows and Linux
