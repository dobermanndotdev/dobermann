<h1 style="text-align: center">Dobermann 🦮 [in beta]</h1>
<p style="text-align: center">Uptime monitoring</p>
<hr>

[![main](https://github.com/flowck/dobermann/actions/workflows/main.yml/badge.svg)](https://github.com/flowck/dobermann/actions/workflows/main.yml)

## Primary Components

- Core service ([./cmd/service](./cmd/service)): The main service responsible for all operations except for recurrent monitor's checks. 
- Worker ([./cmd/worker](./cmd/worker)): A service meant to be replicated across various regions with the goal of recurrently perform checks to the monitor's urls.

## Secondary Components

- PostgreSQL
- RabbitMQ

## Set up in Dev Env

Tools required: 

* Docker (Docker Compose)
* [Taskfile](taskfile.dev) - An enhanced `make/Makefile`

Set up:

* cd `./dobermann`
* boot up everything `docker-compose up` or `docker-compose up -d`
* run tests all kinds of tests `task test:all`

## Project Management

- [GitHub Projects Board](https://github.com/users/flowck/projects/10)

## Design and Modelling

- [Whimsical Board](https://whimsical.com/doberman-mvp-KSeahkKitCd9TMYf7M68ii)

## Continuous Integration and Deployment

* Merges to `main` trigger a tag-based [release](https://github.com/flowck/dobermann/releases) which the deployments are made from.  
* All artifacts are deployed to Fly.io via the GitHub Actions [.github/workflows/deploy-production.yml](.github/workflows/deploy-production.yml).

## License

```
                    GNU AFFERO GENERAL PUBLIC LICENSE
                       Version 3, 19 November 2007
                                   ...
```

The full license: [./LICENSE.md](./LICENSE.md)