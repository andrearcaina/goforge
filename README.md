# goforge

Opinionated CLI/TUI tool for generating Go backend service boilerplate

### Overview

Inspired by the fact that I hate rewriting and copy pasting code that I have already written before, so I decided to create a CLI tool to generate the needed files and folder structure to quickly create a new backend service.

The idea is that these are backend services (like microservices), not an entire backend monolith implementation. For starters though this generates a good strtucture, but for it to be an entire monolith it would need to require extra setup from yourself.

P.S. This was encountered while working on [fáfnir](https://github.com/andrearcaina/fafnir).

### TODO

In no particular order:

- [ ] Allow different server type boilerplates 
    - [x] REST (`go-chi`)
    - [ ] gRPC (`grpc-go`)
    - [ ] GraphQL (`gqlgen`)
- [x] Add SQLc support flag
    - [ ] Potentially add a docker-compose file to spin up a PostgreSQL instance
- [x] Interactive UI for when certain flags are not provided
