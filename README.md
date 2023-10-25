# Demography (test task)

Application provides simple operations with demographic data, enriched by free third party APIs.

## REST API specification

Fully-functional client published in Swagger UI on [GitHub Pages](https://barpav.github.io/demography-api).<br>
Application supports CORS, so requests can be made directly from GitHub Pages.

## How to launch it

Download or clone this repo and run:
```sh
make up
```

To stop the application run:
```sh
make down
```

<details>
    <summary>Why I use "sudo docker" instead of just "docker" in the Makefile</summary>

<br>

I use Docker Engine instead of Docker Desktop and according to the [Docker official documentation](https://docs.docker.com/engine/install/linux-postinstall/#manage-docker-as-a-non-root-user): <br>

> The Docker daemon binds to a Unix socket, not a TCP port. By default it's the root user that owns the Unix socket, and other users can only access it using sudo. The Docker daemon always runs as the root user.<br>
> <br>
> If you don't want to preface the docker command with sudo, create a Unix group called docker and add users to it

And: <br>

> The docker group grants root-level privileges to the user. For details on how this impacts security in your system, see [Docker Daemon Attack Surface](https://docs.docker.com/engine/security/#docker-daemon-attack-surface). <br>

But if you are uncomfortable with `sudo` for some reason and it's unnecessary for your system, instead of `make up` you may run:
```sh
docker-compose up -d
```

And instead of `make down`:
```sh
docker-compose down
```

</details>

## Logs

Just run 
```sh
make logs
```