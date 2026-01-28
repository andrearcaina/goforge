# goforge

Opinionated CLI tool for generating Go backend service boilerplate

### Overview

Inspired by the fact that I hate rewriting and copy pasting code that I have already written before, so I decided to create a CLI tool to generate the needed files and folder structure to quickly create a new microservice. 

P.S. This was encountered while working on [fáfnir](https://github.com/andrearcaina/fafnir).

> The idea is that these are backend services (like microservices), not an entire backend monolith implementation. For starters though this generates a good strtucture, but for it to be an entire monolith it would need to require extra setup from yourself.


### TODO

In no particular order:

- [ ] Add REST, gRPC, and GraphQL templates to my liking
    - [X] REST
    - [ ] gRPC
    - [ ] GraphQL
- [X] Add SQLc support to my liking
- [X] Implement quality of life features like:
    - [X] Flags for every possible option
    - [X] Interactive UI for when flags are not provided